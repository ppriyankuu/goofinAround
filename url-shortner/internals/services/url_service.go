package services

import (
	"log"
	"url-shortener/internals/models"
	"url-shortener/pkg/db"

	"github.com/aidarkhanov/nanoid/v2"
)

func ShortenURL(originalURL string) (*models.URL, error) {
	url := &models.URL{Original: originalURL, Slug: generateUniqueSlug()}
	result := db.DB.Create(url)
	if result.Error != nil {
		return nil, result.Error
	}
	return url, nil
}

func GetURLBySlug(slug string) (*models.URL, error) {
	var url models.URL
	result := db.DB.Where("slug = ?", slug).First(&url)
	if result.Error != nil {
		return nil, result.Error
	}
	return &url, nil
}

func generateUniqueSlug() string {
	for {
		slug, err := nanoid.New()
		if err != nil {
			log.Println("Error generating slug:", err)
			continue
		}
		var url models.URL
		if err := db.DB.Where("slug = ?", slug).First(&url).Error; err != nil {
			return slug
		}
	}
}
