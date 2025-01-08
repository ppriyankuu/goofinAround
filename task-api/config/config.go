package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var SecretKey string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	SecretKey = os.Getenv("JWT_SECRET_KEY")
	if SecretKey == "" {
		log.Fatal("JWT_SECRET_KEY is not set in the environment")
	}
}
