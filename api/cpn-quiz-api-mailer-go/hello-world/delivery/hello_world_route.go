package delivery

import (
	"cpn-quiz-api-mailer-go/constant"
	"cpn-quiz-api-mailer-go/domain"
	"cpn-quiz-api-mailer-go/logger"
	"net/http"

	"github.com/labstack/echo/v4"
	config "github.com/spf13/viper"
)

type helloWorldDelivery struct {
	helloWorldUseCase domain.HelloWorldUseCase
	log               *logger.PatternLogger
	transId           string
}

func NewHelloWorldDelivery(e *echo.Echo, helloWorldUseCase domain.HelloWorldUseCase, log *logger.PatternLogger) {
	handler := &helloWorldDelivery{
		helloWorldUseCase: helloWorldUseCase,
		log:               log,
	}

	eg := e.Group(config.GetString("service.endpoint"))

	eg.GET("/hello", handler.Hello)

	eg.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})
}

func (hello helloWorldDelivery) DoResponse(c echo.Context, useCase domain.UseCaseResult, response domain.Response) error {
	switch useCase.StatusCode {
	case 200:
		response.Message = constant.Success
		response.ResponseData = useCase.Result
		break
	case 400:
		response.Message = constant.BadRequest
		break
	case 401:
		response.Message = constant.UnAuthorization
		break
	case 403:
		response.Message = constant.AccessDenied
		break
	case 404:
		response.Message = constant.NotFound
		response.Message = constant.Success
		response.ResponseData = useCase.Result
		break
	case 500:
		response.Message = constant.InternalServerError
		break
	case 503:
		response.Message = constant.ServiceUnavailable
		break
	default:
		response.Message = useCase.Message
	}

	return c.JSON(useCase.StatusCode, response)
}
