package messages

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
)

func ReactToMessage(state *state.State, message discord.Message, emoji discord.APIEmoji) error {
	err := state.React(message.ChannelID, message.ID, emoji)

	return err
}
