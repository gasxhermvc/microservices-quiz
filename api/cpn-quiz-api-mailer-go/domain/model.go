package domain

type Response struct {
	TransactionId string      `json:"transactionId"`
	Message       string      `json:"msg"`
	Code          string      `json:"code"`
	ResponseData  interface{} `json:"responseData,omitempty"`
}

type ErrorResponse struct {
	Error []string `json:"errorMessage"`
}

type UseCaseResult struct {
	Result     interface{}
	Error      error
	Success    bool
	Message    string
	StatusCode int
}
