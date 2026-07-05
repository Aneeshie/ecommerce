package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	ServerAddress string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")

	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = ":8080"
	}

	if dbURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	return &Config{
		DatabaseURL:   dbURL,
		ServerAddress: serverAddress,
	}, nil
}
