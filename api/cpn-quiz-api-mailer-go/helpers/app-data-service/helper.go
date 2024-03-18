package appdataservice

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

//=>ข้ามการกรอง PI สำหรับ KEY ที่ตรงกับชุด Array lists
func IgnoreParameterName() []string {
	return []string{"app_data_procedure"}
}

//=>Add Prefix PI_
//=>Example: INPUT: DATA|PI_DATA, OUTPUT: PI_DATA|PI_DATA
func AddPrefixParameter(parameters map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	ignoreParameter := IgnoreParameterName()

	for k, v := range parameters {
		if slices.Contains(ignoreParameter, strings.ToLower(k)) {
			result[k] = v
			continue
		}

		if !strings.HasPrefix(k, "PI_") {
			result["PI_"+k] = v
		} else {
			result[k] = v
		}
	}

	return result
}

//=>Clear Prefix PI_
//=>Example: INPUT: PI_DATA, OUTPUT: DATA
func CleanPrefixParameter(parameters []map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	ignoreParameter := IgnoreParameterName()

	for _, m := range parameters {
		for k, v := range m {
			if slices.Contains(ignoreParameter, strings.ToLower(k)) {
				result[k] = v
				continue
			}

			if strings.HasPrefix(k, "PI_") {
				num := len(k)
				result[k[3:num]] = v
			} else {
				result[k] = v
			}
		}
	}

	return result
}

//=>แปลงคำสั่ง SQL Command สำหรับ Execute Store Procedure
func GeneratorExecuteStoreProcedureScript(schemaName string,
	storeProcedureName string,
	parameters map[string]interface{},
	directionParameters []DirectionParameter) string {

	keys := []string{}
	for k := range parameters {
		keys = append(keys, "@"+strings.ToLower(k))
	}

	declare := `DECLARE	@return_value int,
	@PO_STATUS int,
	@PO_STATUS_MSG nvarchar(max)`

	exec := fmt.Sprintf("\r\n\r\nEXEC @return_value = [%s].[%s]\r\n", schemaName, storeProcedureName)

	returnValue := `SELECT	@PO_STATUS as N'@PO_STATUS',
	@PO_STATUS_MSG as N'@PO_STATUS_MSG'
	
	SELECT	'Return Value' = @return_value`

	sqlCommand := declare
	sqlCommand += exec

	commandParameter := []string{}

	for _, direction := range directionParameters {
		_parameterName := strings.ToLower(direction.PARAMETER_NAME)
		_parameterInput := _parameterName[1:len(direction.PARAMETER_NAME)]

		if slices.Contains(keys, _parameterName) && !direction.IS_OUTPUT {
			parameterName := direction.PARAMETER_NAME[1:len(direction.PARAMETER_NAME)]

			if _, key := parameters[parameterName]; key {
				commandParameter = append(commandParameter, fmt.Sprintf("%s = ?", direction.PARAMETER_NAME))
				continue
			}

			for pName := range parameters {
				if strings.ToLower(pName) == _parameterInput {
					commandParameter = append(commandParameter, fmt.Sprintf("%s = ?", direction.PARAMETER_NAME))
					break
				}
			}
		}

		//=>Followup declare
		if direction.IS_OUTPUT {
			commandParameter = append(commandParameter, fmt.Sprintf("%s = %s OUTPUT", direction.PARAMETER_NAME, direction.PARAMETER_NAME))
		}
	}

	sqlCommand += strings.Join(commandParameter, ",\r\n")
	sqlCommand += "\n\r\r\n"

	sqlCommand += returnValue

	return sqlCommand
}

//=>ดึง Binding Variable ที่จะใช้ในการ Execute Store Procedure
func PrepareBindings(parameters map[string]interface{},
	directionParameters []DirectionParameter) []interface{} {
	bindings := []interface{}{}

	keys := []string{}

	for k := range parameters {
		keys = append(keys, "@"+strings.ToLower(k))
	}

	for _, direction := range directionParameters {
		_parameterName := strings.ToLower(direction.PARAMETER_NAME)
		_parameterInput := _parameterName[1:len(direction.PARAMETER_NAME)]

		if slices.Contains(keys, _parameterName) && !direction.IS_OUTPUT {
			parameterName := direction.PARAMETER_NAME[1:len(direction.PARAMETER_NAME)]

			if value, key := parameters[parameterName]; key {
				bindings = append(bindings, value)
				continue
			}

			for pName, pValue := range parameters {
				if strings.ToLower(pName) == _parameterInput {
					bindings = append(bindings, pValue)
					break
				}
			}
		}
	}

	return bindings
}

//=>จัดการ Query Error
func CustomQueryError(queryResult QueryResult, err error) QueryResult {
	queryResult.Message = "Cannot call store procedure."

	if err != nil {
		queryResult.Message = err.Error()
		queryResult.Error = err
	}

	ClearMapValue(queryResult.OutputParameter)
	queryResult.Total = 0
	queryResult.Success = false

	return queryResult
}

func ClearMapValue(maps map[string]interface{}) {
	for k := range maps {
		delete(maps, k)
	}
}
