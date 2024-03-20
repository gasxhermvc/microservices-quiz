package usecase

import (
	"cpn-quiz-api-file-manage-go/domain"
	"cpn-quiz-api-file-manage-go/logger"
)

type appFileUseCase struct {
	appFileRespository domain.AppFileRespository
	log                *logger.PatternLogger
}

func NewAppFileUseCase(appFileRespository domain.AppFileRespository, log *logger.PatternLogger) domain.AppFileUseCase {
	return &appFileUseCase{
		appFileRespository: appFileRespository,
		log:                log,
	}
}

func (appFile *appFileUseCase) UploadFile(param map[string]interface{}) domain.UseCaseResult {
	response := domain.UseCaseResult{}
	return response
}

func (appFile *appFileUseCase) RemoveFile(param map[string]interface{}) domain.UseCaseResult {
	response := domain.UseCaseResult{}
	return response
}

func (appFile *appFileUseCase) DownloadFile(param map[string]interface{}) domain.UseCaseResult {
	response := domain.UseCaseResult{}
	return response
}

func (appFile *appFileUseCase) PreviewFile(param map[string]interface{}) domain.UseCaseResult {
	response := domain.UseCaseResult{}
	return response
}
