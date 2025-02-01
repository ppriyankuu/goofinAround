// cmd/worker/main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/ppriyankuu/task-queue-system/internal/config"
	"github.com/ppriyankuu/task-queue-system/internal/models"
	"github.com/ppriyankuu/task-queue-system/internal/postgres"
	"github.com/ppriyankuu/task-queue-system/internal/redis"
	"github.com/ppriyankuu/task-queue-system/internal/services"
	taskqueue "github.com/ppriyankuu/task-queue-system/pkg/task-queue"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	pgClient, err := postgres.NewPostgresClient(cfg.Postgres.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	if err := pgClient.AutoMigrate(&models.Task{}); err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}

	redisClient := redis.NewRedisClient(cfg.Redis.Address, cfg.Redis.Password, cfg.Redis.DB)

	taskService := services.NewTaskService(pgClient, redisClient)
	retryService := services.NewRetryService(pgClient, redisClient)

	go retryService.StartRetryLoop()

	taskQueue := taskqueue.NewKafkaTaskQueue(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer taskQueue.Close()

	handler := func(message []byte) error {
		var payload map[string]string
		if err := json.Unmarshal(message, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal message: %w", err)
		}

		taskID, err := strconv.Atoi(payload["id"])
		if err != nil {
			return fmt.Errorf("invalid task ID: %w", err)
		}

		task := &models.Task{
			ID:      uint(taskID),
			Payload: payload["task"],
			Status:  "pending",
		}

		if _, err := taskService.CreateTask(task.Payload); err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}

		if err := taskService.ProcessTask(task); err != nil {
			log.Printf("Failed to process task %d: %v", task.ID, err)
			if err := taskService.RetryTask(task); err != nil {
				log.Printf("Failed to retry task %d: %v", task.ID, err)
			}
		}

		return nil
	}

	if err := taskQueue.Consume(context.Background(), cfg.Kafka.Topic, handler); err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}
}
