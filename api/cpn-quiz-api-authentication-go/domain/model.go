package domain

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

type UseCaseResult struct {
	Result     interface{}
	Errors     []string
	Success    bool
	Message    string
	StatusCode string
	Tx         string
}
