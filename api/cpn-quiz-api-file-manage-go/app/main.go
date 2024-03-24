package main

import (
	"cpn-quiz-api-file-manage-go/database"

	"cpn-quiz-api-file-manage-go/logger"
	"time"

	"fmt"
	"os"
	_ "time/tzdata"

	_appFileDelivery "cpn-quiz-api-file-manage-go/app-file/delivery"
	_appFileRepository "cpn-quiz-api-file-manage-go/app-file/repository"
	_appFileUseCase "cpn-quiz-api-file-manage-go/app-file/usecase"
	_auth "cpn-quiz-api-file-manage-go/middleware/service-authorize"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	config "github.com/spf13/viper"
	"github.com/tylerb/graceful"
)

var log = new(logger.PatternLogger).InitLogger("ALL", config.GetString("service.name"), logger.Service, logger.Database)

func init() {
	//=>Get argument of command.
	env := os.Args[1]
	log.Info("", fmt.Sprintf("Server start running on %s environment configuration", env))

	//=>Concat string for get file in config directory.
	config.SetConfigFile(fmt.Sprintf("config/%s.yml", env))
	if err := config.ReadInConfig(); err != nil {
		//=>Cannot load configuration.
		log.Info("", fmt.Sprintf("Fatal error env config (file): %s", env))
		panic(err)
	}

	//=>Load env.yml
	config.SetConfigFile((".env.yml"))
	if err := config.MergeInConfig(); err != nil {
		panic(err)
	}

	//=>Load db config.
	database.LoadConfig()
	_auth.LoadApiKeyLists()
}

func main() {
	db := database.Database{}
	db.Log = *log
	dbConn := db.GetConnectionDB()

	//=>new echo context
	e := echo.New()

	//=>set middleware
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("50M"))

	//=>domain data-access & business logic
	appFileRepository := _appFileRepository.NewAppFileRepository(dbConn)
	appFileUseCase := _appFileUseCase.NewAppFileUseCase(appFileRepository, log)

	//=>domain delivery
	_appFileDelivery.NewAppFileDelivery(e, appFileUseCase, log)

	//=>launch server & port
	e.Server.Addr = ":" + config.GetString("service.port")
	if err := graceful.ListenAndServe(e.Server, 5*time.Second); err != nil {
		panic(err)
	}
}
