package main

import (
	"cpn-quiz-schedule-messenger-go/database"
	"cpn-quiz-schedule-messenger-go/logger"
	"cpn-quiz-schedule-messenger-go/task"

	_emailMessengerRepository "cpn-quiz-schedule-messenger-go/email-messenger/repository"
	_emailMessengerUseCase "cpn-quiz-schedule-messenger-go/email-messenger/usecase"

	"fmt"
	"os"
	"time"
	_ "time/tzdata"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron"
	config "github.com/spf13/viper"
	"github.com/tylerb/graceful"
)

var log = new(logger.PatternLogger).InitLogger("ALL", config.GetString("service.name"), logger.Schedule, logger.Database)

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
	cronJobs := cron.New()
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

	emailMessengerRepository := _emailMessengerRepository.NewEmailMessengerRepository(dbConn)
	emailMessengerUseCase := _emailMessengerUseCase.NewEmailMessengerUseCase(emailMessengerRepository, rdbConn, log)

	task.NewHelloWorldHandler(emailMessengerUseCase, rdbConn, log, cronJobs)
	cronJobs.Run() //=>First run
	//=>launch server & port
	e.Server.Addr = ":" + config.GetString("service.port")
	if err := graceful.ListenAndServe(e.Server, 5*time.Second); err != nil {
		panic(err)
	}
}
