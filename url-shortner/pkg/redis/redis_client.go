package redis

import (
	"context"
	"log"
	"url-shortener/config"

	"github.com/go-redis/redis/v8"
)

var Client *redis.Client

func InitRedis(c config.Config) error {
	client := redis.NewClient(&redis.Options{
		Addr:     c.RedisAddr,
		Password: c.RedisPassword,
		DB:       c.RedisDB,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
		return err
	}
	Client = client
	return nil
}
