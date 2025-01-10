package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthenticationMiddleware checks for the presence of an authorization header
func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if token != "Bearer crankyyy" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user", "mock_user")
		c.Next()
	}
}
