package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabasePath string
	SecretKey    string
	Port         string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		DatabasePath: os.Getenv("DATABASE_PATH"),
		SecretKey:    os.Getenv("SECRET_KEY"),
		Port:         os.Getenv("PORT"),
	}, nil
}
