package api

import (
	"basic-api/api/handlers"
	"basic-api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Public routes
	public := r.Group("/api/v1")
	{
		public.POST("/login", handlers.Login)
		public.GET("/health", handlers.HealthCheck)
	}

	// Authenticated routes
	protected := r.Group("/api/v1/protected")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/user", handlers.GetUser)
	}

	// Rate limited routes
	limited := r.Group("api/v1/limited")
	limited.Use(middlewares.RateLimitMiddleware())
	{
		limited.GET("/data", handlers.GetLimitedData)
	}

	return r
}
