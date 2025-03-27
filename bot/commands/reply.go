package commands

import (
	"context"
	"discord-bot-tickets/bot/services"
	"discord-bot-tickets/bot/tickets"
	logger "discord-bot-tickets/logging"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

func ReplyCommand(ctx context.Context, service *services.BotService, data cmdroute.CommandData) *api.InteractionResponseData {
	options := data.Options
	if len(options) == 0 {
		return &api.InteractionResponseData{
			Content: option.NewNullableString("Please provide a message to send."),
			Flags:   discord.EphemeralMessage,
		}
	}

	message := options[0].String()
	if message == "" {
		return &api.InteractionResponseData{
			Content: option.NewNullableString("Please provide a message to send."),
			Flags:   discord.EphemeralMessage,
		}
	}

	// get the author of the ticket from the channel
	author, err := tickets.GetAuthorFromChannel(service.State(), data.Event.Channel)
	if err != nil {
		return &api.InteractionResponseData{
			Content: option.NewNullableString("There was an error determining whether this channel is a ticket."),
			Flags:   discord.EphemeralMessage,
		}
	}

	if author == nil {
		return &api.InteractionResponseData{
			Content: option.NewNullableString("This channel is not a ticket"),
			Flags:   discord.EphemeralMessage,
		}
	}

	// Update the ticket with the reply
	if err = tickets.UpdateTicket(service.Config(), service.State(), *author, tickets.SlashCommandMessage{Message: message}); err != nil {
		logger.Error(err.Error())
		return &api.InteractionResponseData{
			Content: option.NewNullableString("Error updating ticket."),
			Flags:   discord.EphemeralMessage,
		}
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString("Message sent successfully."),
		Flags:   discord.EphemeralMessage,
	}
}

func GetReplyLocale() map[discord.Language]string {
	return map[discord.Language]string{}
}

func GetReplyDescription() string {
	return "Reply to a ModMail ticket!"
}

func GetReplyOptions() discord.CommandOptions {
	// return a command option with a string type
	return discord.CommandOptions{
		&discord.StringOption{
			OptionName:  "message",
			Description: "The message to reply with",
			Required:    true,
		},
	}
}
