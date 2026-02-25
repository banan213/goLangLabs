package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	UnsplashKey   string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	t := os.Getenv("TELEGRAM_APITOKEN")
	u := os.Getenv("AccessKey")

	if t == "" {
		return nil, fmt.Errorf("TELEGRAM_APITOKEN not set in environment")
	}
	if u == "" {
		return nil, fmt.Errorf("AccessKey (Unsplash) not set in environment")
	}

	return &Config{
		TelegramToken: t,
		UnsplashKey:   u,
	}, nil
}
