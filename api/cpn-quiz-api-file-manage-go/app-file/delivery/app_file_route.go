package delivery

import (
	"cpn-quiz-api-file-manage-go/domain"
	"cpn-quiz-api-file-manage-go/logger"

	_perm "cpn-quiz-api-file-manage-go/middleware/permission-middleware"
	_auth "cpn-quiz-api-file-manage-go/middleware/service-authorize"

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

	auth := _auth.NewCustomAuthorizeGuardOnly("cpn-file-manage", log)
	perm := _perm.NewPermissionMiddleware(log)
	//=>For check permissions
	matches := []_perm.MatchRoute{
		_perm.MatchRoute{
			Route:    "/upload",
			Resource: "cpn-quiz-api-file-manage-create",
		},
		_perm.MatchRoute{
			Route:    "/remove",
			Resource: "cpn-quiz-api-file-manage-delete",
		},
		_perm.MatchRoute{
			Route:    "/download",
			Resource: "cpn-quiz-api-file-manage-query",
		},
		_perm.MatchRoute{
			Route:    "/preview",
			Resource: "cpn-quiz-api-file-manage-query",
		},
	}

	eg.Use(perm.AuthorizePermissions(matches))
	eg.Use(auth.AuthorizeGuard())
	//=>Dynamic upload file
	eg.POST("/upload", handler.UploadFile)
	//=>Dynamic Remove path within config remote path
	eg.DELETE("/remove", handler.RemoveFile)
	//=>Download file
	eg.GET("/download", handler.DownloadFile)
	//=>Preview file
	eg.GET("/preview", handler.PreviewFile)
}
