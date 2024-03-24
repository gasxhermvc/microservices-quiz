package delivery

import (
	"mime/multipart"
	"path/filepath"
	"slices"

	config "github.com/spf13/viper"
)

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

func isExtensionInvalid(allowExtensions []string, fileCollection map[string][]*multipart.FileHeader) bool {
	for _, files := range fileCollection {
		for _, file := range files {
			ext := filepath.Ext(file.Filename)
			if !slices.Contains(allowExtensions, ext) {
				return true
			}
		}
	}

	return false
}
