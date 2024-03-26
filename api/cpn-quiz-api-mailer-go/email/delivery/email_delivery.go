package delivery

import (
	"cpn-quiz-api-mailer-go/constant"
	"cpn-quiz-api-mailer-go/domain"
	_utils "cpn-quiz-api-mailer-go/utils"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	config "github.com/spf13/viper"
)

func (email emailDelivery) SendEmail(c echo.Context) error {
	email.transId = uuid.New().String()
	email.emailUseCase.SetTransaction(email.transId)
	startProcess := time.Now()
	email.log.Info(email.transId, "Start :: SendEmail")

	//=>สร้างและเตรียม Struct สำหรับ Response
	response := domain.Response{}
	response.TransactionId = email.transId
	response.ResponseData = nil
	response.Error = []string{}

	allowed := []string{".zip", ".png", ".jpg", ".jpeg", ".gif"}
	//=>Validate File
	form, err := c.MultipartForm()
	if err != nil {
		response.Message = constant.BadRequest
		response.Code = constant.BadRequestCode
		response.Error = append(response.Error, err.Error())
		return c.JSON(http.StatusBadRequest, response)
	}

	fileCollection := form.File

	//=>Validate Ext
	invalidExt, extErrors := isExtensionInvalid(allowed, fileCollection)
	if invalidExt {
		response.Message = constant.BadRequest
		response.Code = constant.BadRequestCode
		response.Error = extErrors
		return c.JSON(http.StatusBadRequest, response)
	}

	//=>Validate File
	//=>ตรวจสอบไฟล์อัปโหลดต้องไม่เกินที่จำกัดไว้ไม่เกิน 5 ไฟล์
	if isExceedLimitUpload(fileCollection) {
		response.Message = constant.InternalServerError
		response.Code = constant.InternalServerErrorCode
		limit := config.GetInt32("cpn.quiz.upload.limit.file")

		err := fmt.Sprintf("Please upload no more than %d files.", limit)
		response.Error = append(response.Error, err)
		email.log.Error(email.transId, "SendEmail.Error.isExceedLimitUpload: "+err)
		email.log.Info(email.transId, "End :: SendEmail")
		return c.JSON(http.StatusInternalServerError, response)
	}

	//=>ตรวจสอบขนาดไฟล์ของ Request ทั้งหมดต้องไม่เกินที่จำกัดไว้ไม่เกิน 20MB
	if isExceedPerRequest(fileCollection) {
		response.Message = constant.InternalServerError
		response.Code = constant.InternalServerErrorCode
		limitSizePerRequest := config.GetInt64("cpn.quiz.upload.limit.perrequest")
		convertSize := limitSizePerRequest / 1024 / 1024

		err := fmt.Sprintf("Upload all files totaling no more than %d MB.", convertSize)
		response.Error = append(response.Error, err)
		email.log.Error(email.transId, "SendEmail.Error.isExceedPerRequest: "+err)
		email.log.Info(email.transId, "End :: SendEmail")
		return c.JSON(http.StatusInternalServerError, response)
	}

	//=>ตรวจสอบขนาดไฟล์แต่ละไฟล์ต้องไม่เกินที่จำกัดไว้ไม่เกิน 5MB
	if isExceedPerFile(fileCollection) {
		response.Message = constant.InternalServerError
		response.Code = constant.InternalServerErrorCode
		limitSizePerFile := config.GetInt64("cpn.quiz.upload.limit.perfile")
		convertSize := limitSizePerFile / 1024 / 1024

		err := fmt.Sprintf("Please upload a file no larger than %d MB per file.", convertSize)
		response.Error = append(response.Error, err)
		email.log.Error(email.transId, "SendEmail.Error.isExceedPerFile: "+err)
		email.log.Info(email.transId, "End :: SendEmail")
		return c.JSON(http.StatusInternalServerError, response)
	}

	//=>สกัด Email parameter
	configuation := buildEmailConfiguration()
	provider, err := buildEmailProvider(*configuation, c)
	if err != nil {
		response.Message = constant.BadRequest
		response.Code = constant.BadRequestCode

		err := fmt.Sprintf("Can't parser email provider, %s", err.Error())
		response.Error = append(response.Error, err)
		email.log.Error(email.transId, "SendEmail.Error.buildEmailProvider: "+err)
		email.log.Info(email.transId, "End :: SendEmail")
		return c.JSON(http.StatusBadRequest, response)
	}

	jwtToken := c.Request().Header.Get("Authorization")
	path := time.Now().Format("20060102")
	//=>เตรียม Parameter สำหรับ Query file
	fileQueryParameter := make(map[string]interface{})
	fileQueryParameter[config.GetString("cpn.quiz.fileserver.parameter.source")] = "email"
	fileQueryParameter[config.GetString("cpn.quiz.fileserver.parameter.path")] = path
	fileQueryParameter[config.GetString("cpn.quiz.api.mailer.email.attachmentparameter")] = fileCollection

	result := email.emailUseCase.SendEmail(provider, fileQueryParameter, jwtToken)
	result.Tx = email.transId
	response = _utils.Response(result)

	email.log.Info(email.transId, "End :: SendEmail", email.log.GetElapsedTime(startProcess))
	return c.JSON(response.StatusCode, response)
}

func (email emailDelivery) AttachmentFile(c echo.Context) error {
	//=>สร้าง uuid สำหรับ Tracking request ใน Logs
	email.transId = uuid.New().String()
	email.log.Info(email.transId, "Start :: AttachmentFile")

	token := c.Request().Header.Get("Authorization")
	path := strings.Replace(c.Request().URL.Path, strings.Replace(c.Path(), "*", "", -1), "", -1)

	fileQueryParameter := make(map[string]interface{})
	fileQueryParameter[config.GetString("cpn.quiz.fileserver.parameter.source")] = "email"
	fileQueryParameter[config.GetString("cpn.quiz.fileserver.parameter.path")] = filepath.Dir(path)
	fileQueryParameter[config.GetString("cpn.quiz.fileserver.parameter.id")] = filepath.Base(path)

	result := email.emailUseCase.AttachmentFile(fileQueryParameter, token)
	result.Tx = email.transId
	response := _utils.Response(result)

	if result.Errors != nil {
		for _, err := range result.Errors {
			email.log.Info(email.transId, "AttachmentFile.Usecase.Error: ", err)
		}
	}

	file := result.Result.(map[string]interface{})
	if response.StatusCode == 200 {
		email.log.Info(email.transId, "End :: AttachmentFile")
		return c.Blob(file["status_code"].(int), file["mime_type"].(string), file["bytes"].([]byte))
	}

	email.log.Info(email.transId, "End :: AttachmentFile")
	return c.JSON(response.StatusCode, response)
}
