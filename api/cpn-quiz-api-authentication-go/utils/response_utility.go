package utils

import (
	"cpn-quiz-api-authentication-go/domain"
	"strconv"
)

func Response(result domain.UseCaseResult) domain.Response {
	response := domain.Response{}

	code := strconv.Itoa(result.StatusCode)
	response.Code = code
	response.Message = result.Message
	if result.Error != nil {
		response.Message = result.Error.Error()
	}
	response.ResponseData = result.Result
	response.TransactionId = result.Tx

	return response
}
