package main

import (
	"discord-bot-tickets/bot"
	"discord-bot-tickets/config"
	"discord-bot-tickets/database"
	logger "discord-bot-tickets/logging"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if err := logger.Init(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Pass the configuration to the database connection function
	db, err := database.Connect(cfg)

	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
		return
	}

	_ = db

	bot.InitializeBot(cfg)
}
