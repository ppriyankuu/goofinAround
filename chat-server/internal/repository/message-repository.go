package repository

import (
	"chat-server/internal/models"
	"chat-server/pkg/db"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageRepository struct {
	collection *mongo.Collection
}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{
		collection: db.DB.Database("chat").Collection("messages"),
	}
}

func (r *MessageRepository) Save(message *models.Message) error {
	message.Timestamp = time.Now().Unix()
	_, err := r.collection.InsertOne(context.Background(), message)
	return err
}

func (r *MessageRepository) FindByGroupID(groupID string) ([]*models.Message, error) {
	var messages []*models.Message
	cursor, err := r.collection.Find(context.Background(), bson.M{"group_id": groupID})
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		var message models.Message
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	return messages, nil
}
