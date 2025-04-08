package mongodb

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewClient creates a new MongoDB client
func NewClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get MongoDB connection string from environment variable or use default
	mongoURI := "mongodb://localhost:27017"
	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		mongoURI = uri
	}

	// Create MongoDB client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")
	return client, nil
}
