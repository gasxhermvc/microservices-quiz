//=>Support for SQL Server Only
package appdataservice

import (
	"errors"
	"fmt"

	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type AppDataService struct {
	db *gorm.DB
}

func NewAppDataService(DB *gorm.DB) IAppDataService {
	return &AppDataService{
		db: DB,
	}
}

//=>คิวรี่ SP
func (appDataService AppDataService) ExecuteProcedure(queryParameter QueryParameter) QueryResult {
	result := QueryResult{}

	//=>Step1
	//===>Apply prefix PI_ to parameter
	parameters := AddPrefixParameter(queryParameter.Parameters)

	//=>Step2
	//===>Get store procedure name
	spName := parameters["APP_DATA_PROCEDURE"]

	if spName == "" || spName == nil {
		err := errors.New("{APP_DATA_PROCEDURE} parameter doesn't exist or empty value in the QueryParameter")
		return CustomQueryError(result, err)
	}

	//===>Delete APP_DATA_PROCEDURE from parameter
	delete(parameters, "APP_DATA_PROCEDURE")

	//=>Step3
	//===>Get schema of store procedure
	schemaName, err := GetSchemeOfStoreProcedure(appDataService.db, spName.(string))
	if err != nil {
		return CustomQueryError(result, err)
	}

	//=>Step4
	//===>Get parameters of store procedure
	directionParameters, err := GetStoreProcedureParameter(appDataService.db, spName.(string))
	if err != nil {
		return CustomQueryError(result, err)
	}

	//=>Step5
	//===>Render Execute Command
	rawSqlCommand := GeneratorExecuteStoreProcedureScript(schemaName, spName.(string), parameters, directionParameters)

	//=>Step6
	//===>Prepare Statement & Execute
	//===>Get binding values
	bindings := PrepareBindings(parameters, directionParameters)

	//===>Prepare statement
	rows, err := appDataService.db.Raw(rawSqlCommand, bindings...).Rows()
	if err != nil {
		return CustomQueryError(result, err)
	}

	//=>Step7
	//===>Add dataset to outputParameter
	outputParameters := make(map[string]interface{})

	var isError bool
	var exception error

	resultDataSets := [][]map[string]interface{}{}

	for cont := true; cont; cont = rows.NextResultSet() {
		if isError {
			break
		}

		dataSets := []map[string]interface{}{}
		for rows.Next() {
			data := make(map[string]interface{})

			scanErr := appDataService.db.ScanRows(rows, &data)
			if err != nil {
				fmt.Println(scanErr)
				isError = true
				exception = scanErr
				break
			}

			if len(data) > 0 {
				dataSets = append(dataSets, data)
			}
		}

		resultDataSets = append(resultDataSets, dataSets)
	}

	for idx, datasets := range resultDataSets {
		columns := []string{}

		//=>ดึงชื่อ Key เพื่อนำมาใช้เช็คชื่อ Column ของ Result แต่ละ Dataset
		for _, data := range datasets {
			//=>หาเฉพาะแถวแรกแล้วจบการวนลูป
			for k := range data {
				columns = append(columns, k)
			}
			break
		}

		if !slices.Contains(columns, "@PO_STATUS") && !slices.Contains(columns, "Return Value") {
			//=>ตรวจสอบหากไม่พบ คอลัมน์ชื่อ @PO_STATUS และ Return Value ให้ทำงานใน Statement ส่วนนี้
			if (idx + 1) == 1 {
				outputParameters["Data"] = datasets
				result.Total = len(datasets)
			} else {
				outputParameters[fmt.Sprintf("%s%d", "Data", (idx+1))] = datasets
			}
		} else if slices.Contains(columns, "@PO_STATUS") {
			//=>ตรวจสอบหากพบ คอลัมน์ชื่อ @PO_STATUS ให้ทำงานใน Statement ส่วนนี้
			for _, data := range datasets {
				//=>ดึงข้อมูลผลลัพธ์การคิวรี่ หรือ Error Message ต่างๆในกรณีพบ เฉพาะแถวแรกและจบการวนลูป
				for key, value := range data {
					if key == "@PO_STATUS" { //=>
						//=>เก็บค่าสถานะการคิวรี่ของ @PO_STATUS ไปที่ outputParameter["Success"]
						outputParameters["Success"] = value
					}

					if key == "@PO_STATUS_MSG" {
						//=>เก็บค่าข้อความของสถานะการคิวรี่ของ @PO_STATUS_MSG ไปที่ outputParameter["Message"]
						outputParameters["Message"] = value
					}
				}
				break
			}
		}
	}

	//=>Step8
	//===>Logic to check query success
	if isError {
		return CustomQueryError(result, exception)
	}

	if value, ok := outputParameters["Success"]; ok {
		if value != nil {
			result.Success = value.(int64) == 1

		} else {
			result.Success = false
		}
	} else {
		result.Success = false
	}

	if value, ok := outputParameters["Message"]; ok {
		if value != nil {
			result.Message = value.(string)
		} else {
			result.Message = ""
		}
	}

	//===>Delete Success & Message
	delete(outputParameters, "Success")
	delete(outputParameters, "Message")
	result.OutputParameter = outputParameters
	return result
}

//=>ดึงข้อมูลชื่อ SCHEMA ของ SP
func GetSchemeOfStoreProcedure(DB *gorm.DB, storeProcedureName string) (string, error) {

	rows, err := DB.Raw(`SELECT s.name as SCHEMA_NAME, pr.name as SP_NAME FROM sys.procedures pr INNER JOIN sys.schemas s ON pr.schema_id = s.schema_id WHERE pr.name = ?`, storeProcedureName).Rows()

	if err != nil {
		return "", err
	}

	var schemas []SchemaProcedure

	for rows.Next() {
		DB.ScanRows(rows, &schemas)
	}

	if len(schemas) == 0 {
		err := errors.New("The schema of stored procedure " + storeProcedureName + " doesn't exist.")
		return "", err
	}

	return schemas[0].SCHEMA_NAME, nil
}

//=>ดึงข้อมูลชื่อ PARAMETER_IN&PARAMETER_OUT ของ SP
func GetStoreProcedureParameter(DB *gorm.DB, storeProcedureName string) ([]DirectionParameter, error) {
	rows, err := DB.Raw(`SELECT o.name as SP_NAME, p.name as PARAMETER_NAME, p.is_output as IS_OUTPUT,p.is_nullable as IS_ALLOW_NULL,p.is_readonly as IS_READONLY,p.system_type_id as SYS_TYPE_ID, t.name as SYS_TYPE_NAMEE
		FROM sys.all_parameters p inner join sys.all_objects o on p.object_id = o.object_id inner join sys.types t on t.system_type_id = p.system_type_id
		WHERE o.type = ? and o.name = ? and p.user_type_id = t.user_type_id`, "P", storeProcedureName).Rows()

	if err != nil {
		return nil, err
	}

	var parameters []DirectionParameter

	for rows.Next() {
		DB.ScanRows(rows, &parameters)
	}

	if len(parameters) == 0 {
		err := errors.New("The stored procedure " + storeProcedureName + " doesn't exist.")
		return nil, err
	}

	return parameters, nil
}
