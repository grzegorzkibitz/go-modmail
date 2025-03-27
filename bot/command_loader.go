package bot

import (
	"context"
	"discord-bot-tickets/bot/commands"
	"discord-bot-tickets/bot/services"
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
)

// CommandHandler represents a Discord command handler
type CommandHandler func(ctx context.Context, service *services.BotService, data cmdroute.CommandData) *api.InteractionResponseData

// CommandRegistry holds all registered commands
var CommandRegistry = map[string]CommandHandler{
	"reply": commands.ReplyCommand,
	"close": commands.CloseCommand,
}

// CommandData holds all command data
var commandData = []api.CreateCommandData{
	{Name: "reply", Description: commands.GetReplyDescription(), DescriptionLocalizations: commands.GetReplyLocale(), Options: commands.GetReplyOptions()},
	{Name: "close", Description: commands.GetCloseDescription(), DescriptionLocalizations: commands.GetCloseLocale()},
}

// RegisterCommands loads and registers all commands
func RegisterCommands(router *cmdroute.Router, service *services.BotService) {
	for name, handler := range CommandRegistry {
		router.AddFunc(name, func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
			return handler(ctx, service, data)
		})
	}

	if err := cmdroute.OverwriteCommands(service.State(), commandData); err != nil {
		log.Fatalln("cannot update commands:", err)
	}

	log.Printf("Registered %d commands", len(CommandRegistry))
}
