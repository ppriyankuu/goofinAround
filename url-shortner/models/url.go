package models

import "time"

type URL struct {
	ID          uint      `gorm:"primaryKey"`
	OriginalURL string    `gorm:"not null"`
	ShortURL    string    `gorm:"unique;not null"`
	Clicks      int       `gorm:"defaul:0"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

type CreateShortURLRequest struct {
	OriginalURL string `json:"original_url" binding:"required"`
}

type AnalyticsResponse struct {
	OriginalURL string `json:"original_url"`
	Clicks      int    `json:"clicks"`
	CreatedAt   string `json:"created_at"`
}
