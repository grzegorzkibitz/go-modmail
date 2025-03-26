package listeners

import (
	"discord-bot-tickets/config"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

type BotContext struct {
	Config *config.Config
	State  *state.State
}

// RegisterListeners registers all listeners for the bot. i.e. messageCreate, messageDelete, etc.
func RegisterListeners(config *config.Config, state *state.State) {
	ctx := &BotContext{
		Config: config,
		State:  state,
	}

	state.AddHandler(func(event *gateway.MessageCreateEvent) {
		HandleMessageCreate(ctx, event)
	})
}
