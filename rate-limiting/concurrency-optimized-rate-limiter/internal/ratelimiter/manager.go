package ratelimiter

import (
	"concurrency-optimized-rate-limiter/internal/ratelimiter/tokenbucket"
	"context"
	"sync"
	"time"
)

type RateLimiter struct {
	mu      sync.RWMutex
	buckets map[string]*tokenbucket.TokenBucket
	config  *Config
	metrics *Metrics
}

func New(config *Config) *RateLimiter {
	return &RateLimiter{
		buckets: make(map[string]*tokenbucket.TokenBucket),
		config:  config,
		metrics: NewMetrics(),
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, clientID string, tokens int) (bool, error) {
	start := time.Now()

	bucket := rl.getOrCreateBucket(clientID)
	allowed, err := bucket.Allow(ctx, tokens)

	rl.metrics.RecordRequest(allowed, time.Since(start))
	return allowed, err
}

func (rl *RateLimiter) getOrCreateBucket(clientID string) *tokenbucket.TokenBucket {
	rl.mu.RLock()
	bucket, exists := rl.buckets[clientID]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		defer rl.mu.Unlock()
		bucket = tokenbucket.New(rl.config.GetDefaultRate(), rl.config.GetDefaultCapacity())
		rl.buckets[clientID] = bucket
	}

	return bucket
}

func (rl *RateLimiter) UpdateClient(clientID string, rate, capacity *int) {
	rl.mu.RLock()
	bucket, exists := rl.buckets[clientID]
	rl.mu.RUnlock()

	if exists {
		bucket.Update(rate, capacity)
	}
}

func (rl *RateLimiter) Metrics() *Metrics {
	return rl.metrics
}
