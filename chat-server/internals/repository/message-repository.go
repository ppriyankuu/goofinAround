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

// InitRespository initialises the MongoDB collection for storing messages.
// This function should be called once at application startup.
func InitRespository() {
	collection = db.GetCollection("messages")
}

// SaveMessage inserts a new message into the MongoDB collection.
// It updates the messageID with the inserted documentID.
func SaveMessage(message models.Message) error {
	result, err := collection.InsertOne(context.TODO(), message)
	if err != nil {
		log.Printf("Failed to save message: %v", err)
	}

	// Type assertion to ensure the inserted ID is an ObjectID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		message.ID = oid.Hex()
	} else {
		log.Println("Inserted ID is not an ObjectID")
		return mongo.ErrInvalidIndexValue
	}

	return nil
}

// GemMessagesByRoomID retrieves all messages for a specific roomID.
func GetMessagesByRoomID(roomID string) ([]models.Message, error) {
	var messages []models.Message
	filter := bson.M{"room_id": roomID}

	// Fetching messages from the database.
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("Failed to find messages: %v", err)
		return nil, err
	}

	// Ensuring the cursor is closed on function end.
	defer cursor.Close(context.TODO())

	// Decoding retrieved messages into the slice
	if err := cursor.All(context.TODO(), &messages); err != nil {
		log.Printf("Failed to decode messages: %v", err)
		return nil, err
	}

	return messages, nil
}
