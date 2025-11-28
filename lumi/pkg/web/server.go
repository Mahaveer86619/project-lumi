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
	authLimiter *mid.RateLimiter
	apiLimiter  *mid.RateLimiter
)

func initSystem() {
	config.InitConfig()
	db.InitDB()

	authLimiter = mid.NewRateLimiter(10, 1*time.Minute) // 10 req in 1 min
	apiLimiter = mid.NewRateLimiter(60, 1*time.Minute)  // 10 req in 1 min
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
	// --- Services Initialization ---
	healthService := services.NewHealthService()
	avatarService := services.NewAvatarService()
	authService := services.NewAuthService(avatarService)
	wahaService := services.NewWahaService()
	userService := services.NewUserService(wahaService)

	// --- Route Groups & Middleware ---
	authGroup := e.Group("/auth")
	authGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rateLimitHandler := authLimiter.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { next(c) }))
			rateLimitHandler.ServeHTTP(c.Response(), c.Request())
			return nil
		}
	})

	apiGroup := e.Group("/api/v1")
	apiGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rateLimitHandler := apiLimiter.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { next(c) }))
			rateLimitHandler.ServeHTTP(c.Response(), c.Request())
			return nil
		}
	})

	protectedGroup := apiGroup.Group("")
	protectedGroup.Use(mid.JWTMiddleware)

	wahaGroup := protectedGroup.Group("/whatsapp")

	// Handlers
	handlers.NewHealthHandler(apiGroup, healthService)
	handlers.NewAvatarHandler(apiGroup, avatarService)
	handlers.NewAuthHandler(authGroup, authService)
	handlers.NewUserHandler(protectedGroup, userService)
	handlers.NewWahaHandler(wahaGroup, wahaService)
}
