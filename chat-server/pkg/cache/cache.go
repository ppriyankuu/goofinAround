package cache

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
	ctx context.Context
)

func InitCache(addr string) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // default DB
	})

	// pinging the Redis server to ensure it's up
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
}

func SetMessage(roomID, message string) error {
	return rdb.Set(ctx, roomID, message, 10*time.Minute).Err()
}

func GetMessage(roomID string) (string, error) {
	return rdb.Get(ctx, roomID).Result()
}
