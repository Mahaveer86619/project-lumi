package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/Mahaveer86619/ms/auth/pkg/config"
)

var DB *sql.DB

func InitDB() {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.GConfig.DBUser,
		config.GConfig.DBPassword,
		config.GConfig.DBHost,
		config.GConfig.DBPort,
		config.GConfig.DBName,
	)

	DB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Println("Successfully connected to the database!")
}
