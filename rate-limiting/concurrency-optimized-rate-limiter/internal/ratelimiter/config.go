package ratelimiter

import (
	"sync"
	"time"
)

type Config struct {
	mu             sync.RWMutex
	defaultRate    int
	defaultCap     int
	updateInterval time.Duration
}

func NewConfig(defaultRate, defaultCap int) *Config {
	return &Config{
		defaultRate:    defaultRate,
		defaultCap:     defaultCap,
		updateInterval: 1 * time.Second, // Default token refill interval
	}
}

func (c *Config) GetDefaultRate() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.defaultRate
}

func (c *Config) GetDefaultCapacity() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.defaultCap
}

func (c *Config) UpdateDefaults(newRate, newCap int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.defaultRate = newRate
	c.defaultCap = newCap
}

func (c *Config) GetUpdateInterval() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.updateInterval
}

func (c *Config) SetUpdateInterval(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.updateInterval = interval
}
