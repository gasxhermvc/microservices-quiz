package usecase

import (
	"cpn-quiz-api-authentication-go/domain"
	"cpn-quiz-api-authentication-go/logger"
	"cpn-quiz-api-authentication-go/utils"
	"net/http"
)

type helloWorldUseCase struct {
	helloWorldRepository domain.HelloWorldRepository
	log                  *logger.PatternLogger
}

func NewHelloWorldUseCase(helloWorldRepository domain.HelloWorldRepository, log *logger.PatternLogger) domain.HelloWorldUseCase {
	return &helloWorldUseCase{
		helloWorldRepository: helloWorldRepository,
		log:                  log,
	}
}

func (hello *helloWorldUseCase) Hello(param domain.HelloParameter) domain.UseCaseResult {
	response := domain.UseCaseResult{}

	result := hello.helloWorldRepository.Hello(param)
	response.Success = result.Success
	response.Message = result.Message

	if result.Error != nil {
		response.Error = result.Error
		response.StatusCode = http.StatusUnprocessableEntity
		return response
	}

	data := []domain.HelloResponse{}
	utils.Transform(result.GetOutputParameter("Data"), &data)

	response.Result = data[len(data)-1]
	response.StatusCode = http.StatusOK
	return response
}
