package pubsub

import (
	"chat-server/internal/models"
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
	ctx context.Context
)

func InitPubSub(addr string) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // default DB
	})

	// pinging the Redis server to ensure it's up
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis for Pub/Sub")
}

type PubSub struct {
	Client *redis.Client
}

func Publish(message models.Message) error {
	return rdb.Publish(ctx, message.RoomID, message.Content).Err()
}

func Subscribe(roomID string, handler func(string)) {
	pubsub := rdb.Subscribe(ctx, roomID)
	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			handler(msg.Payload)
		}
	}()
}
