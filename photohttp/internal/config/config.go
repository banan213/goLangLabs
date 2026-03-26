package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	UnsplashKey string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	u := os.Getenv("AccessKey")

	if u == "" {
		return nil, fmt.Errorf("AccessKey (Unsplash) not set in environment")
	}

	return &Config{
		UnsplashKey: u,
	}, nil
}
