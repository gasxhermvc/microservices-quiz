package repository

import (
	"cpn-quiz-api-mailer-go/domain"
	_appDataService "cpn-quiz-api-mailer-go/helpers/app-data-service"

	"gorm.io/gorm"
)

type helloWorldRepository struct {
	db *gorm.DB
}

func NewHelloWorldRepository(db *gorm.DB) domain.HelloWorldRepository {
	return &helloWorldRepository{
		db: db,
	}
}

func (hello *helloWorldRepository) Hello(param domain.HelloParameter) _appDataService.QueryResult {
	response := _appDataService.QueryResult{}

	response.OutputParameter = make(map[string]interface{})
	response.Success = true
	response.AddOutputParameter("Data", []map[string]interface{}{{"Name": "john smith!"}})
	response.Total = 1

	return response
}
