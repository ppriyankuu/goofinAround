package configs

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port         string
	Servers      []string
	HealthCheck  time.Duration
	ProxyTimeout time.Duration
}

func LoadConfig() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		Servers:      getEnvSlice("SERVERS", []string{"http://localhost:8081", "http://localhost:8082"}),
		HealthCheck:  getEnvDuration("HEALTH_CHECK_INTERVAL", 10*time.Second),
		ProxyTimeout: getEnvDuration("PROXY_TIMEOUT", 10*time.Second),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.Split(value, ",")
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := strconv.Atoi(value); err == nil {
			return time.Duration(duration) * time.Second
		}
	}
	return defaultValue
}
