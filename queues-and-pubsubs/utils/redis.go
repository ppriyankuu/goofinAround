package utils

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	RDB *redis.Client
	CTX = context.Background()
)

func InitRedis() {
	redisAddress := os.Getenv("REDIS_ADDRESS")
	if redisAddress == "" {
		log.Fatalf("Redis address not set in .env")
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "",
		DB:       0,
	})

	_, err := RDB.Ping(CTX).Result()
	if err != nil {
		log.Fatalf("Failed to connect to redis %v", err)
	}
	log.Println("Connected to redis.")
}
