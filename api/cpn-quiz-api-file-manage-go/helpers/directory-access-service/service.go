package directoryaccess

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	config "github.com/spf13/viper"
)

type DirectoryAccessService struct {
}

type DirectoryAccess struct {
	RootPath             string
	Path                 string
	DestinationDirectory string
	Permission           AccessPermission
	FileList             []string
}

func initConfiguration() FileServerConfig {

	config := FileServerConfig{
		FileSourceParameter: config.GetString("cpn.quiz.fileserver.parameter.source"),
		FilePathParameter:   config.GetString("cpn.quiz.fileserver.parameter.path"),
		FileIdParameter:     config.GetString("cpn.quiz.fileserver.parameter.id"),
		FileListParameter:   config.GetString("cpn.quiz.fileserver.parameter.list"),
		DefaultFileSource:   config.GetString("cpn.quiz.fileserver.parameter.default"),
		Filesource:          config.GetStringMap("cpn.quiz.fileserver.filesource"),
		Permission: AccessPermission{
			File: PermissionAccessType{
				Public:  config.GetUint32("cpn.quiz.fileserver.permission.file.public"),
				Private: config.GetUint32("cpn.quiz.fileserver.permission.file.private"),
			},
			Dir: PermissionAccessType{
				Public:  config.GetUint32("cpn.quiz.fileserver.permission.dir.public"),
				Private: config.GetUint32("cpn.quiz.fileserver.permission.dir.private"),
			},
		},
	}

	return config
}

func (fs DirectoryAccessService) CreateAccessDirectory(source map[string]interface{}) (*DirectoryAccess, error) {
	configuation := initConfiguration()

	//=>ตรวจสอบหากพบ parameter filesource และเป็นค่าว่างจะคืน Error
	if val, ok := source[configuation.FileSourceParameter]; ok && val == "" {
		return nil, fmt.Errorf("parameter '%s' not empty.", configuation.FileSourceParameter)
	}

	//=>ตรวจสอบไม่พบ parameter filesource จะเช็ทเป็น default ให้
	if _, ok := source[configuation.FileSourceParameter]; !ok {
		source[configuation.FileSourceParameter] = configuation.DefaultFileSource
	}

	path := ""
	//=>ตรวจสอบ parameter path ไม่ว่าง
	if _, ok := source[configuation.FilePathParameter]; ok {
		path = source[configuation.FilePathParameter].(string)
	}

	fileId := ""
	//=>ตรวจสอบ parameter file id ไม่ว่าง
	if _, ok := source[configuation.FileIdParameter]; ok {
		fileId = source[configuation.FileIdParameter].(string)
	}

	var fileList []string
	//=>ตรวจสอบ parameter file list ไม่ว่าง
	if val, ok := source[configuation.FileListParameter]; ok && val.(string) != "" {
		fileId = ""
		fileList = strings.Split(strings.Trim(source[configuation.FileListParameter].(string), " "), ",")
	}

	filesource, found := configuation.Filesource[source[configuation.FileSourceParameter].(string)]
	//=>หากไม่ตรงกับที่ Mapping จะคืน Error
	if !found {
		return nil, fmt.Errorf("filesource '%s' dosen't exist.", configuation.FileSourceParameter)
	}

	//=>กำหนดใช้งาน Filesource
	//=>แปลงเป็น JSON Bytes
	bytes, err := json.Marshal(filesource)
	if err != nil {
		return nil, err
	}

	var activeFilesource Filesource
	//=>นำ JSON Bytes มาแปลงเป็นค่าของ Filesource
	err = json.Unmarshal(bytes, &activeFilesource)
	if err != nil {
		return nil, err
	}

	//=>ดึง root path ของโปรเจกต์
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	rootPath := protectPath(filepath.Join(pwd, strings.Trim(activeFilesource.RemotePath, " ")))
	destinationPath := protectPath(filepath.Join(rootPath, strings.Trim(path, " "), strings.Trim(fileId, " ")))

	provider := DirectoryAccess{
		DestinationDirectory: destinationPath,
		RootPath:             rootPath,
		Path:                 prettyPath(source[configuation.FilePathParameter].(string)),
		Permission:           configuation.Permission,
		FileList:             fileList,
	}

	return &provider, nil
}
