package main

import (
	"task-api/database"
	"task-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()

	router := gin.Default()

	routes.TaskRoutes(router)

	router.Run(":8080")
}
