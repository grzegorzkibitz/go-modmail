package main

import (
	"discord-bot-tickets/bot"
	"discord-bot-tickets/config"
	"discord-bot-tickets/database"
	"fmt"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Println("Config loaded successfully:", cfg)

	// Pass the configuration to the database connection function
	db, err := database.Connect(cfg)

	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
		return
	}

	_ = db

	bot.InitializeBot(cfg)
}
