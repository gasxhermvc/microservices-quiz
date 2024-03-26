package domain

import "github.com/golang-jwt/jwt/v5"

type Response struct {
	TransactionId string      `json:"transactionId"`
	Message       string      `json:"msg"`
	Code          string      `json:"code"`
	ResponseData  interface{} `json:"responseData,omitempty"`
	StatusCode    int         `json:"-"`
	ErrorResponse
}

type ErrorResponse struct {
	Error []string `json:"errors,omitempty"`
}

type Token struct {
	Username   string      `json:"username"`
	UserInfo   *UserInfo   `json:"userInfo"`
	Permission Permissions `json:"permission"`
	jwt.RegisteredClaims
}

type UserInfo struct {
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Sub               string `json:"sub"`
}

type Permissions struct {
	Roles []string `json:"roles"`
}

//=>App struct.
type UseCaseResult struct {
	Result     interface{}
	Errors     []string
	Success    bool
	Message    string
	StatusCode string
	Tx         string
}
