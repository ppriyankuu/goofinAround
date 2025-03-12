package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"url-shortener/pkg/redis"

	"github.com/gin-gonic/gin"
)

func RateLimitingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit : %s", ip)
		ctx := context.Background()

		count, err := redis.Client.Incr(ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit error"})
			c.Abort()
			return
		}

		if count == 1 {
			err = redis.Client.Expire(ctx, key, 1*time.Minute).Err()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit error"})
				c.Abort()
				return
			}
		}

		if count > 50 {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}
