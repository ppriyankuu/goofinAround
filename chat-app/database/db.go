package database

import (
	"chat-app/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(&models.ChatRoom{}, &models.Message{})
	if err != nil {
		log.Fatalf("Failed to migrate models: %v", err)
	}
	log.Println("Database connected and models migrated.")
}
