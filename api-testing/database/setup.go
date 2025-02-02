package database

import (
	"api-testing/interfaces"
	"api-testing/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() interfaces.DBInterface {
	var err error
	dsn := "host=localhost user=postgres password=yourpassword dbname=myapi port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	DB.AutoMigrate(&models.User{})
	return DB
}
