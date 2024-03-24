package directoryaccess

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// =>เขียนไฟล์
func (access DirectoryAccess) Put(fileCollection map[string][]*multipart.FileHeader) ([]FileStream, error) {
	results := []FileStream{}

	for param, files := range fileCollection {
		for _, file := range files {
			fileId := fmt.Sprintf("%s%s", generateRenameFile(), filepath.Ext(file.Filename))
			path := access.ProtectionPublicPath(protectPath(filepath.Join(access.DestinationDirectory, fileId)))
			result := FileStream{
				FileId:         fileId,
				FileParameter:  param,
				Filename:       file.Filename,
				FullNameAccess: protectPath(filepath.Join(access.DestinationDirectory, fileId)),
				FilePath:       path,
				MimeType:       mime.TypeByExtension(filepath.Ext(file.Filename)),
				Extension:      filepath.Ext(file.Filename),
			}

			buf := bytes.NewBuffer(nil)
			src, _ := file.Open()
			defer src.Close()
			if _, err := io.Copy(buf, src); err != nil {
				result.Error = access.ProtectionPublicPath(err.Error())
			} else {
				result.Size = buf.Len()
				result.RawBytes = buf.Bytes()
				os.WriteFile(result.FullNameAccess, result.RawBytes, fs.FileMode(access.Permission.File.Public))
			}
			results = append(results, result)
		}
	}

	return results, nil
}

// =>ลบไฟล์เดียว
func (access DirectoryAccess) RemoveFile() (string, error) {
	var fileErros string
	var errors error

	if access.FileNameExist(access.DestinationDirectory) {
		err := os.Remove(protectPath(access.DestinationDirectory))
		if err != nil {
			fileErros = filepath.Base(access.DestinationDirectory)
			errors = err
		}
	}

	return fileErros, errors
}

// =>ลบหลายไฟล์
func (access DirectoryAccess) RemoveFileName(files []string) ([]string, []error) {
	var fileErros []string
	var errors []error
	for _, file := range files {
		fullAccess := protectPath(filepath.Join(access.DestinationDirectory, strings.Trim(file, " ")))
		if access.FileNameExist(fullAccess) {
			err := os.Remove(fullAccess)
			if err != nil {
				fileErros = append(fileErros, file)
				errors = append(errors, err)
			}
		} else {
			fileErros = append(fileErros, file)
			errors = append(errors, fmt.Errorf("file '%s' dosen't exists.", file))
		}
	}
	return fileErros, errors
}

// =>ลบ Folder เดียว
func (access DirectoryAccess) RemoveDirectory(filename string, recursive bool) (bool, error) {
	return true, nil
}

// =>ดึงไฟล์เดียว
func (access DirectoryAccess) Get(filename string) FileStream {
	return FileStream{}
}

// =>ดึงไฟล์เดียว
func (access DirectoryAccess) GetFile() FileStream {
	filename := filepath.Base(access.DestinationDirectory)
	path := access.ProtectionPublicPath(protectPath(access.DestinationDirectory))
	result := FileStream{
		Filename:       filename,
		FullNameAccess: access.DestinationDirectory,
		FilePath:       path,
		MimeType:       mime.TypeByExtension(filepath.Ext(filename)),
		Extension:      filepath.Ext(filename),
	}

	buf := bytes.NewBuffer(nil)
	src, _ := os.Open(access.DestinationDirectory)
	defer src.Close()
	if _, err := io.Copy(buf, src); err != nil {
		result.Error = err.Error()
		return result
	}
	result.Size = buf.Len()
	result.RawBytes = buf.Bytes()
	return result
}

// =>ดึงข้อมูลไฟล์ทั้งหมด
func (access DirectoryAccess) Files() ([]string, error) {
	file, err := os.Open(access.DestinationDirectory)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	names, err := file.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	return names, nil
}

// =>ดึงชื่อ Folder
func (access DirectoryAccess) GetDirectory() ([]string, error) {
	return nil, nil
}

