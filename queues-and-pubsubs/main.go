package main

import (
	"log"
	"queues-and-pubsubs/controllers"
	"queues-and-pubsubs/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	utils.InitRedis()

	router := gin.Default()

	router.POST("/publish/:topic", controllers.PublishMessage)
	router.POST("/subscribe/:topic", controllers.SubscribeToTopic)

	log.Println("Server running on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
