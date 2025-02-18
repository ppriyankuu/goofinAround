package services

import (
	"chat-server/internal/models"
	"chat-server/internal/repository"
	"chat-server/pkg/cache"
	"chat-server/pkg/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatService struct {
	messageRepo *repository.MessageRepository
	cache       *cache.Cache
	pubsub      *pubsub.PubSub
}

func NewChatService(messageRepo *repository.MessageRepository, cache *cache.Cache, pubsub *pubsub.PubSub) *ChatService {
	return &ChatService{
		messageRepo: messageRepo,
		cache:       cache,
		pubsub:      pubsub,
	}
}

func (s *ChatService) SendMessage(userID, groupID, content string) (*models.Message, error) {
	message := &models.Message{
		ID:      primitive.NewObjectID(),
		UserID:  userID,
		GroupID: groupID,
		Content: content,
	}

	if err := s.messageRepo.Save(message); err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("message:%s", message.ID.Hex())
	s.cache.Set(context.Background(), cacheKey, message, time.Hour)
	if err := s.pubsub.Publish(context.Background(), fmt.Sprintf("group:%s", groupID), message); err != nil {
		log.Printf("Failed to publish message to group %s: %v", groupID, err)
	}

	return message, nil
}

func (s *ChatService) GetMessagesByGroupID(groupID string) ([]*models.Message, error) {
	cacheKey := fmt.Sprintf("messages:%s", groupID)
	cachedMessages, err := s.cache.Get(context.Background(), cacheKey)
	if err == nil {
		var messages []*models.Message
		if err := json.Unmarshal([]byte(cachedMessages), &messages); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached messages: %v", err)
		}
		return messages, nil
	}

	messages, err := s.messageRepo.FindByGroupID(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages from database: %v", err)
	}

	jsonMessages, err := json.Marshal(messages)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal messages for caching: %v", err)
	}

	expiration := 1 * time.Hour
	if err := s.cache.Set(context.Background(), cacheKey, string(jsonMessages), expiration); err != nil {
		log.Printf("Failed to set messages in cache: %v", err)
	}

	return messages, nil
}
