package api

import (
	"net/http"
	"redis-queue-worker/internal/redis"

	"github.com/gin-gonic/gin"
)

type EnqueueHandler struct {
	RedisClient *redis.RedisClient
}

func NewEnqueueHandler(redisClient *redis.RedisClient) *EnqueueHandler {
	return &EnqueueHandler{RedisClient: redisClient}
}

func (h *EnqueueHandler) Enqueue(c *gin.Context) {
	var request struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.RedisClient.PushToQueue("task_queue", request.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue message."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message enqueued"})
}
