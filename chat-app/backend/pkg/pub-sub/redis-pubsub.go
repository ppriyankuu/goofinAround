package pubsub

import (
	"chat-app/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

type RedisPubSub struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisPubSub(client *redis.Client) *RedisPubSub {
	return &RedisPubSub{
		client: client,
		ctx:    context.Background(),
	}
}

func (rps *RedisPubSub) Publish(channel string, msg *models.Message) error {
	// marshals the message into json format
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// publishes the message to the specified Redis Channel
	err = rps.client.Publish(rps.ctx, channel, msgBytes).Err()
	if err != nil {
		return fmt.Errorf("failed to publish message to channel '%s': %w", channel, err)
	}

	log.Printf("Message publised successfully to channel '%s'", channel)
	return nil
}

func (rps *RedisPubSub) Subscribe(channel string) (*redis.PubSub, error) {
	pubsub := rps.client.Subscribe(rps.ctx, channel)

	_, err := pubsub.Receive(rps.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to channel '%s': %w", channel, err)
	}

	log.Printf("Successfully subscribed to channel '%s'", channel)
	return pubsub, nil
}

func (rps *RedisPubSub) Receive(pubsub *redis.PubSub) (*models.Message, error) {
	msg, err := pubsub.ReceiveMessage(rps.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to receive message from PubSub: %w", err)
	}

	var message models.Message
	err = json.Unmarshal([]byte(msg.Payload), &message)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message payload: %w", err)
	}

	log.Printf("Received message on channel '%s': %+v", msg.Channel, message)
	return &message, nil
}
