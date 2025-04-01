package services

import (
	"chat-server/internals/models"
	"chat-server/internals/repository"
	"chat-server/pkg/pubsub"
)

// SaveMessage stores a message in the database and publishes it via Redis Pub/Sub.
func SaveMessage(message models.Message) error {
	// Save message to database
	if err := repository.SaveMessage(message); err != nil {
		return err
	}

	// Publish message to Pub/Sub system.
	if err := pubsub.Publish(message); err != nil {
		return err
	}

	return nil
}

// GetMessagesByRoomID retrieves all messages for a given roomID from the database.
func GetMessagesByRoomID(roomID string) ([]models.Message, error) {
	return repository.GetMessagesByRoomID(roomID)
}
