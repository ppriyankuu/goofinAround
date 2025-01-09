package main

import (
	"url-shortner/controllers"
	"url-shortner/database"
	"url-shortner/helpers"

	"github.com/gin-gonic/gin"
)

func main() {

	database.ConnectDB()
	helpers.ConnectRedis()

	router := gin.Default()

	router.POST("/create", controllers.CreateShortURL)
	router.GET("/:short", controllers.RedirectToOriginalURL)
	router.GET("/analytics/:short", controllers.GetAnalytics)

	router.Run(":8080")
}
