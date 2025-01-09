package controllers

import (
	"context"
	"net/http"
	"time"
	"url-shortner/database"
	"url-shortner/helpers"
	"url-shortner/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var CTX = context.Background()

func CreateShortURL(c *gin.Context) {
	var req models.CreateShortURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if shortURL, err := helpers.RDB.Get(CTX, req.OriginalURL).Result(); err == nil {
		c.JSON(http.StatusOK, gin.H{"short_url": shortURL})
		return
	}

	shortURL := uuid.New().String()[:8] // just the first 8 characters
	url := models.URL{
		OriginalURL: req.OriginalURL,
		ShortURL:    shortURL,
	}

	if err := database.DB.Create(&url).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save URL"})
		return
	}

	helpers.RDB.Set(CTX, req.OriginalURL, shortURL, 24*time.Hour)
	helpers.RDB.Set(CTX, shortURL, req.OriginalURL, 24*time.Hour)

	c.JSON(http.StatusOK, gin.H{"short_url": shortURL})
}

func RedirectToOriginalURL(c *gin.Context) {
	shortURL := c.Param("short")

	var url models.URL
	if err := database.DB.Where("short_url = ?", shortURL).First(&url).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	// incrementing click count
	url.Clicks++
	database.DB.Save(&url)

	// save to redis (cache)
	helpers.RDB.Set(CTX, shortURL, url.OriginalURL, 24*time.Hour)

	c.Redirect(http.StatusFound, url.OriginalURL)
}

func GetAnalytics(c *gin.Context) {
	shortURL := c.Param("short")

	var url models.URL
	if err := database.DB.Where("short_url = ?", shortURL).First(&url).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.JSON(http.StatusOK, models.AnalyticsResponse{
		OriginalURL: url.OriginalURL,
		Clicks:      url.Clicks,
		CreatedAt:   url.CreatedAt.Format(time.RFC3339),
	})
}
