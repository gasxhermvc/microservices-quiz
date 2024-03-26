package delivery

import (
	"cpn-quiz-api-mailer-go/domain"
	"cpn-quiz-api-mailer-go/logger"
	"crypto/tls"
	"net/http"
	"strconv"
	"time"

	_perm "cpn-quiz-api-mailer-go/middleware/permission-middleware"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	config "github.com/spf13/viper"
	gomail "gopkg.in/mail.v2"
)

type emailDelivery struct {
	emailUseCase domain.EmailUseCase
	log          *logger.PatternLogger
	transId      string
}

func NewEmailDelivery(e *echo.Echo, emailUseCase domain.EmailUseCase, log *logger.PatternLogger) {
	handler := &emailDelivery{
		emailUseCase: emailUseCase,
		log:          log,
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
	//=>For check permissions
	matches := []_perm.MatchRoute{
		_perm.MatchRoute{
			Route:    "/send-email",
			Resource: "cpn-quiz-api-mailer-create",
		},
		_perm.MatchRoute{
			Route:    "/send-email-test",
			Resource: "cpn-quiz-api-mailer-create",
		},
		_perm.MatchRoute{
			Route:    "/attachments/*",
			Resource: "cpn-quiz-api-mailer-query",
		},
	}

	eg.Use(perm.AuthorizePermissions(matches))
	//=>ส่ง Email แบบ Dynamic
	eg.POST("/send-email", handler.SendEmail)
	//=>ทดสอบส่ง Email
	eg.POST("/send-email-test", func(ctx echo.Context) error {
		to := ctx.QueryParam("TO")
		m := gomail.NewMessage()

		// Set E-Mail sender
		m.SetHeader("From", config.GetString("cpn.quiz.api.mailer.email.username"))

		// Set E-Mail receivers
		m.SetHeader("To", to)

		// Set E-Mail subject
		m.SetHeader("Subject", "Gomail test subject")

		// Set E-Mail body. You can set plain text or html with text/html
		m.SetBody("text/plain", "This is Gomail test body")

		// Settings for SMTP server
		d := gomail.NewDialer(config.GetString("cpn.quiz.api.mailer.email.server"), config.GetInt("cpn.quiz.api.mailer.email.port"), config.GetString("cpn.quiz.api.mailer.email.username"), config.GetString("cpn.quiz.api.mailer.email.password"))

		if config.GetString("cpn.quiz.api.mailer.email.timeout") == "" {
			d.Timeout = time.Duration(180) * time.Second
		} else {
			timeout, _ := strconv.Atoi(config.GetString("cpn.quiz.api.mailer.email.timeout"))
			d.Timeout = time.Duration(timeout) * time.Second
		}
		// This is only needed when SSL/TLS certificate is not valid on server.
		// In production this should be set to false.
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

		// Now send E-Mail
		if err := d.DialAndSend(m); err != nil {
			log.Error("", "DialAndSend", err)
			panic(err)
		}
		return ctx.JSON(http.StatusOK, nil)
	})
	//=>ดูภาพเอกสารไฟล์แนบ กรณีจะทำ Public endpoint สำหรับ Email
	eg.GET("/attachments/*", handler.AttachmentFile)
}
