package discord

import (
	"time"

	"github.com/skwair/harmony/optional"
)

// Invite represents a code that when used, adds a user to a guild or group DM channel.
type Invite struct {
	Code                     string  `json:"code"`
	Guild                    Guild   `json:"guild"` // Nil if this invite is for a group DM channel.
	Channel                  Channel `json:"channel"`
	ApproximatePresenceCount int     `json:"approximate_presence_count"`
	ApproximateMemberCount   int     `json:"approximate_member_count"`

	InviteMetadata
}

// InviteMetadata contains additional information about an Invite.
type InviteMetadata struct {
	Inviter   User      `json:"inviter"`
	Uses      int       `json:"uses"`
	MaxUses   int       `json:"max_uses"`
	MaxAge    int       `json:"max_age"`
	Temporary bool      `json:"temporary"`
	CreatedAt time.Time `json:"created_at"`
	Revoked   bool      `json:"revoked"`
}

// InviteSettings describes how to create a channel invite. All fields are optional.
type InviteSettings struct {
	MaxAge    *optional.Int  `json:"max_age,omitempty"`
	MaxUses   *optional.Int  `json:"max_uses,omitempty"`
	Temporary *optional.Bool `json:"temporary,omitempty"`
	Unique    *optional.Bool `json:"unique,omitempty"`
}

// InviteSetting is a function that configures a channel.
type InviteSetting func(*InviteSettings)

// NewInviteSettings returns new InviteSettings to modify a a channel.
func NewInviteSettings(opts ...InviteSetting) *InviteSettings {
	s := &InviteSettings{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithMaxAge sets the delay before an invitation expires.
func WithMaxAge(age time.Duration) InviteSetting {
	return func(s *InviteSettings) {
		s.MaxAge = optional.NewInt(int(age.Seconds()))
	}
}

// WithMaxUses sets the maximum number of uses of an invitation.
func WithMaxUses(uses int) InviteSetting {
	return func(s *InviteSettings) {
		s.MaxUses = optional.NewInt(uses)
	}
}

// WithTemporary sets the maximum number of uses of an invitation.
func WithTemporary(yes bool) InviteSetting {
	return func(s *InviteSettings) {
		s.Temporary = optional.NewBool(yes)
	}
}

// WithUnique determines if we should try to reuse a similar existing invite or
// not (useful for creating many unique one time use invites).
func WithUnique(yes bool) InviteSetting {
	return func(s *InviteSettings) {
		s.Unique = optional.NewBool(yes)
	}
}
