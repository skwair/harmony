package invite

import "github.com/skwair/harmony/optional"

// Settings describes how to modify a guild role. All fields are optional.
type Settings struct {
	MaxAge    *optional.Int  `json:"max_age,omitempty"`
	MaxUses   *optional.Int  `json:"max_uses,omitempty"`
	Temporary *optional.Bool `json:"temporary,omitempty"`
	Unique    *optional.Bool `json:"unique,omitempty"`
}

// Setting is a function that configures a guild role.
type Setting func(*Settings)

// NewSettings returns new Settings to modify a a guild role.
func NewSettings(opts ...Setting) *Settings {
	s := &Settings{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithMaxAge sets the maximum age (in seconds) of a channel invite.
func WithMaxAge(max int) Setting {
	return func(s *Settings) {
		s.MaxAge = optional.NewInt(max)
	}
}

// WithMaxUses sets maximum number of times this invite can be used.
func WithMaxUses(max int) Setting {
	return func(s *Settings) {
		s.MaxUses = optional.NewInt(max)
	}
}

// WithTemporary sets whether this invite only grants temporary membership.
func WithTemporary(yes bool) Setting {
	return func(s *Settings) {
		s.Temporary = optional.NewBool(yes)
	}
}

// WithUnique sets whether this invite is unique. If set to true, don't try to
// reuse a similar invite (useful for creating many unique one time use invites).
func WithUnique(yes bool) Setting {
	return func(s *Settings) {
		s.Unique = optional.NewBool(yes)
	}
}
