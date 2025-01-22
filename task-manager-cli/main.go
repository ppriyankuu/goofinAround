package main

import (
	"task-manager/auth-service/config"
	"task-manager/auth-service/database"
	"task-manager/auth-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()
	database.InitDB()

	r := gin.Default()

	// Middleware
	// r.Use(middlewares.AuthMiddleware())

	// Routes
	routes.AuthRoutes(r)

	r.Run(config.GetPort())
}
