package repository

import (
	"cpn-quiz-api-mailer-go/domain"

	"gorm.io/gorm"
)

type emailRepository struct {
	db *gorm.DB
}

func NewEmailRepository(db *gorm.DB) domain.EmailRespository {
	return &emailRepository{
		db: db,
	}
}
