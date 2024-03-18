package repository

import (
	"cpn-quiz-schedule-messenger-go/domain"

	"gorm.io/gorm"
)

type helloWorldRepository struct {
	db *gorm.DB
}

func NewHelloWorldRepository(db *gorm.DB) domain.HelloWorldRespository {
	return &helloWorldRepository{
		db: db,
	}
}

// func (hello helloWorldRepository) Hello(param domain.HelloParameter) (_appDataService.QueryResult, error) {
// 	response := _appDataService.QueryResult{}

// 	response.OutputParameter = make(map[string]interface{})
// 	response.Success = true
// 	response.AddOutputParameter("Data", []map[string]interface{}{{"Name": "john smith!"}})
// 	response.Total = 1

// 	return response, nil
// }

// func (r helloWorldRepository) InsertLogDetail(params map[string]interface{}) error {
// 	var PARAMETER_IN = make(map[string]interface{})
// 	utils.Transform(params, &PARAMETER_IN)
// 	PARAMETER_IN["APP_DATA_PROCEDURE"] = "WEB_A_CPM_LOG_PROGRAM_DETAIL"
// 	queryParameter := _appDataService.NewQueryParameter(PARAMETER_IN)
// 	appDataService := _appDataService.NewAppDataService(r.db)
// 	result := appDataService.ExecuteProcedure(queryParameter)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	return nil
// }

// func (r helloWorldRepository) InsertLog(params map[string]interface{}) (int, error) {
// 	// params
// 	params["APP_DATA_PROCEDURE"] = "WEB_A_CPM_LOG_PROGRAM_STATUS"
// 	queryParameter := _appDataService.NewQueryParameter(params)
// 	appDataService := _appDataService.NewAppDataService(r.db)
// 	result := appDataService.ExecuteProcedure(queryParameter)
// 	if result.Error != nil {
// 		return 0, result.Error
// 	}

// 	var id int
// 	var dataDetail []domain.DataDetailStore
// 	_, errTrasform := utils.Transform(result.GetOutputParameter("Data"), &dataDetail)

// 	if errTrasform != nil {
// 		return id, errTrasform
// 	}
// 	if result.Total == 0 {
// 		return id, errTrasform
// 	}

// 	if result.Error != nil {

// 		return id, result.Error
// 	}

// 	parseID, errParse := strconv.Atoi(dataDetail[0].ID)
// 	if errParse != nil {
// 		return id, errParse
// 	}

// 	id = parseID

// 	return id, nil
// }

// func (r helloWorldRepository) UpdateLog(params map[string]interface{}) error {
// 	// params
// 	params["APP_DATA_PROCEDURE"] = "WEB_U_CPM_LOG_PROGRAM_STATUS"
// 	queryParameter := _appDataService.NewQueryParameter(params)
// 	appDataService := _appDataService.NewAppDataService(r.db)
// 	result := appDataService.ExecuteProcedure(queryParameter)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	return nil
// }

// func (r helloWorldRepository) GetLogProgramStatus() ([]domain.LogProgramStatusStore, error) {
// 	var PARAMETER_IN = make(map[string]interface{})

// 	PARAMETER_IN["APP_DATA_PROCEDURE"] = "WEB_Q_CPM_LOG_PROGRAM_STATUS"
// 	PARAMETER_IN["PROGRAM_ID"] = config.GetInt("cpm.application.programId")

// 	queryParameter := _appDataService.NewQueryParameter(PARAMETER_IN)
// 	appDataService := _appDataService.NewAppDataService(r.db)
// 	result := appDataService.ExecuteProcedure(queryParameter)
// 	// have error
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}

// 	var logProgramStatusStore []domain.LogProgramStatusStore
// 	_, errTransForm := utils.Transform(result.GetOutputParameter("Data"), &logProgramStatusStore)

// 	// transform error
// 	if errTransForm != nil {
// 		return nil, errTransForm
// 	}

// 	// not found
// 	if len(logProgramStatusStore) == 0 {
// 		return nil, nil
// 	}
// 	return logProgramStatusStore, nil
// }
