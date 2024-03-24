package delivery

import (
	"cpn-quiz-api-file-manage-go/constant"
	"cpn-quiz-api-file-manage-go/domain"
	_directoryAccessService "cpn-quiz-api-file-manage-go/helpers/directory-access-service"
	_utils "cpn-quiz-api-file-manage-go/utils"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	config "github.com/spf13/viper"
)

func (appFile appFileDelivery) UploadFile(c echo.Context) error {
	//=>สร้าง uuid สำหรับ Tracking request ใน Logs
	appFile.transId = uuid.New().String()
	appFile.log.Info(appFile.transId, "Start :: UploadFile")

	//=>สร้างและเตรียม Struct สำหรับ Response
	response := domain.Response{}
	response.TransactionId = appFile.transId
	response.ResponseData = nil
	response.Error = []string{}

	//=>เตรียม Parameter สำหรับ Query file
	fileQueryParameter := make(map[string]interface{})
	fileQueryParameter[config.GetString("cpn.quiz.fileserver.parameter.source")] = c.Request().FormValue(config.GetString("cpn.quiz.fileserver.parameter.source"))
	fileQueryParameter[config.GetString("cpn.quiz.fileserver.parameter.path")] = c.Request().FormValue(config.GetString("cpn.quiz.fileserver.parameter.path"))

	//=>เตรียม Struct files
	form := c.Request().MultipartForm
	files := form.File

	//=>Validate File
	//=>ตรวจสอบไฟล์อัปโหลดต้องไม่เกินที่จำกัดไว้
	if isExceedLimitUpload(files) {
		response.Message = constant.InternalServerError
		response.Code = constant.InternalServerErrorCode
		limit := config.GetInt32("cpn.quiz.upload.limit.file")

		err := fmt.Sprintf("Please upload no more than %d files.", limit)
		response.Error = append(response.Error, err)
		appFile.log.Error(appFile.transId, "UploadFile.Error.isExceedLimitUpload: "+err)
		appFile.log.Info(appFile.transId, "End :: UploadFile")
		return c.JSON(http.StatusInternalServerError, response)
	}

	//=>ตรวจสอบขนาดไฟล์ของ Request ทั้งหมดต้องไม่เกินที่จำกัดไว้
	if isExceedPerRequest(files) {
		response.Message = constant.InternalServerError
		response.Code = constant.InternalServerErrorCode
		limitSizePerRequest := config.GetInt64("cpn.quiz.upload.limit.perrequest")
		convertSize := limitSizePerRequest / 1024 / 1024

		err := fmt.Sprintf("Upload all files totaling no more than %d MB.", convertSize)
		response.Error = append(response.Error, err)
		appFile.log.Error(appFile.transId, "UploadFile.Error.isExceedPerRequest: "+err)
		appFile.log.Info(appFile.transId, "End :: UploadFile")
		return c.JSON(http.StatusInternalServerError, response)
	}

	//=>ตรวจสอบขนาดไฟล์แต่ละไฟล์ต้องไม่เกินที่จำกัดไว้
	if isExceedPerFile(files) {
		response.Message = constant.InternalServerError
		response.Code = constant.InternalServerErrorCode
		limitSizePerFile := config.GetInt64("cpn.quiz.upload.limit.perfile")
		convertSize := limitSizePerFile / 1024 / 1024

		err := fmt.Sprintf("Please upload a file no larger than %d MB per file.", convertSize)
		response.Error = append(response.Error, err)
		appFile.log.Error(appFile.transId, "UploadFile.Error.isExceedPerFile: "+err)
		appFile.log.Info(appFile.transId, "End :: UploadFile")
		return c.JSON(http.StatusInternalServerError, response)
	}

	result := appFile.appFileUseCase.UploadFile(fileQueryParameter, files)
	result.Tx = appFile.transId
	response = _utils.Response(result)

	if result.Errors != nil {
		for _, err := range result.Errors {
			appFile.log.Error(appFile.transId, "UploadFile.Usecase.Error: "+err)
		}
	}

	appFile.log.Info(appFile.transId, "End :: UploadFile")
	return c.JSON(response.StatusCode, response)
}

