package main

import (
	"context"
	"fmt"
	"log"

	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	redisAddr := os.Getenv("REDIS_ADDRESS")
	if redisAddr == "" {
		log.Fatalf("REDIS_ADDR not set in .env")
	}

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0, // Assuming default DB
	})

	subscriber := rdb.Subscribe(ctx, "random") // Subscribe to the "random" topic
	defer subscriber.Close()

	fmt.Println("Subscribed to topic: random")

	ch := subscriber.Channel()
	for msg := range ch {
		fmt.Printf("Received message: %s\n", msg.Payload)
	}
}
