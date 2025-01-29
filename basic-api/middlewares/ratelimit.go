package middlewares

import (
	"basic-api/config"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var limiter = NewIPRateLimiter(rate.Limit(config.RateLimit()), 1)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := limiter.GetLimiter(c.ClientIP())
		if !limiter.Allow() {
			c.AbortWithStatusJSON(429, gin.H{"error": "Too many request. Slow down!"})
			return
		}
		c.Next()
	}
}

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}
