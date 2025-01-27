package main

import (
	"log"
	"redis-queue-worker/configs"
	"redis-queue-worker/internal/redis"
	"redis-queue-worker/internal/worker"
)

func main() {
	config := configs.LoadConfig()

	redisClient := redis.NewRedisClient(config.RedisAddr)

	workerInstance := worker.NewWorker(redisClient)
	log.Println("Worker started, waiting for messages...")
	workerInstance.Start("task_queue")
}
