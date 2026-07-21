package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	ServerAddress string
	JwtSecret     string
	Env           string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")

	JwtSecret := os.Getenv("JWT_SECRET")

	serverAddress := os.Getenv("SERVER_ADDRESS")

	env := os.Getenv("ENV")

	if serverAddress == "" {
		serverAddress = ":8080"
	}

	if dbURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	if JwtSecret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}

	if env == "" {
		env = "dev"
	}

	return &Config{
		DatabaseURL:   dbURL,
		ServerAddress: serverAddress,
		JwtSecret:     JwtSecret,
		Env:           env,
	}, nil
}
