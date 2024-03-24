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
		SigningKey: []byte(config.GetString("cpn.quiz.api.jwt.secretkey")),
	}
	eg.Use(echojwt.WithConfig(eConfig))

	//=>Dynamic upload file
	eg.POST("/upload", handler.UploadFile)

	//=>Dynamic Remove path within config remote path
	eg.DELETE("/remove", handler.RemoveFile)

	//=>Download file
	eg.GET("/download", handler.DownloadFile)

	//=>Preview file
	eg.GET("/preview", handler.PreviewFile)
}
