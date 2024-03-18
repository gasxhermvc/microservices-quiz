package vaildate

import (
	"cpn-quiz-api-authentication-go/constant"
	"cpn-quiz-api-authentication-go/domain"
	"unicode"

	vaildate "gopkg.in/go-playground/validator.v9"
)

func CheckRequest(req interface{}) (response domain.Response) {
	vaild := vaildate.New()
	if err := vaild.Struct(req); err != nil {
		for _, e := range err.(vaildate.ValidationErrors) {
			response = MapErrorResponse(e.Tag(), e.Field(), e.Param())
		}
	}

	return response
}

func MapErrorResponse(tag string, field string, param string) (response domain.Response) {
	field = string(unicode.ToLower([]rune(field)[0])) + field[1:]

	var statusCode string
	var message string

	switch tag {
	case constant.Required:
		statusCode = constant.BadRequestCode
		message = "Error : " + field + " is required."
	case constant.Min:
		statusCode = constant.BadRequestCode
		message = "Error : " + field + " must be greater than " + param
	case constant.Max:
		statusCode = constant.BadRequestCode
		message = "Error : " + field + " is not over " + param
	}

	response.Code = statusCode
	response.Message = message
	return
}
