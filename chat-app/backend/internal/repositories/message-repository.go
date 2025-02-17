package repositories

import (
	"chat-app/internal/models"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepository struct {
	collection *mongo.Collection
}

func NewMessageRepository(client *mongo.Client) *MessageRepository {
	if client == nil {
		panic("MongoDB client cannot be nil")
	}
	collection := client.Database("chatdb").Collection("messages")
	return &MessageRepository{collection: collection}
}

// Save saves a message to the database
func (mr *MessageRepository) Save(msg *models.Message) error {
	if msg == nil {
		return errors.New("message cannot be nil")
	}
	if msg.RoomID == "" || msg.UserID == "" || msg.Content == "" {
		return errors.New("RoomID, UserID, and Content must not be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := mr.collection.InsertOne(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to insert message into MongoDB: %w", err)
	}

	return nil
}

// GetByRoom retrieves all messages for a specific room
func (mr *MessageRepository) GetByRoom(roomID string) ([]*models.Message, error) {
	if roomID == "" {
		return nil, errors.New("roomID must not be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"room_id": roomID}
	findOptions := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})
	cursor, err := mr.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages for room %s: %w", roomID, err)
	}
	defer cursor.Close(ctx)

	var messages []*models.Message
	for cursor.Next(ctx) {
		var msg models.Message
		if err := cursor.Decode(&msg); err != nil {
			return nil, fmt.Errorf("failed to decode message: %w", err)
		}
		messages = append(messages, &msg)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error while retrieving messages: %w", err)
	}

	return messages, nil
}
