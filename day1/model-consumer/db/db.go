package db

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client    *mongo.Client
	initOnce  sync.Once
	mongoURI  = "mongodb://admin:admin@mongo:27017/predictions_db?authSource=admin"
	database  = "predictions_db" // Database name
	clientErr error
)

// InitMongoDB initializes the MongoDB client (singleton pattern)
func InitMongoDB() {
	initOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		clientOptions := options.Client().ApplyURI(mongoURI)
		client, clientErr = mongo.Connect(ctx, clientOptions)
		if clientErr != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", clientErr)
		}

		// Check the connection
		if err := client.Ping(ctx, nil); err != nil {
			log.Fatalf("MongoDB connection error: %v", err)
		}

		log.Println("Connected to MongoDB successfully")
	})
}

// GetCollection returns a reference to a specific collection
func GetCollection(collectionName string) *mongo.Collection {
	if client == nil {
		InitMongoDB()
	}
	return client.Database(database).Collection(collectionName)
}
