package utils

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr, pass string, db int) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return &RedisClient{client: client}
}

func (r *RedisClient) Publish(queueName, message string) error {
	return r.client.Publish(ctx, queueName, message).Err()
}

func (r *RedisClient) Subscribe(queueName string) (*redis.PubSub, error) {
	return r.client.Subscribe(ctx, queueName), nil
}
