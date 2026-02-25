package main

import (
	"log"

	"github.com/banan213/photobot/internal/bot"
	"github.com/banan213/photobot/internal/config"
	"github.com/banan213/photobot/internal/images"
	"github.com/banan213/photobot/internal/telegram"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	tg, err := telegram.New(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("telegram init: %v", err)
	}

	img := images.New(cfg.UnsplashKey)

	bl := bot.New(tg, img)

	if err := bl.Run(); err != nil {
		log.Fatalf("bot run: %v", err)
	}
}
