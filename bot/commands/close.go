package commands

import (
	"context"
	"discord-bot-tickets/bot/commands/helpers/language"
	"discord-bot-tickets/bot/services"
	"discord-bot-tickets/bot/tickets"
	logger "discord-bot-tickets/logging"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

func CloseCommand(ctx context.Context, service *services.BotService, data cmdroute.CommandData) *api.InteractionResponseData {
	// Get the channel where the command was used
	channel, err := service.State().Channel(data.Event.ChannelID)
	if err != nil {
		return &api.InteractionResponseData{
			Content: option.NewNullableString(language.GetTranslation("general.errors.channel")),
			Flags:   discord.EphemeralMessage,
		}
	}

	// Check if this is a ticket channel
	isTicket, err := tickets.IsChannelTicket(*service.State(), channel)
	if err != nil {
		return &api.InteractionResponseData{
			Content: option.NewNullableString(language.GetTranslation("general.errors.generic")),
			Flags:   discord.EphemeralMessage,
		}
	}

	if !isTicket {
		return &api.InteractionResponseData{
			Content: option.NewNullableString(language.GetTranslation("general.errors.not_a_ticket")),
			Flags:   discord.EphemeralMessage,
		}
	}

	// Get the ticket owner from the channel topic
	ticketOwner, err := tickets.GetAuthorFromChannel(service.State(), channel)
	if err != nil {
		return &api.InteractionResponseData{
			Content: option.NewNullableString(language.GetTranslation("general.errors.owner")),
			Flags:   discord.EphemeralMessage,
		}
	}

	if ticketOwner == nil {
		return &api.InteractionResponseData{
			Content: option.NewNullableString(language.GetTranslation("general.errors.not_a_ticket")),
			Flags:   discord.EphemeralMessage,
		}
	}

	// Create an embed to notify the user
	embed := discord.Embed{
		Title:       language.GetTranslation("embeds.ticket_closed.title"),
		Description: language.GetTranslation("embeds.ticket_closed.description"),
		Color:       0xFF0000, // Red color
		Footer: &discord.EmbedFooter{
			Text: language.GetTranslation("embeds.ticket_closed.footer"),
		},
	}

	// Create DM channel with the user
	dmChannel, err := service.State().CreatePrivateChannel(ticketOwner.ID)
	if err != nil {
		logger.Error("Failed to create DM channel with user: " + err.Error())
	} else {
		// Send the embed to the user
		_, err = service.State().SendMessage(dmChannel.ID, "", embed)
		if err != nil {
			logger.Error("Failed to send close notification to user: " + err.Error())
		}
	}

	// Log the ticket closure
	logger.Info("Ticket closed by staff member " + data.Event.Member.User.ID.String())

	// Delete the channel with audit log reason
	err = service.State().DeleteChannel(channel.ID, api.AuditLogReason("Ticket closed by "+data.Event.Member.User.Tag()))
	if err != nil {
		return &api.InteractionResponseData{
			Content: option.NewNullableString(language.GetTranslation("commands.close.error")),
			Flags:   discord.EphemeralMessage,
		}
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString(language.GetTranslation("commands.close.success")),
		Flags:   discord.EphemeralMessage,
	}
}

func GetCloseLocale() map[discord.Language]string {
	return map[discord.Language]string{}
}

func GetCloseDescription() string {
	return "Close a ModMail ticket"
}

func GetCloseOptions() discord.CommandOptions {
	return discord.CommandOptions{}
}
