package repository

import (
	"cpn-quiz-schedule-messenger-go/domain"

	"gorm.io/gorm"
)

type emailMessengerRepository struct {
	db *gorm.DB
}

func NewEmailMessengerRepository(db *gorm.DB) domain.EmailMessengerRepository {
	return &emailMessengerRepository{
		db: db,
	}
}
