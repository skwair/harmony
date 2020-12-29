package discord

import (
	"time"

	"github.com/skwair/harmony/optional"
)

// InviteSettings describes how to create a channel invite. All fields are optional.
type InviteSettings struct {
	MaxAge    *optional.Int  `json:"max_age,omitempty"`
	MaxUses   *optional.Int  `json:"max_uses,omitempty"`
	Temporary *optional.Bool `json:"temporary,omitempty"`
	Unique    *optional.Bool `json:"unique,omitempty"`
}

// InviteSetting is a function that configures an invite.
type InviteSetting func(*InviteSettings)

// NewInviteSettings returns new InviteSettings to create an invite.
func NewInviteSettings(opts ...InviteSetting) *InviteSettings {
	s := &InviteSettings{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithInviteMaxAge sets the delay before an invitation expires.
func WithInviteMaxAge(age time.Duration) InviteSetting {
	return func(s *InviteSettings) {
		s.MaxAge = optional.NewInt(int(age.Seconds()))
	}
}

// WithInviteMaxUses sets the maximum number of uses of an invitation.
func WithInviteMaxUses(uses int) InviteSetting {
	return func(s *InviteSettings) {
		s.MaxUses = optional.NewInt(uses)
	}
}

// WithInviteTemporary sets whether this invite only grants a temporary membership.
func WithInviteTemporary(yes bool) InviteSetting {
	return func(s *InviteSettings) {
		s.Temporary = optional.NewBool(yes)
	}
}

// WithInviteUnique determines if we should try to reuse a similar existing
// invite or not (enable this to create many single-use invites).
func WithInviteUnique(yes bool) InviteSetting {
	return func(s *InviteSettings) {
		s.Unique = optional.NewBool(yes)
	}
}
