package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string

	// DB config
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// jwt
	JWTSecret string

	// Gemini config
	GeminiAPIKey string

	// ms
	WahaServiceURL string
	WahaAPIKey     string

	// default session
	WahaSessionName string
}

var GConfig *Config

func InitConfig() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error: %v", err)
		log.Println("No .env file found, relying on system environment variables")
	}

	GConfig = &Config{
		Port: getEnv("PORT"),

		// DB config
		DBHost:     getEnv("DB_HOST"),
		DBPort:     getEnv("DB_PORT"),
		DBUser:     getEnv("DB_USER"),
		DBPassword: getEnv("DB_PASSWORD"),
		DBName:     getEnv("DB_NAME"),

		// jwt
		JWTSecret: getEnv("JWT_SECRET"),

		// Gemini config
		GeminiAPIKey: getEnv("GEMINI_API_KEY"),

		// ms
		WahaServiceURL: getEnv("WAHA_SERVICE_URL"),
		WahaAPIKey:     getEnv("WAHA_API_KEY"),

		// session
		WahaSessionName: getEnv("WAHA_SESSION_NAME"),
	}
}

func getEnv(key string, defaultVal ...string) string {
	val := os.Getenv(key)
	if val == "" {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		log.Fatalf("Key %s not found in .env file", key)
	}

	return val
}
