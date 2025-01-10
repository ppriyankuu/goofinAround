package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/throttled/throttled"
	"github.com/throttled/throttled/store/memstore"
)

// RateLimiter uses a third-party library to limit requests per client
func RateLimiterFunc() gin.HandlerFunc {
	store, err := memstore.New(65536)
	if err != nil {
		log.Fatalf("Failed to create rate limiter store: %v", err)
	}

	// configuring the rate limiter (1 reqs per second per client)
	quota := throttled.RateQuota{MaxRate: throttled.PerSec(1), MaxBurst: 5}
	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		log.Fatalf("Failed to create rate limiter: %v", err)
	}

	return func(c *gin.Context) {
		limited, _, _ := rateLimiter.RateLimit(c.ClientIP(), 1)

		if limited {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}
