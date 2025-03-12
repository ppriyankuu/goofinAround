package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"url-shortener/internals/services"

	"github.com/gin-gonic/gin"
)

func LogVisitHandler(c *gin.Context) {
	var input struct {
		URLID uint `json:"url_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	geolocation, err := getGeolocation(clientIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch geolocation"})
		return
	}

	if err := services.LogVisit(input.URLID, clientIP, userAgent, geolocation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Visit logged successfully"})
}

func getGeolocation(ip string) (string, error) {
	resp, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data struct {
		Country string  `json:"country"`
		City    string  `json:"city"`
		Region  string  `json:"regionName"`
		Lat     float64 `json:"lat"`
		Lon     float64 `json:"lon"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s, %s", data.City, data.Country), nil
}
