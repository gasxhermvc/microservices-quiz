package main

import (
	"cpn-quiz-api-file-manage-go/database"

	"cpn-quiz-api-file-manage-go/logger"
	"time"

	"fmt"
	"os"
	_ "time/tzdata"

	_helloWorldDelivery "cpn-quiz-api-file-manage-go/hello-world/delivery"
	_helloWorldRepository "cpn-quiz-api-file-manage-go/hello-world/repository"
	_helloWorldUseCase "cpn-quiz-api-file-manage-go/hello-world/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	config "github.com/spf13/viper"
	"github.com/tylerb/graceful"
	// _appDataService "cpn-quiz-api-file-manage-go/helpers/app-data-service"
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
	e.Use(middleware.BodyLimit("1M"))

	//=>domain data-access & business logic
	helloWorldRepository := _helloWorldRepository.NewHelloWorldRepository(dbConn)
	helloWorldUseCase := _helloWorldUseCase.NewHelloWorldUseCase(helloWorldRepository, log)

	//=>domain delivery
	_helloWorldDelivery.NewHelloWorldDelivery(e, helloWorldUseCase, log)

	//=>launch server & port
	e.Server.Addr = ":" + config.GetString("service.port")
	if err := graceful.ListenAndServe(e.Server, 5*time.Second); err != nil {
		panic(err)
	}
}
