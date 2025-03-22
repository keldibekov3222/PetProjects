package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	MongoURI   string
	MongoDB    string
}

func LoadConfig() *Config {
	err := godotenv.Load("config/config.env") // или абсолютный путь
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	return &Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		MongoURI:   os.Getenv("MONGO_URI"),
		MongoDB:    os.Getenv("MONGO_DB"),
	}
}
