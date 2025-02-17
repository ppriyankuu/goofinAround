package main

import (
	"chat-app/internal/config"
	"chat-app/internal/handlers"
	"chat-app/internal/repositories"
	"chat-app/internal/services"
	"chat-app/internal/utils"
	"chat-app/pkg/db"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize MongoDB
	mongoClient, err := db.InitMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Initialize Redis Pub/Sub
	redisClient, err := utils.InitRedisClient(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize services and repositories
	messageRepo := repositories.NewMessageRepository(mongoClient)
	messageService := services.NewMessageService(messageRepo)
	roomService := services.NewRoomService(redisClient, messageService)

	// Initialize HTTP server
	r := gin.Default()
	chatHandler := handlers.NewChatHandler(roomService, messageService)

	// Register all routes via ChatHandler
	chatHandler.RegisterRoutes(r)

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
