package db

import (
	"fmt"
	"log"
	"time"

	"github.com/Mahaveer86619/lumi/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

const (
	maxRetries = 10
	retryDelay = 2 * time.Second
)

func InitDB() {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.GConfig.DBUser,
		config.GConfig.DBPassword,
		config.GConfig.DBHost,
		config.GConfig.DBPort,
		config.GConfig.DBName,
	)

	var err error

	for i := 1; i <= maxRetries; i++ {
		log.Printf("Attempt %d of %d: Connecting to database...", i, maxRetries)

		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err == nil {
			sqlDB, pingErr := DB.DB()
			if pingErr == nil {
				if err := sqlDB.Ping(); err == nil {
					log.Println("Successfully connected to the database via GORM!")
					return
				} else {
					log.Println(fmt.Errorf("ping failed: %v", err))
				}
			} else {
				err = fmt.Errorf("failed to get generic db interface: %v", pingErr)
			}
		}

		log.Printf("Failed to connect to database (attempt %d): %v", i, err)

		if i < maxRetries {
			log.Printf("Retrying in %v...", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	log.Fatalf("Failed to connect to the database after %d attempts. Last error: %v", maxRetries, err)
}
