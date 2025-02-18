package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

type PubSub struct {
	Client *redis.Client
}

func NewPubSub(addr, password string, db int) *PubSub {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &PubSub{
		Client: client,
	}
}

func (p *PubSub) Publish(ctx context.Context, channel string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return p.Client.Publish(ctx, channel, data).Err()
}

func (p *PubSub) Subscribe(ctx context.Context, channel string, handler func(message string)) error {
	subscription := p.Client.Subscribe(ctx, channel)
	defer func() {
		if err := subscription.Close(); err != nil {
			log.Printf("Failed to close subscription: %v", err)
		}
	}()

	// Processing incoming messages.
	ch := subscription.Channel()
	for {
		select {
		case <-ctx.Done():
			// Gracefully exits if the context is canceled.
			return ctx.Err()
		case msg := <-ch:
			if msg == nil {
				// Channel closed or error occurred.
				return fmt.Errorf("channel closed unexpectedly")
			}
			handler(msg.Payload)
		}
	}
}
