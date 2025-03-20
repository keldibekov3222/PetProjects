package database

import (
	"context"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

type MongoInstance struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var mg *MongoInstance
var once sync.Once

const dbName = "MongoDB"

var MongoURI string

func init() {
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	MongoURI = viper.GetString("MONGO_URI")
}

func connect() {
	clientOptions := options.Client().ApplyURI(MongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("Could not ping MongoDB: %v", err)
	}

	db := client.Database(dbName)
	mg = &MongoInstance{
		Client: client,
		DB:     db,
	}
}

func GetMongoInstance() *MongoInstance {
	once.Do(connect)
	return mg
}
