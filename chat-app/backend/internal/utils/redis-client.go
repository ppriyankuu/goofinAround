package utils

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// Initialises and returns a Redis client.
// It accepts a single argument: the Redis server address
func InitRedisClient(addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Printf("Failed to connect to Redis at %s: %v", addr, err)
		return nil, err
	}

	log.Printf("Successfully connected to Redis at %s", addr)
	return client, nil
}
