package listeners

import (
	"discord-bot-tickets/bot/services"
	"discord-bot-tickets/bot/tickets"
	logger "discord-bot-tickets/logging"

	"github.com/diamondburned/arikawa/v3/gateway"
)

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

		logger.Info("Removed ticket ", event.Channel.ID.String(), " from cache as it was deleted")
	}
}
