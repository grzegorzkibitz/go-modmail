package commands

import (
	"context"
	"discord-bot-tickets/bot/commands/helpers"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
)

func PingCommand(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	embed := discord.Embed{
		Color: helpers.Color("ff3333"),
		Fields: []discord.EmbedField{
			{
				Value: "Pong! 🏓",
			},
		},
	}
	return &api.InteractionResponseData{
		Embeds: &[]discord.Embed{embed},
	}
}

func GetPingLocale() map[discord.Language]string {
	return map[discord.Language]string{
		discord.EnglishUK: "Ping!",
	}
}

func GetPingDescription() string {
	return "Ping!"
}
