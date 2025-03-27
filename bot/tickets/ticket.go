package tickets

import (
	"discord-bot-tickets/bot/commands/helpers/colors"
	"discord-bot-tickets/config"
	"strings"
	"sync"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
)

// Ticket represents a ModMail ticket with its channel and owner
type Ticket struct {
	Channel *discord.Channel
	Author  *discord.User
}

// TicketCache stores active tickets in memory
type TicketCache struct {
	tickets map[discord.UserID]*Ticket
	mu      sync.RWMutex
}

var ticketCache = &TicketCache{
	tickets: make(map[discord.UserID]*Ticket),
}

// AddTicket adds a ticket to the cache
func (c *TicketCache) AddTicket(ticket *Ticket) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tickets[ticket.Author.ID] = ticket
}

// GetTicket retrieves a ticket from the cache
func (c *TicketCache) GetTicket(userID discord.UserID) *Ticket {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.tickets[userID]
}

// RemoveTicket removes a ticket from the cache
func (c *TicketCache) RemoveTicket(userID discord.UserID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.tickets, userID)
}

// RemoveTicketFromCache removes a ticket from the global cache
func RemoveTicketFromCache(userID discord.UserID) {
	ticketCache.RemoveTicket(userID)
}

// MessageContent represents a message that can be either a regular message or a slash command message
type MessageContent interface {
	GetContent() string
	IsPrivateChat() bool
	GetAuthor() discord.User
}

// RegularMessage implements MessageContent for regular Discord messages
type RegularMessage struct {
	Message discord.Message
}

func (m RegularMessage) GetContent() string {
	return m.Message.Content
}

func (m RegularMessage) IsPrivateChat() bool {
	return !m.Message.GuildID.IsValid()
}

func (m RegularMessage) GetAuthor() discord.User {
	return m.Message.Author
}

// SlashCommandMessage implements MessageContent for slash command messages
type SlashCommandMessage struct {
	Message string
	Author  discord.User
}

func (m SlashCommandMessage) GetContent() string {
	return m.Message
}

func (m SlashCommandMessage) IsPrivateChat() bool {
	return false // Slash commands are always from guild channels
}

func (m SlashCommandMessage) GetAuthor() discord.User {
	return m.Author
}

// CreateTicket creates a ticket for a user
//
// Returns: a pointer to a Ticket and an error if any
func CreateTicket(config *config.Config, state *state.State, author discord.User, message discord.Message) (*Ticket, error) {
	data := api.CreateChannelData{
		Name:       author.Username,
		Type:       discord.GuildText,
		Topic:      "User: " + author.ID.String(),
		CategoryID: config.Discord.CategoryID,
	}

	channel, err := state.CreateChannel(config.Discord.GuildID, data)
	if err != nil {
		return nil, err
	}

	embed := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: author.Username,
			Icon: author.AvatarURL(),
		},
		Fields: []discord.EmbedField{
			{
				Value: message.Content,
			},
		},
		Footer: &discord.EmbedFooter{
			Text: "ModMail",
		},
	}

	_, err = state.SendEmbeds(channel.ID, embed)
	if err != nil {
		return nil, err
	}

	ticket := &Ticket{
		Channel: channel,
		Author:  &author,
	}
	// Add to cache
	ticketCache.AddTicket(ticket)
	return ticket, nil
}

// UpdateTicket updates the ticket with the latest message
//
// Returns: an error if any
func UpdateTicket(config *config.Config, state *state.State, user discord.User, message MessageContent) error {
	var embedColor discord.Color

	ticket, err := GetActiveTicket(config, state, &user)
	if err != nil {
		return err
	}

	// Determine which user to use for the author field
	messageAuthor := message.GetAuthor()
	if message.IsPrivateChat() {
		messageAuthor = *ticket.Author
	}

	if message.IsPrivateChat() {
		embedColor = colors.GetColor(colors.Yellow)
	} else {
		embedColor = colors.GetColor(colors.Green)
	}

	ticketChannelEmbed := discord.Embed{
		Color: embedColor,
		Author: &discord.EmbedAuthor{
			Name: messageAuthor.Username,
			Icon: messageAuthor.AvatarURL(),
		},
		Fields: []discord.EmbedField{
			{
				Value: message.GetContent(),
			},
		},
		Timestamp: discord.NowTimestamp(),
		Footer: &discord.EmbedFooter{
			Text: "ModMail",
		},
	}

	_, err = state.SendEmbeds(ticket.Channel.ID, ticketChannelEmbed)
	if err != nil {
		return err
	}

	// Only send DM if it's not a private chat reply
	if !message.IsPrivateChat() {
		privateChannel, err := state.CreatePrivateChannel(ticket.Author.ID)
		if err != nil {
			return err
		}

		privateChannelEmbed := discord.Embed{
			Color: embedColor,
			Author: &discord.EmbedAuthor{
				Name: messageAuthor.Username,
				Icon: messageAuthor.AvatarURL(),
			},
			Fields: []discord.EmbedField{
				{
					Value: message.GetContent(),
				},
			},
			Timestamp: discord.NowTimestamp(),
			Footer: &discord.EmbedFooter{
				Text: "ModMail",
			},
		}

		_, err = state.SendEmbeds(privateChannel.ID, privateChannelEmbed)
		if err != nil {
			return err
		}
	}

	return nil
}

// IsChannelTicket checks if a specific channel is a ticket
//
// Returns: a boolean and an error if any
func IsChannelTicket(state state.State, channel *discord.Channel) (isTicket bool, err error) {
	if !strings.Contains(channel.Topic, "User: ") {
		return false, nil
	}

	parts := strings.Split(channel.Topic, "User: ")
	if len(parts) < 2 {
		return false, nil
	}

	stripped := parts[1]

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
// Returns: a pointer to a Ticket and an error if any
func GetActiveTicket(config *config.Config, state *state.State, Author *discord.User) (*Ticket, error) {
	// First check the cache
	if ticket := ticketCache.GetTicket(Author.ID); ticket != nil {
		return ticket, nil
	}

	// If not in cache, search through Discord channels
	channels, err := state.Channels(config.Discord.GuildID)
	if err != nil {
		return nil, err
	}

	// Check if the user has an active ticket using the topic
	for _, channel := range channels {
		if channel.Topic == "User: "+Author.ID.String() {
			channelCopy := channel
			ticket := &Ticket{
				Channel: &channelCopy,
				Author:  Author,
			}
			// Add to cache
			ticketCache.AddTicket(ticket)
			return ticket, nil
		}
	}

	return nil, nil
}

// GetAuthorFromChannel gets the author of a ticket if the channel is a ticket
func GetAuthorFromChannel(state *state.State, channel *discord.Channel) (*discord.User, error) {
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
