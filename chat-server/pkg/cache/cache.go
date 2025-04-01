package cache

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	rdb  *redis.Client   // Global Redis client instance
	ctx  context.Context // Global context for Redis operations
	once sync.Once       // Ensures InitCache runs only once
)

// InitCache initialises the Redis client connection
// It ensures that only one instance of the client is created.
// This function should be called once at application startup.
func InitCache(addr string) {
	once.Do(func() {
		ctx = context.Background()
		rdb = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "", // no password set
			DB:       0,  // default DB
		})

		// Pinging the Redis server to ensure it's up
		if _, err := rdb.Ping(ctx).Result(); err != nil {
			log.Fatalf("Failed to connect to Redis: %v", err)
		}

		log.Println("Connected to Redis")
	})
}

// SetMessage stores a message in Redis with a 10-minute expiration.
func SetMessage(roomID, message string) error {
	return rdb.Set(ctx, roomID, message, 10*time.Minute).Err()
}

// GetMessage retrieves a stored message from Redis.
func GetMessage(roomID string) (string, error) {
	return rdb.Get(ctx, roomID).Result()
}
