package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Mahaveer86619/ms/auth/pkg/config"
	"github.com/Mahaveer86619/ms/auth/pkg/db"
	"github.com/Mahaveer86619/ms/auth/pkg/handlers"
	"github.com/Mahaveer86619/ms/auth/pkg/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	initSystem()

	startServer()
}

func initSystem() {
	config.InitConfig()
	db.InitDB()
}

func startServer() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())

	registerServices(e)

	serverAddress := fmt.Sprintf(":%s", config.GConfig.Port)
	if err := e.Start(serverAddress); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}
}

func registerServices(e *echo.Echo) {

	// API group
	apiGroup := e.Group("/api/v1")

	healthService := services.NewHealthService()
	handlers.NewHealthHandler(apiGroup, healthService)

}
