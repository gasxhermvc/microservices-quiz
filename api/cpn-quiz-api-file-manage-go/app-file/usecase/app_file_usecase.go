package usecase

import (
	"cpn-quiz-api-file-manage-go/constant"
	"cpn-quiz-api-file-manage-go/domain"
	"cpn-quiz-api-file-manage-go/logger"
	"fmt"
	"mime/multipart"

	_directoryAccessService "cpn-quiz-api-file-manage-go/helpers/directory-access-service"

	config "github.com/spf13/viper"
)

type appFileUseCase struct {
	appFileRespository domain.AppFileRespository
	log                *logger.PatternLogger
}

func NewAppFileUseCase(appFileRespository domain.AppFileRespository, log *logger.PatternLogger) domain.AppFileUseCase {
	return &appFileUseCase{
		appFileRespository: appFileRespository,
		log:                log,
	}
}

func (appFile *appFileUseCase) UploadFile(param map[string]interface{}, files map[string][]*multipart.FileHeader) domain.UseCaseResult {
	response := domain.UseCaseResult{}

	//=>สร้าง instance strcut สำหรับใช้จัดการไฟล์
	directoryAccessService := _directoryAccessService.DirectoryAccessService{}
	directoryAccess, err := directoryAccessService.CreateAccessDirectory(param)
	if err != nil {
		//=>สร้าง directory access ไม่สำเร็จ
		response.Message = constant.InternalServerError
		response.StatusCode = constant.InternalServerErrorCode
		response.Errors = append(response.Errors, err.Error())
		return response
	}

	//=>ตรวจสอบ Directory ที่เกิดจาก Query file ต้องไม่ว่าง
	if !directoryAccess.DirectoryExist() {
		//=>หากว่างจะมีการสร้าง Directory ทั้งหมดให้
		_, err := directoryAccess.MakeAllDirectory()
		if err != nil {
			response.Message = constant.InternalServerError
			response.StatusCode = constant.InternalServerErrorCode
			response.Errors = append(response.Errors, err.Error())
			return response
		}
	}

	//=>อัปโหลดไฟล์ทั้งหมดเก็บเข้าไปตามที่ระบุใน Query file
	fileResult, err := directoryAccess.Put(files)
	if err != nil {
		response.Message = constant.InternalServerError
		response.StatusCode = constant.InternalServerErrorCode
		response.Errors = append(response.Errors, err.Error())
		return response
	}

	response.Message = constant.Success
	response.StatusCode = constant.SuccessCode
	response.Result = fileResult

	return response
}

func (appFile *appFileUseCase) RemoveFile(param map[string]interface{}) domain.UseCaseResult {
	response := domain.UseCaseResult{}

	//=>สร้าง instance strcut สำหรับใช้จัดการไฟล์
	directoryAccessService := _directoryAccessService.DirectoryAccessService{}
	directoryAccess, err := directoryAccessService.CreateAccessDirectory(param)
	if err != nil {
		//=>สร้าง directory access ไม่สำเร็จ
		response.Message = constant.InternalServerError
		response.StatusCode = constant.InternalServerErrorCode
		response.Errors = append(response.Errors, directoryAccess.ProtectionPublicPath(err.Error()))
		return response
	}

	isDir, err := directoryAccess.IsDir()
	if err != nil {
		//=>ไม่สามารถระบุได้ว่าเป็น File หรือ Directory จะได้รับ Error
		response.Message = constant.InternalServerError
		response.StatusCode = constant.InternalServerErrorCode
		response.Errors = append(response.Errors, directoryAccess.ProtectionPublicPath(err.Error()))
		return response
	}

	//=>หากเป็น Directory และมีไฟล์ List จะลบไฟล์ทั้งหมดตามที่ระบุ
	//=>เนื่องจากเป็น Dynamic จะไม่อนุญาตให้ลบ Folder แต่สามารถลบได้หากเป็นการเขียน Function ขึ้นเองปกติ
	if isDir && len(directoryAccess.FileList) > 0 {
		//=>ลบหลายไฟล์
		errFiles, errMessages := directoryAccess.RemoveFileName(directoryAccess.FileList)
		if errMessages != nil {
			response.Message = constant.InternalServerError
			response.StatusCode = constant.InternalServerErrorCode
			for i, err := range errMessages {
				response.Errors = append(response.Errors, fmt.Sprintf("Filename: '%s', Error: %s", errFiles[i], directoryAccess.ProtectionPublicPath(err.Error())))
			}
			return response
		}
	} else {
		//=>ลบไฟล์เดียว
		file, err := directoryAccess.RemoveFile()
		if err != nil {
			response.Message = constant.InternalServerError
			response.StatusCode = constant.InternalServerErrorCode
			response.Errors = append(response.Errors, fmt.Sprintf("Filename: '%s', Error: %s", file, directoryAccess.ProtectionPublicPath(err.Error())))
			return response
		}
	}

	response.Message = constant.Success
	response.StatusCode = constant.SuccessCode

	return response
}

func (appFile *appFileUseCase) DownloadFile(param map[string]interface{}) domain.UseCaseResult {
	response := domain.UseCaseResult{}

	//=>สร้าง instance strcut สำหรับใช้จัดการไฟล์
	directoryAccessService := _directoryAccessService.DirectoryAccessService{}
	directoryAccess, err := directoryAccessService.CreateAccessDirectory(param)
	if err != nil {
		//=>สร้าง directory access ไม่สำเร็จ
		response.Message = constant.InternalServerError
		response.StatusCode = constant.InternalServerErrorCode
		response.Errors = append(response.Errors, err.Error())
		return response
	}

	//=>ดึงข้อมูลไฟล์ที่ตรงกับ Query file
	file := directoryAccess.GetFile()
	if file.Error != "" {
		//=>พบ Error
		response.Message = constant.NotFound
		response.StatusCode = constant.NotFoundCode
		response.Errors = append(response.Errors, fmt.Sprintf("file '%s' dosen't exist.", param[config.GetString("cpn.quiz.fileserver.parameter.id")]))
		return response
	}

	//=>กำหนด Response success
	response.Message = constant.Success
	response.StatusCode = constant.SuccessCode
	response.Result = file
	return response
}

func (appFile *appFileUseCase) PreviewFile(param map[string]interface{}) domain.UseCaseResult {
	response := domain.UseCaseResult{}

	//=>สร้าง instance strcut สำหรับใช้จัดการไฟล์
	directoryAccessService := _directoryAccessService.DirectoryAccessService{}
	directoryAccess, err := directoryAccessService.CreateAccessDirectory(param)
	if err != nil {
		//=>สร้าง directory access ไม่สำเร็จ
		response.Message = constant.InternalServerError
		response.StatusCode = constant.InternalServerErrorCode
		response.Errors = append(response.Errors, err.Error())
		return response
	}

	//=>ดึงข้อมูลไฟล์ที่ตรงกับ Query file
	file := directoryAccess.GetFile()
	if file.Error != "" {
		response.Message = constant.NotFound
		response.StatusCode = constant.NotFoundCode
		response.Errors = append(response.Errors, fmt.Sprintf("file '%s' dosen't exist.", param[config.GetString("cpn.quiz.fileserver.parameter.id")]))
		return response
	}

	//=>กำหนด Response success
	response.Message = constant.Success
	response.StatusCode = constant.SuccessCode
	response.Result = file
	return response
}
