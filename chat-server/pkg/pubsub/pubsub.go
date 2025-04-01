package pubsub

import (
	"chat-server/internals/models"
	"context"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	rdb  *redis.Client   // Global Redis client instance
	ctx  context.Context // Global context for Redis operations
	once sync.Once       // Ensures InitPubSub runs only once
)

// InitPubSub initializes the Redis client for Pub/Sub functionality.
// It ensures that only one instance of the client is created.
// This function should be called once at application startup.
func InitPubSub(addr string) {
	once.Do(func() {
		ctx = context.Background()
		rdb = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "", // No password set
			DB:       0,  // Default database
		})

		// Pinging the Redis server to verify connection
		if _, err := rdb.Ping(ctx).Result(); err != nil {
			log.Fatalf("Failed to connect to Redis: %v", err)
		}

		log.Println("Connected to Redis for Pub/Sub")
	})
}

type PubSub struct {
	Client *redis.Client
}

// Publish sends a message to a Redis channel (room)
func Publish(message models.Message) error {
	return rdb.Publish(ctx, message.RoomID, message.Content).Err()
}

// Subscribe listens for messages on a given roomID and invokes the handler function on new messages.
func Subscribe(roomID string, handler func(string)) {
	// Subscribe to the given roomID (channel) in Redis
	pubsub := rdb.Subscribe(ctx, roomID)
	ch := pubsub.Channel()

	// Start a new goroutine to listen to messages
	go func() {
		for msg := range ch {
			// Invoke the provided handler function when a new message arrives
			handler(msg.Payload)
		}
	}()
}
