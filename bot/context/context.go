// context/context.go
package context

import (
	"discord-bot-tickets/config"
	"github.com/diamondburned/arikawa/v3/session"
	"github.com/diamondburned/arikawa/v3/state"
)

type Context struct {
	Config  *config.Config
	State   *state.State
	Session *session.Session
}
