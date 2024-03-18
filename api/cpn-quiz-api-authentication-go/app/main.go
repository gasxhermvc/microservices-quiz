package main

import (
	// "cpn-quiz-api-authentication-go/database"

	"cpn-quiz-api-authentication-go/logger"
	"net/http"
	"time"

	"fmt"
	"os"
	_ "time/tzdata"

	_helloWorldDelivery "cpn-quiz-api-authentication-go/hello-world/delivery"
	_helloWorldRepository "cpn-quiz-api-authentication-go/hello-world/repository"
	_helloWorldUseCase "cpn-quiz-api-authentication-go/hello-world/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	config "github.com/spf13/viper"
	"github.com/tylerb/graceful"
	// _appDataService "cpn-quiz-api-authentication-go/helpers/app-data-service"
)

var log = new(logger.PatternLogger).InitLogger("ALL", config.GetString("service.name"), logger.Service, logger.Database)

func init() {
	env := os.Args[1]
	log.Info("", fmt.Sprintf("Server start running on %s environment configuration", env))

	config.SetConfigFile(fmt.Sprintf("config/%s.yml", env))
	if err := config.ReadInConfig(); err != nil {
		log.Info("", fmt.Sprintf("Fatal error env config (file): %s", env))
		panic(err)
	}

	config.SetConfigFile((".env.yml"))
	if err := config.MergeInConfig(); err != nil {
		panic(err)
	}

	log.Info("", config.GetString("service.port"))

	//=>ต้อง Connected VPN ด้วยถึงจะเรียก LoadConfig ได้
	// database.LoadConfig()
}

func main() {
	// db := database.Database{}
	// db.Log = *log
	// dbConn := db.GetConnectionDB()

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("150M"))

	//=>หามี Path สำหรับแสดง Statics Files ให้เปิดใช้งาน
	// e.Static(config.GetString("service.endpoint")+"/static", "cpm/static")

	//=>First Route
	eg := e.Group(config.GetString("service.endpoint"))
	eg.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World")
	})

	helloWorldRepository := _helloWorldRepository.NewHelloWorldRepository(nil)
	helloWorldUseCase := _helloWorldUseCase.NewHelloWorldUseCase(helloWorldRepository, log)

	_helloWorldDelivery.NewHelloWorldDelivery(e, helloWorldUseCase, log)

	e.Server.Addr = ":" + config.GetString("service.port")
	if err := graceful.ListenAndServe(e.Server, 5*time.Second); err != nil {
		panic(err)
	}

}
