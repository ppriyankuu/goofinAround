package main

import (
	"chat-server/internal/config"
	"chat-server/internal/repository"
	"chat-server/internal/services"
	"chat-server/internal/websockets"
	"chat-server/pkg/cache"
	"chat-server/pkg/db"
	"chat-server/pkg/pubsub"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.NewConfig()
	db.Connect(cfg)

	messageRepo := repository.NewMessageRepository()
	cache := cache.NewCache(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	pubsub := pubsub.NewPubSub(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	chatService := services.NewChatService(messageRepo, cache, pubsub)
	wsServer := websockets.NewWebSocketServer(chatService)

	r := gin.Default()
	r.GET("/ws", wsServer.ServeWS)
	r.GET("/messages/:groupID", func(c *gin.Context) {
		groupID := c.Param("groupID")
		messages, err := chatService.GetMessagesByGroupID(groupID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, messages)
	})

	go wsServer.Run()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// creating a channel to listen for OS signal (e.g., SIGINT, SIGTERM)
	quit := make(chan os.Signal, 1)
	// Notifying the `quit` channel when the process receives a SIGINT (Ctrl+C) or SIGTERM signal
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// Blocking the main goroutine until a signal is received on the `quit` channel
	<-quit
	log.Println("Shutting down server...")

	// context with a timeout of 5 seconds to allow the server to gracefully shut down
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// gracefully shuts down the HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
