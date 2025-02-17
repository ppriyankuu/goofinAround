package config

import (
	"os"
)

type Config struct {
	Port      string
	MongoURI  string
	RedisAddr string
}

func LoadConfig() Config {
	return Config{
		Port:      getEnv("PORT", "8080"),
		MongoURI:  getEnv("MONGO_URI", "mongodb://localhost:27017/chatdb"),
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
	}
}

// Helper function to get environment variables with default values
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
