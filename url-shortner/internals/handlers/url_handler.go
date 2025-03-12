package handlers

import (
	"context"
	"net/http"
	"url-shortener/internals/services"
	"url-shortener/pkg/redis"

	"github.com/gin-gonic/gin"
)

func ShortenURLHandler(c *gin.Context) {
	var input struct {
		OriginalURL string `json:"original_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := services.ShortenURL(input.OriginalURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"slug": url.Slug})
}

func GetURLBySlugHandler(c *gin.Context) {
	slug := c.Param("slug")
	url, err := services.GetURLBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	ctx := context.Background()
	cachedURL, err := redis.Client.Get(ctx, slug).Result()
	if err == nil && cachedURL != "" {
		c.Redirect(http.StatusMovedPermanently, cachedURL)
		return
	}

	err = redis.Client.Set(ctx, slug, url.Original, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusMovedPermanently, url.Original)
}
