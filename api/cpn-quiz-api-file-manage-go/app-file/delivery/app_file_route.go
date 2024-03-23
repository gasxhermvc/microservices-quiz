package delivery

import (
	"cpn-quiz-api-file-manage-go/domain"
	"cpn-quiz-api-file-manage-go/logger"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	config "github.com/spf13/viper"
)

type appFileDelivery struct {
	appFileUseCase domain.AppFileUseCase
	log            *logger.PatternLogger
	transId        string
}

func NewAppFileDelivery(e *echo.Echo, appFileUseCase domain.AppFileUseCase, log *logger.PatternLogger) {
	handler := &appFileDelivery{
		appFileUseCase: appFileUseCase,
		log:            log,
	}

	eg := e.Group(config.GetString("service.endpoint"))
	eConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(domain.Token)
		},
		SigningKey: []byte(config.GetString("dcc.api.jwt.token.sign")),
	}
	eg.Use(echojwt.WithConfig(eConfig))
	eg.POST("/upload", handler.UploadFile)
	eg.DELETE("/remove", handler.RemoveFile)
	eg.GET("/download", handler.DownloadFile)
	eg.GET("/preview", handler.PreviewFile)
}
