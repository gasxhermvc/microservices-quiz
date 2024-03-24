package service_authorize

import (
	"cpn-quiz-api-file-manage-go/logger"
	"sync"
)

type (
	CustomAuthorizeGuard struct {
		Authorization string `json:"authorization"`
		ClientId      string `json:"client_id"`
		Only          string `json:"only"`
		Mutex         sync.Mutex
		Log           *logger.PatternLogger
	}

	SecretKey struct {
		ClientId string `gorm:"column:client_id" json:"client_id"`
		ApiKey   string `gorm:"column:api_key" json:"api_key"`
	}
)
