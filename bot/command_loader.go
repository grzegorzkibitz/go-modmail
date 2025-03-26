package bot

import (
	"context"
	"discord-bot-tickets/bot/commands"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/state"
	"log"
)

// CommandHandler represents a Discord command handler
type CommandHandler func(context.Context, cmdroute.CommandData) *api.InteractionResponseData

// CommandRegistry holds all registered commands
var CommandRegistry = map[string]CommandHandler{}

// CommandData holds all command data
var commandData = []api.CreateCommandData{
	{Name: "ping", Description: commands.GetPingDescription(), DescriptionLocalizations: commands.GetPingLocale()},
}

// RegisterCommands loads and registers all commands
func RegisterCommands(router *cmdroute.Router, state *state.State) {
	for _, cmd := range commandData {
		switch cmd.Name {
		case "ping":
			CommandRegistry[cmd.Name] = commands.PingCommand
			router.AddFunc("ping", func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
				return commands.PingCommand(ctx, data)
			})
		}
	}

	if err := cmdroute.OverwriteCommands(state, commandData); err != nil {
		log.Fatalln("cannot update commands:", err)
	}

	log.Printf("Registered %d commands", len(CommandRegistry))
}
