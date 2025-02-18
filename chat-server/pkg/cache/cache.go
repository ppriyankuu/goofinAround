package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	Client *redis.Client
}

func NewCache(addr, password string, db int) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &Cache{
		Client: client,
	}
}

func (c *Cache) Close() error {
	if c.Client == nil {
		return nil
	}
	return c.Client.Close()
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = c.Client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		log.Printf("Failed to set key %s in cache: %v", key, err)
	}
	return err
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s not found in cache", key)
	} else if err != nil {
		log.Printf("Failed to get key %s from cache: %v", key, err)
		return "", err
	}
	return val, nil
}
