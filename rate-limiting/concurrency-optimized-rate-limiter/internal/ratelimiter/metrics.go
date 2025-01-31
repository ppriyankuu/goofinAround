package ratelimiter

import (
	"sync"
	"time"
)

type Metrics struct {
	mu               sync.Mutex
	totalRequests    uint64
	rejectedRequests uint64
	totalDuration    time.Duration
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) RecordRequest(allowed bool, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalRequests++
	if !allowed {
		m.rejectedRequests++
	}
	m.totalDuration += duration
}

func (m *Metrics) SnapShot() (uint64, uint64, time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.totalRequests, m.rejectedRequests, m.totalDuration
}
