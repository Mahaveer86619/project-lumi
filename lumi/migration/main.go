package main

import (
	"github.com/Mahaveer86619/lumi/pkg/config"
	"github.com/Mahaveer86619/lumi/pkg/db"
	"github.com/Mahaveer86619/lumi/pkg/models"
	"github.com/labstack/gommon/log"
)

var version string = "dev"

func main() {
	log.Printf("Starting migration, version %s...", version)

	config.InitConfig()
	db.InitDB()

	tables := []interface{}{
		&models.UserProfile{},
		&models.WhatsAppSession{},
	}

	log.Info("Running AutoMigrate...")
	if err := db.DB.AutoMigrate(tables...); err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Info("DB Migration completed successfully!")
}
