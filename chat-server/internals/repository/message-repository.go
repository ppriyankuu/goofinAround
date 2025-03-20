package repository

import (
	"chat-server/internals/models"
	"chat-server/pkg/db"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection

func InitRespository() {
	collection = db.GetCollection("messages")
}

func SaveMessage(message models.Message) error {
	result, err := collection.InsertOne(context.TODO(), message)
	if err != nil {
		log.Printf("Failed to save message: %v", err)
	}

	message.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func GetMessagesByRoomID(roomID string) ([]models.Message, error) {
	var messages []models.Message
	filter := bson.M{"room_id": roomID}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("Failed to find messages: %v", err)
		return nil, err
	}

	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &messages); err != nil {
		log.Printf("Failed to decode messages: %v", err)
		return nil, err
	}

	return messages, nil
}
