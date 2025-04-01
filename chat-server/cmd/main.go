package main

import (
	"chat-server/internals/config"
	"chat-server/internals/repository"
	"chat-server/internals/services"
	"chat-server/internals/websockets"
	"chat-server/pkg/cache"
	"chat-server/pkg/db"
	"chat-server/pkg/pubsub"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load application configuration
	cfg := config.NewConfig()

	// Initialise database, cache, and pub/sub services
	db.InitDB(cfg.DBURL)
	cache.InitCache(cfg.RedisAddr)
	pubsub.InitPubSub(cfg.RedisAddr)

	repository.InitRespository()

	// Create a new Gin router instance
	r := gin.Default()

	// WebSocket route for real-time messaging
	r.GET("/ws", func(c *gin.Context) { websockets.HandleWebSocket(c) })

	// REST API route to fetch messages for a room
	r.GET("/api/messages/:room_id", func(c *gin.Context) {
		roomID := c.Param("room_id")

		// Fetch the messages from the database
		messages, err := services.GetMessagesByRoomID(roomID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return fetched messages.
		c.JSON(http.StatusOK, messages)
	})

	// Start WebSocket server in a seperate goroutine
	go func() {
		if err := r.Run(fmt.Sprintf(":%d", cfg.WebSocketPort)); err != nil {
			log.Fatalf("Failed to start WebSocket server: %v", err)
		}
	}()

	// Start REST API server
	if err := r.Run(fmt.Sprintf(":%d", cfg.RESTAPIPort)); err != nil {
		log.Fatalf("Failed to start REST API server: %v", err)
	}
}
