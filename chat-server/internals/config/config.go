package config

import "os"

type Config struct {
	DBURL         string
	RedisAddr     string
	WebSocketPort int
	RESTAPIPort   int
}

func NewConfig() *Config {
	return &Config{
		DBURL:         os.Getenv("MONGODB_URI"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		WebSocketPort: 8080,
		RESTAPIPort:   8081,
	}
}
