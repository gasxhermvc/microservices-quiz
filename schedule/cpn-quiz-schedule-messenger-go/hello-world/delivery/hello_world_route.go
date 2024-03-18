package delivery

import (
	"cpn-quiz-schedule-messenger-go/domain"
	"cpn-quiz-schedule-messenger-go/logger"

	"github.com/labstack/echo/v4"
)

type helloWorldDelivery struct {
	helloWorldUsecase domain.HelloWorldUseCase
	log               *logger.PatternLogger
	transId           string
}

func NewHelloWorldDelivery(e *echo.Echo, helloWorldUsecase domain.HelloWorldUseCase, log *logger.PatternLogger) {
	// transId := uuid.New().String()
	// handler := &helloWorldDelivery{
	// 	helloWorldUsecase: helloWorldUsecase,
	// 	log:               log,
	// 	transId:           transId,
	// }
	// handler.helloWorldUsecase.SetTransactionID(transId)

	// eg := e.Group(config.GetString("service.endpoint"))
	// eConfig := middleware.JWTConfig{
	// 	Claims:     &domain.Token{},
	// 	SigningKey: []byte(config.GetString("cpm.api.jwt.token.sign")),
	// }
	// eg.Use(middleware.JWTWithConfig(eConfig))
	// eg.GET("/triggerHelloWorld", handler.TriggerJob)
}
