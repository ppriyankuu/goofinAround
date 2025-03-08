package services

import (
	"encoding/json"
	"log"
	"real-time-pipeline/internals/models"
	"real-time-pipeline/internals/utils"
)

type ConsumerService struct {
	redisClient *utils.RedisClient
	queueName   string
}

func NewConsumerService(redisClient *utils.RedisClient, queueName string) *ConsumerService {
	return &ConsumerService{redisClient: redisClient, queueName: queueName}
}

func (c *ConsumerService) StartConsuming() {
	pubsub, err := c.redisClient.Subscribe(c.queueName)
	if err != nil {
		log.Fatalf("Failed to subscribe to queue: %v", err)
	}
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		var data models.DataModel
		if err := json.Unmarshal([]byte(msg.Payload), &data); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}
		log.Printf("Recieved data: %+v", data)
		// Data Processing here :)
	}
}
