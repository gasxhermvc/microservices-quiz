package delivery

import (
	"cpn-quiz-api-authentication-go/domain"
	"cpn-quiz-api-authentication-go/logger"

	echojwt "github.com/labstack/echo-jwt/v4"

	"github.com/labstack/echo/v4"
	config "github.com/spf13/viper"
)

type authDelivery struct {
	authUsecase domain.AuthUseCase
	log         *logger.PatternLogger
	transId     string
}

func NewAuthDelivery(e *echo.Echo, authUsecase domain.AuthUseCase, log *logger.PatternLogger) {
	handler := &authDelivery{
		authUsecase: authUsecase,
		log:         log,
	}

	r := e.Group(config.GetString("service.endpoint"))

	g := r.Group("/auth")
	eConfig := echojwt.Config{
		KeyFunc: getKey,
	}

	g.Use(echojwt.WithConfig(eConfig))
	g.POST("/token", handler.AuthToken)
}
