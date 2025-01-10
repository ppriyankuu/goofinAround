package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger logs all incoming requests with method, path, and processing time
func RequestLoggerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		log.Printf("[%s] %s %s | %d %s",
			c.Request.Method,   // HTTP method (get, post, etc)
			c.Request.URL.Path, // Request path
			time.Since(start),  // Processing time
			c.Writer.Status(),  // Response status code
			c.ClientIP(),       // client IP address
		)
	}
}
