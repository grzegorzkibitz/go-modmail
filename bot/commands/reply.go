package commands

import (
	"context"
	"discord-bot-tickets/bot/commands/helpers"
	"discord-bot-tickets/bot/listeners"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
)

func ReplyCommand(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	//isTicket, err := listeners.IsChannelTicket()
}

func GetReplyLocale() map[discord.Language]string {
	return map[discord.Language]string{}
}

func GetReplyDescription() string {
	return "Reply to a ModMail ticket!"
}
