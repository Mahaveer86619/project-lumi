package web

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Mahaveer86619/lumi/pkg/config"
	"github.com/Mahaveer86619/lumi/pkg/db"
	"github.com/Mahaveer86619/lumi/pkg/handlers"
	mid "github.com/Mahaveer86619/lumi/pkg/middleware"
	"github.com/Mahaveer86619/lumi/pkg/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	apiLimiter  *mid.RateLimiter
)

func initSystem() {
	config.InitConfig()
	db.InitDB()

	apiLimiter = mid.NewRateLimiter(10, 1*time.Minute)  // 10 req in 1 min
}

func StartServer() {
	initSystem()

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

	apiGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rateLimitHandler := apiLimiter.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { next(c) }))
			rateLimitHandler.ServeHTTP(c.Response(), c.Request())
			return nil
		}
	})

	// Services
	healthService := services.NewHealthService()
	wahaService := services.NewWahaService()

	// Handlers
	handlers.NewHealthHandler(apiGroup, healthService)
	handlers.NewWahaHandler(apiGroup, wahaService)
}
