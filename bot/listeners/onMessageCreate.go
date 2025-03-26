package listeners

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

func HandleMessageCreate(ctx *BotContext, event *gateway.MessageCreateEvent) {
	if event.Author.Bot || event.GuildID.IsValid() {
		return
	}

	channel, err := hasActiveTicket(ctx, event.Author)
	if err != nil {
		return
	}

	if channel != nil {
		err = updateTicket(ctx.State, *channel, event.Author, event.Message)
		if err != nil {
			return
		}
		return
	}

	err = createTicket(ctx, event.Author, event.Message)
	if err != nil {
		return
	}
}

func hasActiveTicket(ctx *BotContext, Author discord.User) (*discord.Channel, error) {
	channels, err := ctx.State.Channels(ctx.Config.Discord.GuildID)
	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		if channel.Topic == "User: "+Author.ID.String() {
			return &channel, nil
		}
	}

	return nil, nil
}

func createTicket(ctx *BotContext, Author discord.User, message discord.Message) error {
	data := api.CreateChannelData{
		Name:  "ticket-" + Author.Username,
		Type:  discord.GuildText,
		Topic: "User: " + Author.ID.String(),
	}
	_, err := ctx.State.CreateChannel(ctx.Config.Discord.GuildID, data)
	return err
}

func updateTicket(state *state.State, channel discord.Channel, user discord.User, message discord.Message) error {
	embed := discord.Embed{
		Title: "Ticket",
		Fields: []discord.EmbedField{
			{
				Name:  "User",
				Value: user.Mention(),
			},
			{
				Name:  "Message",
				Value: message.Content,
			},
		},
		Timestamp: discord.NowTimestamp(),
	}

	_, err := state.SendEmbeds(channel.ID, embed)
	return err
}
