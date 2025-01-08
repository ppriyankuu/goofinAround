package middlewares

import (
	"task-api/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime)
		utils.Logger.WithFields(logrus.Fields{
			"status":      c.Writer.Status(),
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"duration_ms": duration.Milliseconds(),
		}).Info("HTTP Request")
	}
}
