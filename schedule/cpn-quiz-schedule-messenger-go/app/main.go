package main

import (
	"cpn-quiz-schedule-messenger-go/database"
	"cpn-quiz-schedule-messenger-go/logger"
	"cpn-quiz-schedule-messenger-go/task"
	"fmt"
	"os"
	"time"
	_ "time/tzdata"

	// _helloWorldDelivery "cpn-quiz-schedule-messenger-go/hello-world/delivery"
	// _helloWorldRepository "cpn-quiz-schedule-messenger-go/hello-world/repository"
	// _helloWorldUseCase "cpn-quiz-schedule-messenger-go/hello-world/usecase"

	"github.com/labstack/echo/v4"
	"github.com/robfig/cron"
	config "github.com/spf13/viper"
	"github.com/tylerb/graceful"
)

var log = new(logger.PatternLogger).InitLogger("ALL", config.GetString("service.name"), logger.Schedule, logger.Database)

func init() {
	env := os.Args[1]
	log.Info("", fmt.Sprintf("Server start running on %s environment configuration", env))

	config.SetConfigFile(fmt.Sprintf("config/%s.yml", env))
	if err := config.ReadInConfig(); err != nil {
		log.Info("", fmt.Sprintf("Fatal error env config (file): %s", env))
		panic(err)
	}

	config.SetConfigFile(".env.yml")
	if err := config.MergeInConfig(); err != nil {
		panic(err)
	}

	database.LoadConfig() // Load Config
}

func main() {
	// cronJobs := cron.New()
	// db := database.Database{}
	db.Log = *log
	dbPool := db.GetConnectionDB()

	e := echo.New()

	// helloWorldRepository := _helloWorldRepository.NewHelloWorldRepository(dbPool)
	// helloWorldUseCase := _helloWorldUseCase.NewHelloWorldUseCase(helloWorldRepository, log)
	// _helloWorldDelivery.NewHelloWorldDelivery(e, helloWorldUseCase, log)

	// task.NewHelloWorldHandler(helloWorldUseCase, log, cronJobs)

	e.Server.Addr = ":" + config.GetString("service.port")
	err := graceful.ListenAndServe(e.Server, 5*time.Second)
	if err != nil {
		panic(err)
	}
}
