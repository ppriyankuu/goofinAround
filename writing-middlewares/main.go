package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // Allow all origins (use "*" for development)
		AllowMethods:     []string{"GET", "POST", "PUT"},                      // Allowed HTTP methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allowed headers
		ExposeHeaders:    []string{"Content-Length"},                          // Exposed headers
		AllowCredentials: true,                                                // Allow cookies
	}))

	router.Use(RateLimiterFunc())
	router.Use(RequestLoggerFunc())
	router.Use(AuthenticationMiddleware())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome!!!"})
	})

	router.GET("/secure", func(c *gin.Context) {
		user := c.MustGet("user").(string) // Get user info from context
		c.JSON(http.StatusOK, gin.H{"message": "You are authenticated!", "user": user})
	})

	router.Run(":8080")
}
