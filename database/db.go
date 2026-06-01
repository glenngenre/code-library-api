package database

import (
	"fmt"
	"log"
	"os"

	"github.com/glenngenre/code-paste-service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	DB = db

	db.AutoMigrate(&models.Snippet{})
	db.AutoMigrate(&models.RateLimit{})
	db.AutoMigrate(&models.SnippetView{})
	log.Println("✓ Database schema migrated successfully")

	log.Println("✓ Connected to PostgreSQL successfully")
	return nil
}

func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		sqlDB.Close()
		log.Println("✓ Database connection closed")
	}
	return nil
}
