package delivery

import "cpn-quiz-api-file-manage-go/domain"

func (appFile *appFileDelivery) InitConfiguration() *domain.FileServerConfig {
	config := domain.FileServerConfig{}
	return &config
}
