package main

import (
	"bookstore/kafka"
	"bookstore/orders/controllers"
	"bookstore/orders/models"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=bookstore password=password dbname=bookstore_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&models.Order{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	controllers.SetDatabase(db)

	kafka.InitProducer([]string{"localhost:9092"}, "orders")

	defer kafka.CloseKafka()

	r := gin.Default()

	r.POST("/orders", controllers.CreateOrder)
	r.GET("/orders", controllers.GetOrders)
	r.GET("/orders/:id", controllers.GetOrderById)

	log.Println("Orders service running on http://localhost:8082")
	if err := r.Run(":8082"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
