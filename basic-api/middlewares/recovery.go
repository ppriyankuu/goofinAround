package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, err interface{}) {
		logger.Error("Panic occurred",
			zap.Any("error", err),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)
		c.AbortWithStatusJSON(500, gin.H{
			"message": "Internal server error",
		})
	})
}
