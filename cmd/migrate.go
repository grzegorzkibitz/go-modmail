package main

import (
	"discord-bot-tickets/config"
	"discord-bot-tickets/database"
	"log"
	"os"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	var fresh = false

	// Check for the --fresh flag in command-line arguments
	for _, arg := range os.Args[1:] {
		if arg == "--fresh" {
			fresh = true
			break
		} else if arg == "--rollback" {
			database.RollbackMigration(cfg)
			return
		}
	}

	database.MigrateDatabase(cfg, fresh)

	log.Println("Migrations applied successfully!")
}
