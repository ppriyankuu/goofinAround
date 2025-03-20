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
	cfg := config.NewConfig()

	db.InitDB(cfg.DBURL)
	cache.InitCache(cfg.RedisAddr)
	pubsub.InitPubSub(cfg.RedisAddr)
	repository.InitRespository()

	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) { websockets.HandleWebSocket(c) })

	r.GET("/api/messages/:room_id", func(c *gin.Context) {
		roomID := c.Param("room_id")
		messages, err := services.GetMessagesByRoomID(roomID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, messages)
	})

	go func() {
		if err := r.Run(fmt.Sprintf(":%d", cfg.WebSocketPort)); err != nil {
			log.Fatalf("Failed to start WebSocket server: %v", err)
		}
	}()

	if err := r.Run(fmt.Sprintf(":%d", cfg.RESTAPIPort)); err != nil {
		log.Fatalf("Failed to start REST API server: %v", err)
	}
}
