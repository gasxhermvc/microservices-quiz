package directoryaccess

type FileServer struct {
	FilesourceParameter string `json:"FileSourceParameter"`
	FilePathParameter   string `json:"FilePathParameter"`
	FileIdParameter     string `json:"FileIdParameter"`
	FileListParameter   string `json:"FileListParameter"`
	DefaultFilesource   string `json:"DefaultFileSource"`
	Filesource          Filesource
}

type FileServerConfig struct {
	FileSourceParameter string           `json:"file_source_parameter"`
	FilePathParameter   string           `json:"file_path_parameter"`
	FileIdParameter     string           `json:"file_id_parameter"`
	FileListParameter   string           `json:"file_list_parameter"`
	DefaultFileSource   string           `json:"default_file_source"`
	Filesource          map[string]any   `json:"filesource"`
	Permission          AccessPermission `json:"permission"`
}

type Filesource struct {
	Domain     string `json:"domain"`
	RemotePath string `json:"path"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

type FileStream struct {
	FileId         string `json:"file_id"`
	FileParameter  string `json:"file_parameter"`
	Filename       string `json:"filename"`
	FullNameAccess string `json:"-"`
	FilePath       string `json:"path"`
	MimeType       string `json:"mime_type"`
	Extension      string `json:"extension"`
	Size           int    `json:"size"`
	Error          string `json:"-"`
	RawBytes       []byte `json:"-"`
}

type AccessPermission struct {
	File PermissionAccessType `json:"file"`
	Dir  PermissionAccessType `json:"dir"`
}

type PermissionAccessType struct {
	Public  uint32 `json:"public"`
	Private uint32 `json:"private"`
}
