package service_authorize

import (
	config "github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var AuthorizeApiKeyLists []SecretKey

func LoadApiKeyLists() {
	//=>Load data from [CPM].[CPM_EXTERNAL_AUTHORIZE]
	conString := config.GetString("cpn.quiz.postgresql.connection.string")
	db, err := gorm.Open(postgres.Open(conString), &gorm.Config{})
	sql, _ := db.DB()
	defer sql.Close()
	if err != nil {
		panic("database connection failed : " + err.Error())
	}

	db.Raw("SELECT client_id,api_key FROM public.cq_key_authentication").Scan(&AuthorizeApiKeyLists)
}
