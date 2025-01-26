package routes

import (
	"authentication/controllers"
	"authentication/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	protected := r.Group("/admin")
	protected.Use(middlewares.JWTAuthMiddleware(), middlewares.RBACMiddlware("admin"))
	protected.GET("/dashboard", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome Admin"})
	})
}
