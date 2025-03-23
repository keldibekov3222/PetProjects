package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	DBUser        string
	DBPassword    string
	DBHost        string
	DBPort        string
	DBName        string
	MongoURI      string
	MongoDB       string
	RedisAddr     string // Адрес Redis
	RedisPassword string // Пароль Redis
	RedisDB       int    // Номер базы данных Redis
}

func LoadConfig() *Config {
	err := godotenv.Load("config/config.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	return &Config{
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBName:        os.Getenv("DB_NAME"),
		MongoURI:      os.Getenv("MONGO_URI"),
		MongoDB:       os.Getenv("MONGO_DB"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       atoi(os.Getenv("REDIS_DB")),
	}
}

func atoi(str string) int {
	if str == "" {
		return 0
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		log.Fatalf("Error converting string to int: %v", err)
	}
	return i
}
