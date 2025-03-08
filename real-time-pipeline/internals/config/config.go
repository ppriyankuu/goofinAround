package config

import (
	"os"
	"strconv"
)

type Config struct {
	RedisAddr string
	RedisPass string
	RedisDB   int
	QueueName string
}

func NewConfig() (*Config, error) {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPass := os.Getenv("REDIS_PASS")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	queueName := os.Getenv("QUEUE_NAME")

	return &Config{
		RedisAddr: redisAddr,
		RedisPass: redisPass,
		RedisDB:   redisDB,
		QueueName: queueName,
	}, nil
}
