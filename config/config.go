package config

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strconv"
)

type Config struct {
	Discord DiscordConfig
	DB      MySqlConfig
	Port    string
}

type MySqlConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Table    string
}

type DiscordConfig struct {
	Token      string
	GuildID    discord.GuildID
	CategoryID string
}

type ErrMissingEnvVar string

func (e ErrMissingEnvVar) Error() string {
	return fmt.Sprintf("missing environment variable: %s", string(e))
}

// LoadConfig loads the configuration from environment variables.
//
// Returns: a pointer to the Config struct and an error if any
func LoadConfig() (*Config, error) {

	if os.Getenv("DISCORD_GUILD_ID") == "" {
		return nil, ErrMissingEnvVar("DISCORD_GUILD_ID")
	}

	guildID, err := strconv.ParseUint(os.Getenv("DISCORD_GUILD_ID"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid DISCORD_GUILD_ID: %v", err)
	}

	cfg := &Config{
		Discord: DiscordConfig{
			Token:      os.Getenv("DISCORD_TOKEN"),
			GuildID:    discord.GuildID(guildID),
			CategoryID: os.Getenv("DISCORD_CATEGORY_ID"),
		},
		DB: MySqlConfig{
			Username: os.Getenv("MYSQL_USER"),
			Password: os.Getenv("MYSQL_PASSWORD"),
			Host:     os.Getenv("MYSQL_HOST"),
			Port:     os.Getenv("MYSQL_PORT"),
			Table:    os.Getenv("MYSQL_TABLE"),
		},
	}

	//make sure all values are set
	if cfg.DB.Username == "" {
		return nil, ErrMissingEnvVar("MYSQL_USER")
	}

	if cfg.DB.Host == "" {
		return nil, ErrMissingEnvVar("MYSQL_HOST")
	}

	if cfg.DB.Port == "" {
		return nil, ErrMissingEnvVar("MYSQL_PORT")
	}

	if cfg.DB.Table == "" {
		return nil, ErrMissingEnvVar("MYSQL_TABLE")
	}

	if cfg.Discord.Token == "" {
		return nil, ErrMissingEnvVar("DISCORD_TOKEN")
	}

	return cfg, nil
}
