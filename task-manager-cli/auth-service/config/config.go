package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Initialize the environment by loading the .env file
func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// GetJWTSecret reads the JWT secret from the .env file
func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set in the .env file")
	}
	return secret
}

// GetPort reads the PORT value from the .env file, with a default fallback
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("PORT is not set in the .env file, using default port 3000")
		return ":3000" // Default port if none is provided
	}
	return ":" + port
}
