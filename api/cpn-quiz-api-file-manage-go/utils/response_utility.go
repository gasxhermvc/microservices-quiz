package utils

import (
	"cpn-quiz-api-file-manage-go/domain"
	"strconv"
)

func Response(result domain.UseCaseResult) domain.Response {
	response := domain.Response{}

	code := result.StatusCode
	statusCode, _ := strconv.Atoi(code)
	response.Code = code
	response.Message = result.Message
	if result.Errors != nil {
		response.Message = result.Message
		response.Error = result.Errors
	}
	response.StatusCode = statusCode
	response.ResponseData = result.Result
	response.TransactionId = result.Tx

	return response
}
