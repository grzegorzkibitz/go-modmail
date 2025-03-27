package listeners

import (
	"discord-bot-tickets/bot/services"

	"github.com/diamondburned/arikawa/v3/gateway"
)

// RegisterListeners registers all listeners for the bot. i.e. messageCreate, messageDelete, etc.
func RegisterListeners(service *services.BotService) {
	service.State().AddHandler(func(event *gateway.MessageCreateEvent) {
		HandleMessageCreate(service, event)
	})

	service.State().AddHandler(func(event *gateway.ChannelDeleteEvent) {
		HandleChannelDelete(service, event)
	})
}
