package main

import (
	"real-time-pipeline/internals/config"
	"real-time-pipeline/internals/services"
	"real-time-pipeline/internals/utils"
)

func main() {
	cfg, _ := config.NewConfig()
	redisClient := utils.NewRedisClient(cfg.RedisAddr, cfg.RedisPass, cfg.RedisDB)
	consumerService := services.NewConsumerService(redisClient, cfg.QueueName)
	consumerService.StartConsuming()
}
