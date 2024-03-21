package domain

// _appDataService "gitlab.com/pea-developer/std/go/app-data-service"

type HelloWorldUseCase interface {
	// CronJob() error

	// //=>cpn Interface logs
	// InsertLogDetail(int, string) error
	// InsertLog() (int, error)
	// UpdateLog(int, int, int, int, int) error
	// GetTransactionID() string
	// SetTransactionID(transactionID string)
	// SetUsername(string)
	// CheckBackOfficeIsRunning() (bool, error)
}

type HelloWorldRespository interface {
	// Hello(param HelloParameter) (_appDataService.QueryResult, error)

	// InsertLogDetail(map[string]interface{}) error
	// InsertLog(map[string]interface{}) (int, error)
	// UpdateLog(map[string]interface{}) error
	// GetLogProgramStatus() ([]LogProgramStatusStore, error)
}

// =>Binder not sensitive case
// =>Ref: https://echo.labstack.com/guide/binding/
type HelloParameter struct {
	Name string
}

type HelloResponse struct {
	Name string
}

type UpdateResult struct {
	Success int `json:"SUCCESS"`
	Error   int `json:"ERROR"`
	Total   int `json:"TOTAL"`
}
