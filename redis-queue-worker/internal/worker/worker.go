package worker

import (
	"log"
	"redis-queue-worker/internal/redis"
	"time"

	Redis "github.com/go-redis/redis/v8"
)

type Worker struct {
	RedisClient *redis.RedisClient
}

func NewWorker(redisClient *redis.RedisClient) *Worker {
	return &Worker{RedisClient: redisClient}
}

func (w *Worker) Start(queue string) {
	for {
		message, err := w.RedisClient.PopFromQueue(queue)
		if err == Redis.Nil {
			time.Sleep(1 * time.Second)
			continue
		} else if err != nil {
			log.Printf("Error redis from queue: %v\n", err)
			continue
		}

		log.Printf("Processing message: %s\n", message)
		time.Sleep(2 * time.Second)
	}
}
