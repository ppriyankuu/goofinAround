package services

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheService interface {
	Set(key string, value string, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

type cacheService struct {
	client *redis.Client
}

func NewCacheService(client *redis.Client) CacheService {
	return &cacheService{client: client}
}

func (c *cacheService) Set(key string, value string, expiration time.Duration) error {
	ctx := context.Background()
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *cacheService) Get(key string) (string, error) {
	ctx := context.Background()
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found")
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (c *cacheService) Delete(key string) error {
	ctx := context.Background()
	return c.client.Del(ctx, key).Err()
}
