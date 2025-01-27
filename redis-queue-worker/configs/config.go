package configs

import "os"

type Config struct {
	RedisAddr string
}

func LoadConfig() *Config {
	return &Config{
		RedisAddr: os.Getenv("REDIS_ADDR"),
	}
}
