package config

import "os"

// Config holds the configuration values for the application.
type Config struct {
	DBURL         string // Database connection URL
	RedisAddr     string // Redis server address
	WebSocketPort int    // Port for WebSocket server
	RESTAPIPort   int    // Port for REST API server
}

// NewConfig initialises and returns a Config struct with values sourced from .env
// and default port values for WebSocket and REST API servers.
func NewConfig() *Config {
	return &Config{
		DBURL:         os.Getenv("MONGODB_URI"), // Retrieves MongoDB connection URL from environment
		RedisAddr:     os.Getenv("REDIS_ADDR"),  // Retrieves Redis address from environment
		WebSocketPort: 8080,                     // Default WebSocket server port
		RESTAPIPort:   8081,                     // Default REST API server port
	}
}