func (appFile appFileDelivery) RemoveFile(c echo.Context) error {
	//=>สร้าง uuid สำหรับ Tracking request ใน Logs
	appFile.transId = uuid.New().String()
	appFile.log.Info(appFile.transId, "Start :: RemoveFile")

	//=>สร้างและเตรียม Struct สำหรับ Response
	response := domain.Response{}
	response.TransactionId = appFile.transId
	response.ResponseData = nil
	response.Error = []string{}

	fileId := c.Request().FormValue(config.GetString("cpn.quiz.fileserver.parameter.id"))
	fileList := c.Request().FormValue(config.GetString("cpn.quiz.fileserver.parameter.list"))
	hasFileParameter := false
	if fileId != "" {
		hasFileParameter = true
	}

	if fileList != "" {
		hasFileParameter = true
	}

	//=>Priority fileList จะสูงกว่า fileId
	if fileId != "" && fileList != "" {
		fileId = ""
	}

	if !hasFileParameter {
		//=>พารามิเตอร์ไม่ถูกต้อง
		err := fmt.Sprintf("parameter '%s' or '%s' is either required.", config.GetString("cpn.quiz.fileserver.parameter.id"), config.GetString("cpn.quiz.fileserver.parameter.list"))
		response.Message = constant.BadRequest
		response.Code = constant.BadRequestCode
		response.Error = append(response.Error, err)

		appFile.log.Info(appFile.transId, "RemoveFile.Parameter.Error: "+err)
		appFile.log.Info(appFile.transId, "End :: RemoveFile")
		return c.JSON(http.StatusBadRequest, response)
	}

	//=>เตรียม Parameter สำหรับ Query file
	fileQueryParameter := make(map[string]interface{})
	fileQueryParameter[config.GetString("cpn.quiz.fileserver.parameter.source")] = c.Request().FormValue(config.GetString("cpn.quiz.fileserver.parameter.source"))
	fileQueryParameter[config.GetString("cpn.quiz.fileserver.parameter.path")] = c.Request().FormValue(config.GetString("cpn.quiz.fileserver.parameter.path"))
	fileQueryParameter[config.GetString("cpn.quiz.fileserver.parameter.id")] = fileId
	fileQueryParameter[config.GetString("cpn.quiz.fileserver.parameter.list")] = fileList

	//=>Response
	result := appFile.appFileUseCase.RemoveFile(fileQueryParameter)
	result.Tx = appFile.transId
	response = _utils.Response(result)

	if result.Errors != nil {
		for _, err := range result.Errors {
			appFile.log.Error(appFile.transId, "RemoveFile.Usecase.Error: "+err)
		}
	}

	appFile.log.Info(appFile.transId, "End :: RemoveFile")
	return c.JSON(response.StatusCode, _utils.Response(result))
}

func (appFile appFileDelivery) PreviewFile(c echo.Context) error {
	//=>สร้าง uuid สำหรับ Tracking request ใน Logs
	appFile.transId = uuid.New().String()
	appFile.log.Info(appFile.transId, "Start :: PreviewFile")

	//=>สร้างและเตรียม Struct สำหรับ Response
	response := domain.Response{}
	response.TransactionId = appFile.transId
	response.ResponseData = nil
	response.Error = []string{}

	//=>เตรียม Parameter สำหรับ Query file
	fileQueryParameter := make(map[string]interface{})
	fileQueryParameter[config.GetString("cpn.quiz.fileserver.parameter.source")] = c.Request().FormValue(config.GetString("cpn.quiz.fileserver.parameter.source"))
	fileQueryParameter[config.GetString("cpn.quiz.fileserver.parameter.path")] = c.Request().FormValue(config.GetString("cpn.quiz.fileserver.parameter.path"))
	fileQueryParameter[config.GetString("cpn.quiz.fileserver.parameter.id")] = c.Request().FormValue(config.GetString("cpn.quiz.fileserver.parameter.id"))

	result := appFile.appFileUseCase.PreviewFile(fileQueryParameter)
	result.Tx = appFile.transId
	response = _utils.Response(result)

	if result.Errors != nil {
		for _, err := range result.Errors {
			appFile.log.Error(appFile.transId, "PreviewFile.Usecase.Error: "+err)
		}
	}

	if response.StatusCode != 200 {
		appFile.log.Info(appFile.transId, fmt.Sprintf("RemoveFile.Usecase.Error: %d", response.StatusCode))
		appFile.log.Info(appFile.transId, "End :: PreviewFile")
		return c.JSON(response.StatusCode, response)
	}

	file := result.Result.(_directoryAccessService.FileStream)
	appFile.log.Info(appFile.transId, "End :: PreviewFile")
	return c.Inline(file.FullNameAccess, file.Filename)
}

func (appFile appFileDelivery) DownloadFile(c echo.Context) error {
	//=>สร้าง uuid สำหรับ Tracking request ใน Logs
	appFile.transId = uuid.New().String()
	appFile.log.Info(appFile.transId, "Start :: DownloadFile")

	//=>สร้างและเตรียม Struct สำหรับ Response
	response := domain.Response{}
	response.TransactionId = appFile.transId
	response.ResponseData = nil
	response.Error = []string{}

	//=>เตรียม Parameter สำหรับ Query file
	fileQueryParameter := make(map[string]interface{})
	fileQueryParameter[config.GetString("cpn.quiz.fileserver.parameter.source")] = c.Request().FormValue(config.GetString("cpn.quiz.fileserver.parameter.source"))
	fileQueryParameter[config.GetString("cpn.quiz.fileserver.parameter.path")] = c.Request().FormValue(config.GetString("cpn.quiz.fileserver.parameter.path"))
	fileQueryParameter[config.GetString("cpn.quiz.fileserver.parameter.id")] = c.Request().FormValue(config.GetString("cpn.quiz.fileserver.parameter.id"))

	result := appFile.appFileUseCase.DownloadFile(fileQueryParameter)
	result.Tx = appFile.transId
	response = _utils.Response(result)

	if result.Errors != nil {
		for _, err := range result.Errors {
			appFile.log.Error(appFile.transId, "DownloadFile.Usecase.Error: "+err)
		}
	}

	if response.StatusCode != 200 {
		appFile.log.Info(appFile.transId, fmt.Sprintf("RemoveFile.Usecase.Error: %d", response.StatusCode))
		appFile.log.Info(appFile.transId, "End :: DownloadFile")
		return c.JSON(response.StatusCode, response)
	}

	file := result.Result.(_directoryAccessService.FileStream)
	appFile.log.Info(appFile.transId, "End :: DownloadFile")
	return c.Attachment(file.FullNameAccess, file.Filename)
}
