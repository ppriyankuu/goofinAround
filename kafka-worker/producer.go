package main

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func produceTasks() {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "tasks",
		Balancer: &kafka.LeastBytes{},
	}

	defer writer.Close()

	for i := 1; i <= 10; i++ {
		message := kafka.Message{
			Key:   []byte(time.Now().String()),
			Value: []byte("Task " + time.Now().Format("15:04:05")),
		}

		for retry := 0; retry < 5; retry++ { // Retry logic
			err := writer.WriteMessages(context.Background(), message)
			if err != nil {
				log.Printf("Failed to write message (attempt %d): %v", retry+1, err)
				time.Sleep(2 * time.Second) // Backoff before retrying
				continue
			}
			log.Printf("Produced: %s\n", string(message.Value))
			break
		}
		time.Sleep(1 * time.Second)
	}
}
