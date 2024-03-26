package usecase

import (
	"bytes"
	"context"
	"cpn-quiz-api-mailer-go/constant"
	"cpn-quiz-api-mailer-go/database"
	"cpn-quiz-api-mailer-go/domain"
	"cpn-quiz-api-mailer-go/helpers/restful-service"
	"cpn-quiz-api-mailer-go/logger"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"path/filepath"

	// "fmt"

	"github.com/google/uuid"
	config "github.com/spf13/viper"
)

var rest *restful.Restful = restful.NewRestful(new(logger.PatternLogger), uuid.New().String())
var ctx = context.Background()

type emailUseCase struct {
	transId          string
	emailRespository domain.EmailRespository
	rdbCon           *database.RedisDatabase
	log              *logger.PatternLogger
}

func NewEmailUseCase(emailRespository domain.EmailRespository, rdbConn *database.RedisDatabase, log *logger.PatternLogger) domain.EmailUseCase {
	return &emailUseCase{
		emailRespository: emailRespository,
		rdbCon:           rdbConn,
		log:              log,
	}
}

func (email *emailUseCase) SendEmail(provider *domain.EmailProvider, param map[string]interface{}, jwtToken string) domain.UseCaseResult {
	response := domain.UseCaseResult{}
	var errors []string

	//=>ตรวจสอบการ connect redis
	err := email.rdbCon.IsConnected()
	if err != nil {
		email.log.Error(email.transId, err.Error())
	}

	//=>สร้าง Struct สำหรับเก็บลง Queue
	emailConfig := domain.EmailQueueParameter{
		From:     provider.Parameter.From,
		To:       provider.Parameter.To,
		Cc:       provider.Parameter.Cc,
		Bcc:      provider.Parameter.Bcc,
		Subject:  provider.Parameter.Subject,
		Body:     provider.Parameter.Body,
		Priority: provider.Parameter.Priority,
		IsHtml:   provider.Parameter.IsHtml,
	}

	var attachments []domain.EmailQueueAttachment
	//=>ตรวจสอบว่ามีการแนบไฟล์ หากไม่มีจะไม่ทำส่วนนี้
	if val, ok := param[config.GetString("cpn.quiz.api.mailer.email.attachmentparameter")]; ok && len(val.(map[string][]*multipart.FileHeader)) > 0 {
		//=>สร้าง endpoint upload
		endpoint := config.GetString("cpn.quiz.endpoint.filemanage")
		endpoint = fmt.Sprintf("%s/upload", endpoint)

		//=>เตรียม Request body
		buf := &bytes.Buffer{}
		write := multipart.NewWriter(buf)

		//=>สร้าง form field
		fields := []string{
			config.GetString("cpn.quiz.fileserver.parameter.source"),
			config.GetString("cpn.quiz.fileserver.parameter.path"),
		}

		//=>pack field ลง form-data
		for _, field := range fields {
			fieldWriter, err := write.CreateFormField(field)
			//=>เขียน request body ไม่สำเร็จ
			if err != nil {
				errors = append(errors, fmt.Sprintf("%s|%s", field, err.Error()))
			}
			//=>กำหนด value ให้กับ field
			fieldWriter.Write([]byte(param[field].(string)))
		}

		//=>อ่านไฟล์ทั้งหมดใน boundary
		for _, files := range param[config.GetString("cpn.quiz.api.mailer.email.attachmentparameter")].(map[string][]*multipart.FileHeader) {
			//=>อ่านไฟล์
			for _, file := range files {
				//=>เปิดไฟล์ใน stream
				src, err := file.Open()
				if err != nil {
					errors = append(errors, fmt.Sprintf("%s|%s", file.Filename, err.Error()))
				}
				defer src.Close()

				//=>เขียนไฟล์ลง Request body หากมีหลายไฟล์จะเป็นแบบ array
				fileWriter, err := write.CreateFormFile(config.GetString("cpn.quiz.api.mailer.email.attachmentparameter"), file.Filename)
				if err != nil {
					errors = append(errors, fmt.Sprintf("%s|%s", file.Filename, err.Error()))
				}

				//=>ทำการ Copy stream ไปเก็บใน Value ของ Request body
				_, err = io.Copy(fileWriter, src)
				if err != nil {
					errors = append(errors, fmt.Sprintf("%s|%s", file.Filename, err.Error()))
				}
			}
		}

		//=>ปิดการเขียน Request body
		write.Close()

		//=>ตรวจสอบหากพบ Errors จะไม่ทำรายการต่อ
		if len(errors) > 0 {
			response.Message = constant.BadRequest
			response.StatusCode = constant.BadRequestCode
			response.Errors = errors
			response.Result = nil
			return response
		}

		//=>Post to upload file
		responseBody, statusCode, err := rest.HttpPostFormDataMultiPart(endpoint, buf, write.FormDataContentType(), jwtToken)
		//=>กรณี Upload completed
		if statusCode != 200 {
			//=>Response fail
			response.Message = constant.ServiceUnavailable
			response.StatusCode = constant.ServiceUnavailableCode
			response.Errors = append(response.Errors, fmt.Sprintf("Can't upload file, status code: %d, %s", statusCode, err.Error()))
			return response
		}

		//=>กรณีเกิด Error ทางฝั่ง Mailer หรืออัปโหลดไม่สำเร็จ
		if err != nil {
			response.Message = constant.InternalServerError
			response.StatusCode = constant.InternalServerErrorCode
			response.Errors = append(response.Errors, fmt.Sprintf("Can't upload file, %s", err.Error()))
			response.Result = nil
			return response
		}

		//=>กำหนด data เพื่อสกัด Response body string
		var data map[string]interface{}
		err = json.Unmarshal(responseBody, &data)
		if err != nil {
			response.Message = constant.InternalServerError
			response.StatusCode = constant.InternalServerErrorCode
			response.Errors = append(response.Errors, fmt.Sprintf("Can't deserialize response body %s", err.Error()))
			response.Result = nil
			return response
		}

		//=>Cast responseData เป็น []interface{}
		uploads := data["responseData"].([]interface{})

		//=>ปั้น Endpoint สำหรับ Download file
		endpoint = config.GetString("cpn.quiz.endpoint.filemanage")

		for _, upload := range uploads {
			//=>Cast ค่าให้เป็น map[string]interface{}
			file := upload.(map[string]interface{})

			//=>สร้าง Struct สำหรับเก็บข้อมูล Email Provider ที่จะส่งไปเก็บใน Queue
			attachments = append(attachments, domain.EmailQueueAttachment{
				FileId:    file["file_id"].(string),
				Filename:  file["filename"].(string),
				MimeType:  file["mime_type"].(string),
				Extension: file["extension"].(string),
				Size:      file["size"].(float64),
				DownloadUrl: fmt.Sprintf("%s/download?%s=%s&%s=%s&%s=%s", endpoint,
					config.GetString("cpn.quiz.fileserver.parameter.source"),
					param[config.GetString("cpn.quiz.fileserver.parameter.source")].(string),
					config.GetString("cpn.quiz.fileserver.parameter.path"),
					param[config.GetString("cpn.quiz.fileserver.parameter.path")].(string),
					config.GetString("cpn.quiz.fileserver.parameter.id"),
					file["file_id"].(string)),
			})
		}

	}

	//=>กรณีมีการแนบไฟล์
	if len(attachments) > 0 {
		emailConfig.Attachment = attachments
	}

	//=>Cast struct ให้เป็น interface{}
	queueVal, err := json.Marshal(emailConfig)
	if err != nil {
		response.Message = constant.InternalServerError
		response.StatusCode = constant.InternalServerErrorCode
		response.Errors = append(response.Errors, fmt.Sprintf("Can't serialize email provider %s", err.Error()))
		response.Result = nil
		return response
	}

	//=>Enqueue
	err = email.rdbCon.Enqueue(config.GetString("cpn.quiz.queue.email"), string(queueVal))
	if err != nil {
		response.Message = constant.ServiceUnavailable
		response.StatusCode = constant.ServiceUnavailableCode
		response.Errors = append(response.Errors, fmt.Sprintf("Can't enqueue to redis,%s", err.Error()))
		response.Result = nil
		return response
	}

	response.Success = true
	response.Message = constant.Success
	response.StatusCode = constant.SuccessCode

	return response
}

