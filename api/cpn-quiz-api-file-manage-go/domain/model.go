package domain

import "github.com/golang-jwt/jwt/v5"

type Response struct {
	TransactionId string      `json:"transactionId"`
	Message       string      `json:"msg"`
	Code          string      `json:"code"`
	ResponseData  interface{} `json:"responseData,omitempty"`
}

type ErrorResponse struct {
	Error []string `json:"errorMessage"`
}

type Token struct {
	Username   string    `json:"username"`
	UserInfo   *UserInfo `json:"userInfo"`
	IsEmployee bool      `json:"isEmployee"`
	IsCapital  bool      `json:"isCapital"` // reserved for กองโค
	jwt.RegisteredClaims
}

type UserInfo struct {
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Sub               string `json:"sub"`
}

//=>App struct.
type UseCaseResult struct {
	Result     interface{}
	Error      error
	Success    bool
	Message    string
	StatusCode int
}

type FileServerConfig struct {
	FileSourceParameter string                `json:"file_source_parameter"`
	FilePathParameter   string                `json:"file_path_parameter"`
	FileIdParameter     string                `json:"file_id_parameter"`
	FileListParameter   string                `json:"file_list_parameter"`
	DefaultFileSource   string                `json:"default_file_source"`
	Filesource          map[string]Filesource `json:"filesource"`
}

type Filesource struct {
	Domain     string `json:"domain"`
	RemotePath string `json:"remote_path"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}
