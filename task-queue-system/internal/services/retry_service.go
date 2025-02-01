package services

import (
	"log"
	"time"

	"github.com/ppriyankuu/task-queue-system/internal/models"
	"github.com/ppriyankuu/task-queue-system/internal/postgres"
	"github.com/ppriyankuu/task-queue-system/internal/redis"
)

type RetryService struct {
	pgClient *postgres.PostgresClient
	redis    *redis.RedisClient
}

func NewRetryService(pgClient *postgres.PostgresClient, redis *redis.RedisClient) *RetryService {
	return &RetryService{
		pgClient: pgClient,
		redis:    redis,
	}
}

func (s *RetryService) StartRetryLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		var tasks []models.Task
		if err := s.pgClient.DB.Where("status = ?", "failed").Find(&tasks).Error; err != nil {
			log.Printf("Failed to fetch failed tasks: %v", err)
			continue
		}

		for _, task := range tasks {
			if err := s.RetryTask(&task); err != nil {
				log.Printf("Failed to retry task %d: %v", task.ID, err)
			}
		}
	}
}

func (s *RetryService) RetryTask(task *models.Task) error {
	if task.Retries >= 3 {
		task.Status = "failed"
		return s.pgClient.Save(task)
	}
	task.Retries++
	task.Status = "pending"
	return s.pgClient.Save(task)
}
