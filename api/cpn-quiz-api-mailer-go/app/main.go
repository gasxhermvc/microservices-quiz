package main

import (
	"cpn-quiz-api-mailer-go/database"

	"cpn-quiz-api-mailer-go/logger"
	"time"

	"fmt"
	"os"
	_ "time/tzdata"

	_emailDelivery "cpn-quiz-api-mailer-go/email/delivery"
	_emailRepository "cpn-quiz-api-mailer-go/email/repository"
	_emailUseCase "cpn-quiz-api-mailer-go/email/usecase"

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
}

func main() {
	db := database.Database{}
	db.Log = *log
	dbConn := db.GetConnectionDB()

	//=>Initial redis
	rdb := database.RedisDatabase{}
	rdb.Log = *log
	rdbConn := rdb.GetConnectionRedisDB()

	//=>new echo context
	e := echo.New()

	//=>set middleware
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("50M"))

	//=>domain data-access & business logic
	emailRepository := _emailRepository.NewEmailRepository(dbConn)
	emailUseCase := _emailUseCase.NewEmailUseCase(emailRepository, rdbConn, log)

	//=>domain delivery
	_emailDelivery.NewEmailDelivery(e, emailUseCase, log)

	//=>launch server & port
	e.Server.Addr = ":" + config.GetString("service.port")
	if err := graceful.ListenAndServe(e.Server, 5*time.Second); err != nil {
		panic(err)
	}
}
