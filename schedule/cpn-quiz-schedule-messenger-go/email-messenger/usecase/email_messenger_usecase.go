package usecase

import (
	"context"
	"cpn-quiz-schedule-messenger-go/constant"
	"cpn-quiz-schedule-messenger-go/database"
	"cpn-quiz-schedule-messenger-go/domain"
	"cpn-quiz-schedule-messenger-go/helpers/restful-service"
	"cpn-quiz-schedule-messenger-go/logger"
	"crypto/tls"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	config "github.com/spf13/viper"
	gomail "gopkg.in/mail.v2"
)

var rest *restful.Restful = restful.NewRestful(new(logger.PatternLogger), uuid.New().String())
var ctx = context.Background()

type emailMessengerUseCase struct {
	emailMessengerRepository domain.EmailMessengerRepository
	rdbCon                   *database.RedisDatabase
	log                      *logger.PatternLogger
	transId                  string
}

func NewEmailMessengerUseCase(emailMessengerRepository domain.EmailMessengerRepository, rdbConn *database.RedisDatabase, log *logger.PatternLogger) domain.EmailMessengerUseCase {
	return &emailMessengerUseCase{
		emailMessengerRepository: emailMessengerRepository,
		rdbCon:                   rdbConn,
		log:                      log,
	}
}

func (sm *emailMessengerUseCase) SetTransaction(transId string) {
	sm.transId = transId
}

func (sm *emailMessengerUseCase) SendEmail(provider *domain.EmailProvider) domain.UseCaseResult {
	response := domain.UseCaseResult{}

	if provider.Config.Skip {
		//=>จำลองการส่งเมล
		response.Success = true
		response.Message = constant.Success
		response.StatusCode = constant.SuccessCode
		return response
	}

	delimit := ";"
	mail := gomail.NewMessage()

	from := strings.Split(provider.Config.Username, delimit)
	sendAddress := provider.Config.SenderAddress
	isOverrideFrom := false
	if provider.Parameter.From != "" {
		_from := strings.Split(provider.Parameter.From, delimit)
		from = _from
		isOverrideFrom = true
	}

	if isOverrideFrom {
		sendAddress = ""
	}

	//=>setup from
	if len(from) > 0 {
		if isOverrideFrom {
			mail.SetHeaders(buildMailAddress("From", from))
		} else {
			if len(from) == 1 {
				mail.SetHeader("From", mail.FormatAddress(from[0], sendAddress))
			} else {
				mail.SetHeaders(buildMailAddress("From", from))
			}
		}
	}

	//=>setup to
	to := strings.Split(provider.Parameter.To, delimit)
	if len(to) > 0 {
		mail.SetHeaders(buildMailAddress("To", to))
	}

	//=>setup cc
	cc := strings.Split(provider.Parameter.Cc, delimit)
	if len(cc) > 0 {
		mail.SetHeaders(buildMailAddress("Cc", cc))
	}

	//=>setup bcc
	bcc := strings.Split(provider.Parameter.Bcc, delimit)
	if len(bcc) > 0 {
		mail.SetHeaders(buildMailAddress("Bcc", bcc))
	}

	//=>setup subject
	subject := provider.Parameter.Subject
	if len(subject) > 1 {
		mail.SetHeader("Subject", subject)
	}

	//=>setup body
	body := provider.Parameter.Body
	if len(body) > 1 {
		contentType := "text/plain"

		//=>หากมีการ Required ส่งแบบ html จะกำหนด content type เป็น text/html
		if provider.Parameter.IsHtml {
			contentType = "text/html"
		}
		mail.SetBody(contentType, body)
	}

	//=>setup attachments
	var attachErrors []string
	if len(provider.Parameter.Attachment) > 0 {
		for _, value := range provider.Parameter.Attachment {

			buf, statusCode, err := rest.HttpGet(value.DownloadUrl, nil, value.Filename)
			if err != nil {
				attachErrors = append(attachErrors, err.Error())
				continue
			}

			if statusCode != 200 {
				attachErrors = append(attachErrors, err.Error())
				continue
			}

			//=>แนบไฟล์
			mail.Attach(value.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := w.Write(buf)
				attachErrors = append(attachErrors, err.Error())
				return err
			}))
		}
	}

	dial := gomail.NewDialer(provider.Config.Server,
		provider.Config.Port,
		provider.Config.Username,
		provider.Config.Password)

	if config.GetString("cpn.quiz.api.mailer.email.timeout") == "" {
		dial.Timeout = time.Duration(180) * time.Second
	} else {
		timeout, _ := strconv.Atoi(config.GetString("cpn.quiz.api.mailer.email.timeout"))
		dial.Timeout = time.Duration(timeout) * time.Second
	}

	dial.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := dial.DialAndSend(mail); err != nil {
		response.Success = false
		response.Message = constant.ServiceUnavailable
		response.StatusCode = constant.ServiceUnavailableCode
		response.Errors = append(response.Errors, err.Error())
		return response
	}

	response.Success = true
	response.Message = constant.Success
	response.StatusCode = constant.SuccessCode

	return response
}
