package appdataservice

type IAppDataService interface {
	ExecuteProcedure(queryParameter QueryParameter) QueryResult
}

type SchemaProcedure struct {
	SCHEMA_NAME string
	SP_NAME     string
}

type DirectionParameter struct {
	SP_NAME        string
	PARAMETER_NAME string
	IS_OUTPUT      bool
	IS_ALLOW_NULL  bool
	IS_READONLY    bool
	SYS_TYPE_ID    int64
	SYS_TYPE_NAME  string
}
