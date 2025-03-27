package services

import (
	"discord-bot-tickets/config"

	"github.com/diamondburned/arikawa/v3/session"
	"github.com/diamondburned/arikawa/v3/state"
)

// BotService encapsulates the bot's state and configuration
type BotService struct {
	config  *config.Config
	state   *state.State
	session *session.Session
}

// NewBotService creates a new BotService instance
func NewBotService(cfg *config.Config, st *state.State) *BotService {
	return &BotService{
		config:  cfg,
		state:   st,
		session: st.Session,
	}
}

// Config returns the bot's configuration
func (s *BotService) Config() *config.Config {
	return s.config
}

// State returns the bot's state
func (s *BotService) State() *state.State {
	return s.state
}

// Session returns the bot's session
func (s *BotService) Session() *session.Session {
	return s.session
}
