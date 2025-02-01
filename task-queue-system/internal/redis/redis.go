package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(address, password string, db int) *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			DB:       db,
		}),
	}
}

func (r *RedisClient) Set(key, value string) error {
	return r.client.Set(context.Background(), key, value, 0).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	val, err := r.client.Get(context.Background(), key).Result()
	if err != redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}
