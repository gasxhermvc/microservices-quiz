package domain

import "github.com/golang-jwt/jwt"

type Token struct {
	Username   string    `json:"username"`
	UserInfo   *UserInfo `json:"userInfo"`
	IsEmployee bool      `json:"isEmployee"`
	jwt.StandardClaims
}

type UserInfo struct {
	EmpID          int     `json:"emp_id"`
	TitleSDesc     string  `json:"title_s_desc"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	EngNameFull    string  `json:"eng_name_full"`
	PosiText       string  `json:"posi_text"`
	PlansTextShort string  `json:"plans_text_short"`
	LevelCode      string  `json:"level_code"`
	DeptChangeCode string  `json:"dept_change_code"`
	DeptSap        int     `json:"dept_sap"`
	Email          string  `json:"email"`
	DepShort3      *string `json:"dept_short3"`
	DepShort4      *string `json:"dept_short4"`
}

type LogProgramStatusStore struct {
	PROGRAM_ID    int    `json:"PROGRAM_ID"`
	PROGRAM_NAME  string `json:"PROGRAM_NAME"`
	PROGRAM_TYPE  string `json:"PROGRAM_TYPE"`
	LAST_RUN_DATE string `json:"LAST_RUN_DATE"`
	STATUS_ID     int    `json:"STATUS_ID"`
	STATUS_NAME   string `json:"STATUS_NAME"`
	TOTAL         int    `json:"TOTAL"`
	TOTAL_SUCCESS int    `json:"TOTAL_SUCCESS"`
	TOTAL_ERROR   int    `json:"TOTAL_ERROR"`
}

type DataDetailStore struct {
	ID string `json:"ID"`
}

type Response struct {
	TransactionId string `json:"transactionId"`
	Message       string `json:"msg"`
	Code          string `json:"code"`
}
