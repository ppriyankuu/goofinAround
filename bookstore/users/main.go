package main

import (
	"bookstore/users/controllers"
	"bookstore/users/models"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=bookstore password=password dbname=bookstore_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the DB: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	controllers.SetDatabase(db)

	r := gin.Default()

	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.LoginUser)
	r.GET("/users", controllers.GetUsers)

	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
