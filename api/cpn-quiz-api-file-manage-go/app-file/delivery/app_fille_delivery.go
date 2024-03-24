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
		response.Error = append(response.Error, fmt.Sprintf("Please upload no more than %d files.", limit))
		return c.JSON(http.StatusInternalServerError, response)
	}

	//=>ตรวจสอบขนาดไฟล์ของ Request ทั้งหมดต้องไม่เกินที่จำกัดไว้
	if isExceedPerRequest(files) {
		response.Message = constant.InternalServerError
		response.Code = constant.InternalServerErrorCode
		limitSizePerRequest := config.GetInt64("cpn.quiz.upload.limit.perrequest")
		convertSize := limitSizePerRequest / 1024 / 1024
		response.Error = append(response.Error, fmt.Sprintf("Upload all files totaling no more than %d MB.", convertSize))
		return c.JSON(http.StatusInternalServerError, response)
	}

	//=>ตรวจสอบขนาดไฟล์แต่ละไฟล์ต้องไม่เกินที่จำกัดไว้
	if isExceedPerFile(files) {
		response.Message = constant.InternalServerError
		response.Code = constant.InternalServerErrorCode
		limitSizePerFile := config.GetInt64("cpn.quiz.upload.limit.perfile")
		convertSize := limitSizePerFile / 1024 / 1024
		response.Error = append(response.Error, fmt.Sprintf("Please upload a file no larger than %d MB per file.", convertSize))
		return c.JSON(http.StatusInternalServerError, response)
	}

	result := appFile.appFileUseCase.UploadFile(fileQueryParameter, files)
	result.Tx = appFile.transId
	response = _utils.Response(result)

	return c.JSON(response.StatusCode, response)
}

func (appFile appFileDelivery) RemoveFile(c echo.Context) error {
	//=>สร้าง uuid สำหรับ Tracking request ใน Logs
	appFile.transId = uuid.New().String()

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
		response.Message = constant.BadRequest
		response.Code = constant.BadRequestCode
		response.Error = append(response.Error, fmt.Sprintf("parameter '%s' or '%s' is either required.", config.GetString("cpn.quiz.fileserver.parameter.id"), config.GetString("cpn.quiz.fileserver.parameter.list")))
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

	return c.JSON(response.StatusCode, _utils.Response(result))
}

func (appFile appFileDelivery) PreviewFile(c echo.Context) error {
	//=>สร้าง uuid สำหรับ Tracking request ใน Logs
	appFile.transId = uuid.New().String()

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

	reuslt := appFile.appFileUseCase.PreviewFile(fileQueryParameter)
	reuslt.Tx = appFile.transId
	response = _utils.Response(reuslt)

	if response.StatusCode != 200 {
		return c.JSON(response.StatusCode, response)
	}

	file := reuslt.Result.(_directoryAccessService.FileStream)
	return c.Inline(file.FullNameAccess, file.Filename)
}

func (appFile appFileDelivery) DownloadFile(c echo.Context) error {
	//=>สร้าง uuid สำหรับ Tracking request ใน Logs
	appFile.transId = uuid.New().String()

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

	reuslt := appFile.appFileUseCase.DownloadFile(fileQueryParameter)
	reuslt.Tx = appFile.transId
	response = _utils.Response(reuslt)

	if response.StatusCode != 200 {
		return c.JSON(response.StatusCode, response)
	}

	file := reuslt.Result.(_directoryAccessService.FileStream)
	return c.Attachment(file.FullNameAccess, file.Filename)
}
