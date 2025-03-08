package services

import (
	"encoding/json"
	"log"
	"real-time-pipeline/internals/models"
	"real-time-pipeline/internals/utils"
	"time"
)

type ProducerService struct {
	redisClient *utils.RedisClient
	queueName   string
}

func NewProducerService(redisClient *utils.RedisClient, queueName string) *ProducerService {
	return &ProducerService{redisClient: redisClient, queueName: queueName}
}

func (p *ProducerService) StartProducing() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			data := &models.DataModel{ID: 1, Value: " Sample Data"}
			message, err := json.Marshal(data)
			if err != nil {
				log.Printf("Failed to marshal data: %v", err)
				continue
			}
			if err := p.redisClient.Publish(p.queueName, string(message)); err != nil {
				log.Printf("Failed to publish message: %v", err)
			}
		}
	}
}
