package main

import (
	"context"
	"log"
	"url-shortener/config"
	"url-shortener/internals/handlers"
	"url-shortener/internals/middlewares"
	"url-shortener/internals/models"
	"url-shortener/pkg/db"
	"url-shortener/pkg/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := db.InitDB(c); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	if err := redis.InitRedis(c); err != nil {
		log.Fatalf("failed to initialize redis: %v", err)
	}

	db.DB.AutoMigrate(&models.URL{}, &models.Analytics{})

	go processAnalyticsQueue()

	router := gin.Default()
	router.Use(middlewares.RateLimitingMiddleware())

	router.POST("/shorten", handlers.ShortenURLHandler)
	router.GET("/r/:slug", handlers.GetURLBySlugHandler)
	router.POST("/log-visit", handlers.LogVisitHandler)

	log.Printf("Server running on :%s", c.Port)
	if err := router.Run(":" + c.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func processAnalyticsQueue() {
	ctx := context.Background()
	pubsub := redis.Client.Subscribe(ctx, "analytics_queue")
	ch := pubsub.Channel()

	for msg := range ch {
		log.Printf("Processing analytics: %s", msg.Payload)
	}
}
