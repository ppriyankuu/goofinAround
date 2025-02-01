package services

import (
	"fmt"

	"github.com/ppriyankuu/task-queue-system/internal/models"
	"github.com/ppriyankuu/task-queue-system/internal/postgres"
	"github.com/ppriyankuu/task-queue-system/internal/redis"
)

type TaskService struct {
	pgClient *postgres.PostgresClient
	redis    *redis.RedisClient
}

func NewTaskService(pgClient *postgres.PostgresClient, redis *redis.RedisClient) *TaskService {
	return &TaskService{
		pgClient: pgClient,
		redis:    redis,
	}
}

func (s *TaskService) CreateTask(payload string) (*models.Task, error) {
	task := &models.Task{
		Payload: payload,
		Status:  "pending",
	}

	if err := s.pgClient.Save(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) ProcessTask(task *models.Task) error {
	// simulating task processing
	task.Result = fmt.Sprintf("Processed %s", task.Payload)
	task.Status = "completed"
	if err := s.pgClient.Save(task); err != nil {
		return err
	}
	return s.redis.Set(fmt.Sprintf("task:%d", task.ID), task.Status)
}

func (s *TaskService) RetryTask(task *models.Task) error {
	if task.Retries >= 3 {
		task.Status = "failed"
		return s.pgClient.Save(task)
	}
	task.Retries++
	task.Status = "pending"
	return s.pgClient.Save(task)
}
