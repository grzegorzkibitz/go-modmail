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

	if ticket != nil {
		if err = tickets.UpdateTicket(service.Config(), service.State(), event.Author, tickets.RegularMessage{Message: event.Message}); err != nil {
			logger.Error(err.Error())
		}
	} else {
		if _, err = tickets.CreateTicket(service.Config(), service.State(), event.Author, event.Message); err != nil {
			logger.Error(err.Error())
		}
	}

	if err = messages.ReactToMessage(service.State(), event.Message, "✅"); err != nil {
		logger.Error(err.Error())
	}
}
