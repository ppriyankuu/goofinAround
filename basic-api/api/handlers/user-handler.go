package handlers

import (
	"basic-api/config"
	"basic-api/models"
	"basic-api/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Login(c *gin.Context) {
	var creds models.UserCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Authentication logic
	user, err := services.AuthenticateUser(creds.Email, creds.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// Creating JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.JWTSecret()))
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating token"})
		return
	}
	c.JSON(200, gin.H{"token": tokenString})
}

func GetUser(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	user, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	c.JSON(200, user)
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

func GetLimitedData(c *gin.Context) {
	c.JSON(200, gin.H{"message": "This is rate limited data"})
}
