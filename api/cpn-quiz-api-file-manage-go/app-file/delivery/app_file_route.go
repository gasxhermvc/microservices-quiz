package delivery

import (
	"cpn-quiz-api-file-manage-go/domain"
	"cpn-quiz-api-file-manage-go/logger"

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

	eg.POST("/upload", handler.UploadFile)
	eg.DELETE("/remove", handler.RemoveFile)
	eg.GET("/download", handler.DownloadFile)
	eg.GET("/preview", handler.PreviewFile)
}