// =>ดึงชื่อไฟล์ใน Folder
func (access DirectoryAccess) GetFileDirectory() ([]string, error) {
	return nil, nil
}

// =>ดึงชื่อไฟล์ทั้งหมดใน Directory ที่เป็น Sub ด้วย
func (access DirectoryAccess) GetFileAllDirectory() ([]string, error) {
	return nil, nil
}

// =>ตรวจสอบว่ามีไฟล์ใน Path
func (access DirectoryAccess) FileNameExist(filename string) bool {
	if filepath.Ext(filename) == "" {
		return false
	}

	path := protectPath(filepath.Join(access.DestinationDirectory, filename))

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

// =>ตรวจสอบว่ามีไฟล์ใน Path
func (access DirectoryAccess) FileExist() bool {
	if filepath.Ext(access.DestinationDirectory) == "" {
		return false
	}

	path := protectPath(access.DestinationDirectory)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

// =>ตรวจสอบว่ามี Directory ชื่อนี้
func (access DirectoryAccess) DirectoryNameExist(directoryName string) bool {
	path := protectPath(filepath.Join(access.DestinationDirectory, directoryName))

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return true
}

// =>ตรวจสอบว่ามี Directory Destination
func (access DirectoryAccess) DirectoryExist() bool {
	path := protectPath(access.DestinationDirectory)

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return true
}

// =>สร้าง Folder ใหม่
func (access DirectoryAccess) MakeDirectory() (bool, error) {
	path := protectPath(access.DestinationDirectory)
	if access.DirectoryNameExist(cleansingPath(access.DestinationDirectory, path)) {
		return true, nil
	}

	err := os.Mkdir(path, fs.FileMode(access.Permission.Dir.Public))

	if err != nil {
		return false, err
	}

	return true, nil
}

// =>สร้าง Folder ใหม่
func (access DirectoryAccess) MakeNameDirectory(directoryName string) (bool, error) {
	path := protectPath(filepath.Join(access.DestinationDirectory, directoryName))
	if access.DirectoryNameExist(cleansingPath(access.DestinationDirectory, path)) {
		return true, nil
	}

	err := os.Mkdir(path, fs.FileMode(access.Permission.Dir.Public))

	if err != nil {
		return false, err
	}

	return true, nil
}

// =>สร้าง Folder ใหม่
func (access DirectoryAccess) MakeAllDirectory() (bool, error) {
	path := protectPath(access.DestinationDirectory)
	if access.DirectoryNameExist(cleansingPath(access.DestinationDirectory, path)) {
		return true, nil
	}

	err := os.MkdirAll(path, fs.FileMode(access.Permission.Dir.Public))

	if err != nil {
		return false, err
	}

	return true, nil
}

// =>สร้าง Folder ใหม่
func (access DirectoryAccess) MakeNameAllDirectory(directoryName string) (bool, error) {
	path := protectPath(filepath.Join(access.DestinationDirectory, directoryName))
	if access.DirectoryNameExist(cleansingPath(access.DestinationDirectory, path)) {
		return true, nil
	}

	err := os.MkdirAll(path, fs.FileMode(access.Permission.Dir.Public))

	if err != nil {
		return false, err
	}

	return true, nil
}

// =>ใช้เช็คว่ามันคือไฟล์
func (access DirectoryAccess) IsFile() (bool, error) {
	file, err := os.Open(access.DestinationDirectory)
	if err != nil {
		return false, err
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return false, err
	}

	if !fileInfo.IsDir() {
		return true, nil
	} else {
		return false, nil
	}
}

// =>ใช้เช็คว่ามันคือ directory
func (access DirectoryAccess) IsDir() (bool, error) {
	file, err := os.Open(access.DestinationDirectory)
	if err != nil {
		return false, err
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return false, err
	}

	if fileInfo.IsDir() {
		return true, nil
	} else {
		return false, nil
	}
}

// =>สำหรับ Display Response จะตัด Root path ออกให้
func (access DirectoryAccess) ProtectionPublicPath(path string) string {
	return cleansingPath(access.RootPath, path)
}
