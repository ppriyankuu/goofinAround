package services

import (
	"context"
	"url-shortener/internals/models"
	"url-shortener/pkg/db"
	"url-shortener/pkg/redis"
)

func LogVisit(urlID uint, ip, deviceType, geolocation string) error {
	analytics := &models.Analytics{
		URLID:       urlID,
		IP:          ip,
		DeviceType:  deviceType,
		Geolocation: geolocation,
	}
	result := db.DB.Create(analytics)
	if result.Error != nil {
		return result.Error
	}

	ctx := context.Background()
	err := redis.Client.RPush(ctx, "analytics_queue", analytics).Err()
	if err != nil {
		return err
	}
	return nil
}
