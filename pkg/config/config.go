package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration values.
type Config struct {
	OpenAIKey  string
	DBHost     string
	DBName     string
	DBUser     string
	DBPassword string
}

// LoadConfig loads the application configuration from .env file.
func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		OpenAIKey:  os.Getenv("OPENAI_KEY"),
		DBHost:     os.Getenv("DB_HOST"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
	}, nil
}
