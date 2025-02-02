package main

import (
	"api-testing/controllers"
	"api-testing/database"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()
	r := gin.Default()

	userController := controllers.NewUserController()

	r.POST("/users", userController.CreateUser)
	r.GET("/users/:id", userController.GetUserByID)

	r.Run(":8080")
}
