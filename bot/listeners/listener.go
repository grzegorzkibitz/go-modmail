package listeners

import (
	"discord-bot-tickets/bot/context"
	"discord-bot-tickets/config"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

// RegisterListeners registers all listeners for the bot. i.e. messageCreate, messageDelete, etc.
func RegisterListeners(config *config.Config, state *state.State) {
	ctx := &context.Context{
		Config:  config,
		State:   state,
		Session: state.Session,
	}

	state.AddHandler(func(event *gateway.MessageCreateEvent) {
		HandleMessageCreate(ctx, event)
	})
}
