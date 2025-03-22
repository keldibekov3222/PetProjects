package repositories

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	_ "time"
)

type MongoDBRepository struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewMongoDBRepository(uri string, dbName string) (*MongoDBRepository, error) {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
		return nil, err
	}
	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	db := client.Database(dbName)
	log.Println("Successfully connected to MongoDB!")

	return &MongoDBRepository{Client: client, DB: db}, nil
}
