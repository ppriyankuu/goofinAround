package services

import (
	"chat-server/internals/models"
	"chat-server/internals/repository"
	"chat-server/pkg/pubsub"
)

func SaveMessage(message models.Message) error {
	if err := repository.SaveMessage(message); err != nil {
		return err
	}

	if err := pubsub.Publish(message); err != nil {
		return err
	}

	return nil
}

func GetMessagesByRoomID(roomID string) ([]models.Message, error) {
	messages, err := repository.GetMessagesByRoomID(roomID)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
