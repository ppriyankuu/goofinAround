package services

import (
	"errors"
	"fmt"

	"chat-app/internal/models"
	"chat-app/internal/repositories"
)

type MessageService struct {
	repo *repositories.MessageRepository
}

func NewMessageService(repo *repositories.MessageRepository) *MessageService {
	if repo == nil {
		panic("MessageRepository cannot be nil")
	}
	return &MessageService{repo: repo}
}

// SaveMessage saves a message to the database
func (ms *MessageService) SaveMessage(msg *models.Message) error {
	if msg == nil {
		return errors.New("message cannot be nil")
	}
	if msg.RoomID == "" || msg.UserID == "" || msg.Content == "" {
		return errors.New("RoomID, UserID, and Content must not be empty")
	}

	if err := ms.repo.Save(msg); err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	return nil
}

// GetMessagesByRoom retrieves all messages for a specific room
func (ms *MessageService) GetMessagesByRoom(roomID string) ([]*models.Message, error) {
	if roomID == "" {
		return nil, errors.New("roomID must not be empty")
	}

	messages, err := ms.repo.GetByRoom(roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve messages for room %s: %w", roomID, err)
	}

	return messages, nil
}
