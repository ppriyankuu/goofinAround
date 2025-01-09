package helpers

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func ConnectRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if RDB == nil {
		panic("Failed to connect to redis.")
	}
	fmt.Println("Connected to redis successfully.")
}
