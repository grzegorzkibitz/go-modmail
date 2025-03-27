package bot

import (
	"context"
	"discord-bot-tickets/bot/listeners"
	"discord-bot-tickets/config"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"log"
)

func InitializeBot(config *config.Config) {
	var intents = []gateway.Intents{
		gateway.IntentGuilds,
		gateway.IntentGuildMessages,
		gateway.IntentDirectMessages,
	}

	router := cmdroute.NewRouter()

	botState := state.New("Bot " + config.Discord.Token)
	botState.AddInteractionHandler(router)

	for _, intent := range intents {
		botState.AddIntents(intent)
	}

	RegisterCommands(router, botState)
	listeners.RegisterListeners(config, botState)

	if err := botState.Connect(context.TODO()); err != nil {
		log.Println("cannot connect:", err)
	}
}
