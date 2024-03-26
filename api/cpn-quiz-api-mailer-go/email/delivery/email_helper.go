package delivery

import (
	"cpn-quiz-api-mailer-go/domain"
	"fmt"
	"html"
	"mime/multipart"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/labstack/echo/v4"
	config "github.com/spf13/viper"
)

func buildEmailConfiguration() *domain.EmailConfig {
	emailConfig := new(domain.EmailConfig)
	emailConfig.Server = config.GetString("cpn.quiz.api.mailer.email.server")
	emailConfig.Port = config.GetInt("cpn.quiz.api.mailer.email.port")
	emailConfig.EnableSSL = config.GetBool("cpn.quiz.api.mailer.email.enablessl")
	emailConfig.SendWithCredential = config.GetBool("cpn.quiz.api.mailer.email.sendwithcredential")
	emailConfig.Skip = config.GetBool("cpn.quiz.api.mailer.email.skip")
	emailConfig.DefaultCredential = config.GetBool("cpn.quiz.api.mailer.email.defaultcredential")
	emailConfig.Username = config.GetString("cpn.quiz.api.mailer.email.username")
	emailConfig.Password = config.GetString("cpn.quiz.api.mailer.email.password")
	emailConfig.SenderAddress = config.GetString("cpn.quiz.api.mailer.email.senderaddress")
	emailConfig.FromParameter = config.GetString("cpn.quiz.api.mailer.email.fromparameter")
	emailConfig.ToParameter = config.GetString("cpn.quiz.api.mailer.email.toparameter")
	emailConfig.CcParameter = config.GetString("cpn.quiz.api.mailer.email.ccparameter")
	emailConfig.BccParameter = config.GetString("cpn.quiz.api.mailer.email.bccparameter")
	emailConfig.SubjectParameter = config.GetString("cpn.quiz.api.mailer.email.subjectparameter")
	emailConfig.BodyParameter = config.GetString("cpn.quiz.api.mailer.email.bodyparameter")
	emailConfig.PriorityParameter = config.GetString("cpn.quiz.api.mailer.email.priorityparameter")
	emailConfig.AttachmentFileParameter = config.GetString("cpn.quiz.api.mailer.email.attachmentparameter")

	return emailConfig
}

func buildEmailProvider(configuration domain.EmailConfig, c echo.Context) (*domain.EmailProvider, error) {
	emailProvider := new(domain.EmailProvider)

	emailProvider.Config = configuration

	parameter := new(domain.EmailParameter)
	parameter.From = html.EscapeString(c.Request().FormValue(configuration.FromParameter))
	parameter.To = html.EscapeString(c.Request().FormValue(configuration.ToParameter))
	parameter.Cc = html.EscapeString(c.Request().FormValue(configuration.CcParameter))
	parameter.Bcc = html.EscapeString(c.Request().FormValue(configuration.BccParameter))
	parameter.Subject = html.EscapeString(c.Request().FormValue(configuration.SubjectParameter))
	parameter.Body = html.EscapeString(c.Request().FormValue(configuration.BodyParameter))
	parameter.Priority = html.EscapeString(c.Request().FormValue(configuration.PriorityParameter))

	isHTML, err := strconv.ParseBool(html.EscapeString(c.Request().FormValue("is_html")))
	if err != nil {
		isHTML = false
	}
	parameter.IsHtml = isHTML
	form, err := c.MultipartForm()
	if err != nil {
		return emailProvider, err
	}

	//=>Has files
	files := form.File[configuration.AttachmentFileParameter]
	if len(files) > 0 {
		parameter.Attachment = files
	}

	emailProvider.Parameter = *parameter

	return emailProvider, nil
}

// =>อัปโหลดจำนวนไฟล์เกิน
func isExceedLimitUpload(fileCollection map[string][]*multipart.FileHeader) bool {
	var limitUpload int
	for _, files := range fileCollection {
		limitUpload += int(len(files))
	}

	return limitUpload > config.GetInt("cpn.quiz.upload.limit.file")
}

// =>อัปโหลดต่อไฟล์ใหญ่เกิน
func isExceedPerFile(fileCollection map[string][]*multipart.FileHeader) bool {
	for _, files := range fileCollection {
		for _, file := range files {
			fmt.Println(file.Size, config.GetInt64("cpn.quiz.upload.limit.perfile"))
			if file.Size > config.GetInt64("cpn.quiz.upload.limit.perfile") {
				return true
			}
		}
	}

	return false
}

// =>อัปโหลดต่อ Request ใหญ่เกิน
func isExceedPerRequest(fileCollection map[string][]*multipart.FileHeader) bool {
	var totalSize int64
	for _, files := range fileCollection {
		for _, file := range files {
			totalSize += int64(file.Size)
		}
	}

	return totalSize > config.GetInt64("cpn.quiz.upload.limit.perrequest")
}

// =>ตรวจสอบ Ext
func isExtensionInvalid(allowExtensions []string, fileCollection map[string][]*multipart.FileHeader) (bool, []string) {
	var errors []string
	invalid := false
	for _, files := range fileCollection {
		for _, file := range files {
			ext := filepath.Ext(file.Filename)
			if !slices.Contains(allowExtensions, ext) {
				invalid = true
				errors = append(errors, fmt.Sprintf("Can't upload file '%s', Unsupported extension.", file.Filename))
			}
		}
	}

	return invalid, errors
}
