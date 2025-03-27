package tickets

import (
	"discord-bot-tickets/config"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"strings"
)

// CreateTicket creates a ticket for a user
//
// Returns: a pointer to a discord.Channel and an error if any
func CreateTicket(config *config.Config, state *state.State, Author discord.User, message discord.Message) (*discord.Channel, error) {
	data := api.CreateChannelData{
		Name:  "ticket-" + Author.Username,
		Type:  discord.GuildText,
		Topic: "User: " + Author.ID.String(),
	}

	channel, err := state.CreateChannel(config.Discord.GuildID, data)
	if err != nil {
		return nil, err
	}

	embed := discord.Embed{
		Title: "Ticket",
		Fields: []discord.EmbedField{
			{
				Name:  "User",
				Value: Author.Mention(),
			},
			{
				Name:  "Message",
				Value: message.Content,
			},
		},
	}

	_, err = state.SendEmbeds(channel.ID, embed)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

// UpdateTicket updates the ticket with the latest message
//
// Returns: an error if any
func UpdateTicket(state *state.State, channel discord.Channel, user discord.User, message discord.Message) error {
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

// IsChannelTicket checks if a specific channel is a ticket
//
// Returns: a boolean and an error if any
func IsChannelTicket(state state.State, channel *discord.Channel) (isTicket bool, err error) {
	stripped := strings.Split(channel.Topic, "User: ")[1]

	// Make sure that the given string is a valid snowflake/user ID
	user, err := discord.ParseSnowflake(stripped)

	if err != nil {
		return false, err
	}

	author, err := state.User(discord.UserID(user))

	if err != nil {
		return false, err
	}

	if author.Bot {
		return false, nil
	}

	return true, nil
}

// GetActiveTicket gets the active ticket of a user
//
// Returns: a pointer to a discord.Channel and an error if any
func GetActiveTicket(config *config.Config, state *state.State, Author discord.User) (*discord.Channel, error) {
	channels, err := state.Channels(config.Discord.GuildID)
	if err != nil {
		return nil, err
	}

	// Check if the user has an active ticket using the topic
	for _, channel := range channels {
		if channel.Topic == "User: "+Author.ID.String() {
			return &channel, nil
		}
	}

	return nil, nil
}

// GetTicketAuthor gets the author of a ticket
func GetTicketAuthor(state *state.State, channel *discord.Channel) (*discord.User, error) {
	stripped := strings.Split(channel.Topic, "User: ")[1]

	// Make sure that the given string is a valid snowflake/user ID
	user, err := discord.ParseSnowflake(stripped)

	if err != nil {
		return &discord.User{}, err
	}

	author, err := state.User(discord.UserID(user))

	if err != nil {
		return &discord.User{}, err
	}

	return author, nil
}

// SendDirectMessage sendReply sends a reply to the user
func SendDirectMessage(state *state.State, user *discord.User, message discord.Message) {
	embed := discord.Embed{
		Title: "Reply",
		Fields: []discord.EmbedField{
			{
				Name:  "Message",
				Value: message.Content,
			},
		},
		Timestamp: discord.NowTimestamp(),
	}

	// get the user's DM channel
	channel, err := state.CreatePrivateChannel(user.ID)

	if err != nil {
		return
	}

	_, err = state.SendEmbeds(channel.ID, embed)
}
