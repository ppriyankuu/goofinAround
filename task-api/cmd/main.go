package main

import (
	"task-api/database"
	"task-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// loading the .env file
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	database.ConnectDB()

	router := gin.Default()

	routes.AuthRoutes(router)
	routes.TaskRoutes(router)

	router.Run(":8080")
}
