package domain

type AppFileUseCase interface {
	UploadFile(param map[string]interface{}) UseCaseResult
	RemoveFile(param map[string]interface{}) UseCaseResult
	DownloadFile(param map[string]interface{}) UseCaseResult
	PreviewFile(param map[string]interface{}) UseCaseResult
}

type AppFileRespository interface{}
