package repository

import (
	"cpn-quiz-api-authentication-go/domain"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) domain.AuthRepository {
	return &authRepository{
		db: db,
	}
}
