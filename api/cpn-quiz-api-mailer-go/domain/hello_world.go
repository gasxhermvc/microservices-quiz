package domain

import (
	_appDataService "cpn-quiz-api-mailer-go/helpers/app-data-service"
)

type HelloWorldUseCase interface {
	Hello(param HelloParameter) UseCaseResult
}

type HelloWorldRepository interface {
	Hello(param HelloParameter) _appDataService.QueryResult
}

// =>Binder not sensitive case
// =>Ref: https://echo.labstack.com/guide/binding/
type HelloParameter struct {
	//=>Required
	Name string `param:"name" query:"name" json:"name" validate:"required"`
}

type HelloResponse struct {
	Name string `json:"name"`
}
