package config

import (
	"authentication/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := "postgresql://postgres:yourpassword@localhost:5432/url_shortener?sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	DB.AutoMigrate(&models.User{})
	log.Println("DB connected!!!")
}
