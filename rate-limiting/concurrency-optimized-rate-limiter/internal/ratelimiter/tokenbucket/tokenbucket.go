package tokenbucket

import (
	"context"
	"time"
)

type TokenBucket struct {
	requestChan  chan *bucketRequest
	configChan   chan *bucketConfig
	refillRate   int
	capacity     int
	currentToken int
}

type bucketRequest struct {
	tokens     int
	responseCh chan bool
}

type bucketConfig struct {
	rate     *int
	capacity *int
}

func New(refillRate, capacity int) *TokenBucket {
	tb := &TokenBucket{
		requestChan:  make(chan *bucketRequest),
		configChan:   make(chan *bucketConfig),
		refillRate:   refillRate,
		capacity:     capacity,
		currentToken: capacity,
	}

	go tb.run()
	return tb
}

func (tb *TokenBucket) run() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tb.currentToken = min(tb.currentToken+tb.refillRate, tb.capacity)
		case req := <-tb.requestChan:
			if tb.currentToken >= req.tokens {
				tb.currentToken -= req.tokens
				req.responseCh <- true
			} else {
				req.responseCh <- false
			}
		case cfg := <-tb.configChan:
			if cfg.rate != nil {
				tb.refillRate = *cfg.rate
			}
			if cfg.capacity != nil {
				tb.capacity = *cfg.capacity
				tb.currentToken = min(tb.currentToken, tb.capacity)
			}
		}
	}
}

func (tb *TokenBucket) Allow(ctx context.Context, tokens int) (bool, error) {
	respChan := make(chan bool)
	req := &bucketRequest{tokens, respChan}

	select {
	case tb.requestChan <- req:
		select {
		case allowed := <-respChan:
			return allowed, nil
		case <-ctx.Done():
			return false, ctx.Err()
		}
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func (tb *TokenBucket) Update(rate, capacity *int) {
	tb.configChan <- &bucketConfig{
		rate:     rate,
		capacity: capacity,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
