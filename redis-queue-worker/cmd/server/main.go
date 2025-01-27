package main

import (
	"redis-queue-worker/configs"
	"redis-queue-worker/internal/api"
	"redis-queue-worker/internal/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	config := configs.LoadConfig()

	redisClient := redis.NewRedisClient(config.RedisAddr)

	r := gin.Default()
	enqueueHandler := api.NewEnqueueHandler(redisClient)
	r.POST("/enqueue", enqueueHandler.Enqueue)

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}

	// println("server shut down gracefully")
}
