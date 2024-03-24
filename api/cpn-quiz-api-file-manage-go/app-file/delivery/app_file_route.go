package delivery

import (
	"cpn-quiz-api-file-manage-go/domain"
	"cpn-quiz-api-file-manage-go/logger"

	_perm "cpn-quiz-api-file-manage-go/middleware/permission-middleware"

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
		Skipper: func(c echo.Context) bool {
			clientId := c.Request().Header.Get("x-client-id")
			authorization := c.Request().Header.Get("x-api-key")
			return clientId != "" && authorization != ""
		},
	}
	eg.Use(echojwt.WithConfig(eConfig))

	perm := _perm.NewPermissionMiddleware(log)
	eg.Use(perm.AuthorizePermissions("cpn-quiz-api-file-manage-create"))
	//=>Dynamic upload file
	eg.POST("/upload", handler.UploadFile)

	eg.Use(perm.AuthorizePermissions("cpn-quiz-api-file-manage-delete"))
	//=>Dynamic Remove path within config remote path
	eg.DELETE("/remove", handler.RemoveFile)

	eg.Use(perm.AuthorizePermissions("cpn-quiz-api-file-manage-query"))
	//=>Download file
	eg.GET("/download", handler.DownloadFile)

	eg.Use(perm.AuthorizePermissions("cpn-quiz-api-file-manage-query"))
	//=>Preview file
	eg.GET("/preview", handler.PreviewFile)
}
