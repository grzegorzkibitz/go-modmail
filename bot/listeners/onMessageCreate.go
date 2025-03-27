package listeners

import (
	"discord-bot-tickets/bot/context"
	"discord-bot-tickets/bot/tickets"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func HandleMessageCreate(context *context.Context, event *gateway.MessageCreateEvent) {
	if event.Author.Bot {
		return
	}

	// Check if the user has an active ticket
	channel, err := tickets.GetActiveTicket(context.Config, context.State, event.Author)
	if err != nil {
		return
	}

	// If it is a guild message, check if it is a ticket
	if event.GuildID.IsValid() {
		return
	}

	// If the user has an active ticket, update it
	if channel != nil {
		err = tickets.UpdateTicket(context.State, *channel, event.Author, event.Message)
		if err != nil {
			return
		}
		return
	}

	// If the user does not have an active ticket, create one
	_, err = tickets.CreateTicket(context.Config, context.State, event.Author, event.Message)
	if err != nil {
		return
	}
}
