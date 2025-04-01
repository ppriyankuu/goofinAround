package db

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client   // Global MongoDB client instance
	ctx    context.Context // Global context for mongo operations
	once   sync.Once       // Ensures InitDB runs only once
)

// InitDB initialises the MongoDB client connection.
// It ensures that only one instance of the client is created.
// This function should be called once at application startup.
func InitDB(uri string) {
	once.Do(func() { // Prevents multiple initializations
		var err error
		ctx = context.Background()

		// Configure the client with the provided URI
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}

		// Verify connection by pinging the database.
		if err := client.Ping(ctx, nil); err != nil {
			log.Fatalf("Failed to ping MongoDB: %v", err)
		}

		log.Println("Connected to MongoDB")
	})
}

// GetClient returns the initialised MongoDB client instance.
// If the client is not initialised, it may cause a nil pointer dereference.
func GetClient() *mongo.Client {
	if client == nil {
		log.Fatal("MongoDB client is not initialized. Call InitDB first.")
	}
	return client
}

// GetCollection retrieves a MongoDB collection from the "chatdb" database.
// It ensures InitDB to be called before usage to ensure a valid client.
func GetCollection(name string) *mongo.Collection {
	return GetClient().Database("chatdb").Collection(name)
}