func (email *emailUseCase) AttachmentFile(param map[string]interface{}, jwtToken string) domain.UseCaseResult {
	response := domain.UseCaseResult{}

	//=>สร้าง endpoint
	endpoint := config.GetString("cpn.quiz.endpoint.filemanage")
	endpoint = fmt.Sprintf("%s/preview?%s=%s&%s=%s&%s=%s",
		endpoint,
		config.GetString("cpn.quiz.fileserver.parameter.source"),
		param[config.GetString("cpn.quiz.fileserver.parameter.source")].(string),
		config.GetString("cpn.quiz.fileserver.parameter.path"),
		param[config.GetString("cpn.quiz.fileserver.parameter.path")].(string),
		config.GetString("cpn.quiz.fileserver.parameter.id"),
		param[config.GetString("cpn.quiz.fileserver.parameter.id")].(string))

	//=>เตรียม struct
	result := make(map[string]interface{})

	//=>Http Request
	rawData, statusCode, _ := rest.HttpGet(endpoint, nil, param[config.GetString("cpn.quiz.fileserver.parameter.id")].(string), jwtToken)

	//=>กรณีดาวน์โหลดไม่สำเร็จ
	if statusCode != 200 {
		response.Message = constant.Success
		response.StatusCode = constant.SuccessCode
		response.Errors = append(response.Errors, fmt.Sprintf("Can't download file, status code: %d", statusCode))
		return response
	}

	//=>คืน Response สำหรับแสดงผลไฟล์
	filename := param[config.GetString("cpn.quiz.fileserver.parameter.id")].(string)
	result["filename"] = filename
	result["bytes"] = rawData
	result["status_code"] = statusCode
	result["mime_type"] = mime.TypeByExtension(filepath.Ext(filename))

	response.Message = constant.Success
	response.StatusCode = constant.SuccessCode
	response.Result = result

	return response
}

func (email *emailUseCase) SetTransaction(transId string) {
	email.transId = transId
}
