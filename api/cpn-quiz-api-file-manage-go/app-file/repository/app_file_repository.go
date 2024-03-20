package repository

import (
	"cpn-quiz-api-file-manage-go/domain"

	"gorm.io/gorm"
)

type appFileRepository struct {
	db *gorm.DB
}

func NewAppFileRepository(db *gorm.DB) domain.AppFileRespository {
	return &appFileRepository{
		db: db,
	}
}
