package listeners

import (
	"discord-bot-tickets/bot/commands/helpers/messages"
	"discord-bot-tickets/bot/services"
	"discord-bot-tickets/bot/tickets"
	logger "discord-bot-tickets/logging"

	"github.com/diamondburned/arikawa/v3/gateway"
)

func HandleMessageCreate(service *services.BotService, event *gateway.MessageCreateEvent) {
	if event.Author.Bot {
		return
	}

	// Check if the user has an active ticket
	ticket, err := tickets.GetActiveTicket(service.Config(), service.State(), &event.Author)
	if err != nil {
		return
	}

	if event.GuildID.IsValid() {
		return
	}

	// If the user has an active ticket, update it
	if ticket != nil {
		if err = tickets.UpdateTicket(service.Config(), service.State(), event.Author, tickets.RegularMessage{Message: event.Message}); err != nil {
			logger.Error(err.Error())
		}

		if err = messages.ReactToMessage(service.State(), event.Message, "✅"); err != nil {
			logger.Error(err.Error())
		}

		return
	}

	// If the user does not have an active ticket, create one
	if _, err = tickets.CreateTicket(service.Config(), service.State(), event.Author, event.Message); err != nil {
		logger.Error(err.Error())
	}

	if err = messages.ReactToMessage(service.State(), event.Message, "✅"); err != nil {
		logger.Error(err.Error())
	}
}

// HandleChannelDelete handles channel deletion events and cleans up the ticket cache
func HandleChannelDelete(service *services.BotService, event *gateway.ChannelDeleteEvent) {
	// Check if the deleted channel was a ticket
	isTicket, err := tickets.IsChannelTicket(*service.State(), &event.Channel)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if isTicket {
		// Get the author from the channel topic
		author, err := tickets.GetAuthorFromChannel(service.State(), &event.Channel)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		if author != nil {
			// Remove the ticket from the cache
			tickets.RemoveTicketFromCache(author.ID)
		}
	}
}
