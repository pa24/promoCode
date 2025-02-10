package config

import (
	"github.com/pkg/errors"
	"os"
)

type Config struct {
	DatabaseURL string
}

func LoadConfig() (*Config, error) {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}
	return &Config{DatabaseURL: databaseUrl}, nil
}
