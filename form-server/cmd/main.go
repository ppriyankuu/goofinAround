package main

import (
	"fmt"
	"form-server/internals/config"
	"form-server/internals/handlers"
	"form-server/internals/models"
	"form-server/internals/repositories"
	"form-server/internals/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	db.AutoMigrate(&models.FormResponse{})

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})

	formRepo := repositories.NewFormRepo(db)
	formService := services.NewFormService(formRepo)
	cacheService := services.NewCacheService(redisClient)

	formHandler := handlers.NewFormHandler(formService, cacheService)

	router := gin.Default()
	SetupRoutes(router, formHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
