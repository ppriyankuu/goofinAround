package controllers

import (
	"task-manager/auth-service/models"
	"task-manager/auth-service/utils"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if err := models.CreateUser(&user); err != nil {
		c.JSON(500, gin.H{"error": "Could not create user"})
		return
	}

	c.JSON(201, gin.H{"message": "User registered successfully!"})
}

func Login(c *gin.Context) {
	var credentials models.Credentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	user, err := models.GetUserByEmail(credentials.Email)
	if err != nil {
		c.JSON(401, gin.H{"error": "User not found!"})
		return
	}

	if !utils.CheckPasswordHash(credentials.Password, user.Password) {
		c.JSON(401, gin.H{"error": "Invalid password"})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not generate token"})
		return
	}
	c.JSON(200, gin.H{"token": token})
}
