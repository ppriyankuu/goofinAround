package main

import (
	"basic-api/database"
	"basic-api/models"
	"basic-api/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()

	err := database.DB.AutoMigrate(&models.Post{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	r := gin.Default()

	r.GET("/posts", routes.GetPosts)
	r.GET("/posts/:id", routes.GetPost)
	r.POST("/posts", routes.CreatePost)
	r.PUT("/posts/:id", routes.UpdatePost)
	r.DELETE("/posts/:id", routes.DeletePost)

	r.Run(":8080")
}
