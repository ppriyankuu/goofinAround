package main

import (
	"bookstore/books/controllers"
	"bookstore/books/models"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	dsn := "host=localhost user=bookstore password=password dbname=bookstore_db port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to Database")
	}

	if err := db.AutoMigrate(&models.Book{}); err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	controllers.SetDatabase(db)

	r := gin.Default()

	r.GET("/books")
	r.POST("/books")
	r.PUT("/books/:id")
	r.DELETE("/books/:id")

	r.Run(":8080")
}
