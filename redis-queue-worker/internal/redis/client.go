package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(redisAddr string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	return &RedisClient{Client: rdb}
}

func (r *RedisClient) PushToQueue(queue string, message string) error {
	return r.Client.LPush(ctx, queue, message).Err()
}

func (r *RedisClient) PopFromQueue(queue string) (string, error) {
	return r.Client.RPop(ctx, queue).Result()
}
