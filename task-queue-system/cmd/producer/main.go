package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ppriyankuu/task-queue-system/internal/config"
	taskqueue "github.com/ppriyankuu/task-queue-system/pkg/task-queue"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	taskqueue := taskqueue.NewKafkaTaskQueue(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer taskqueue.Close()

	for i := 0; i < 10; i++ {
		payload := map[string]string{
			"id":   fmt.Sprintf("%d", i),
			"task": fmt.Sprintf("Task %d", i),
			"time": time.Now().String(),
		}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Failed to marshal payload: %v", err)
			continue
		}

		if err := taskqueue.Produce(context.Background(), cfg.Kafka.Topic, jsonPayload); err != nil {
			log.Printf("Failed to send message: %v", err)
		}
	}
}
