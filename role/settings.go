package role

import "github.com/skwair/discord/optional"

// Settings describes how to modify a guild role. All fields are optional.
type Settings struct {
	Name        *optional.String `json:"name,omitempty"`
	Permissions *optional.Int    `json:"permissions,omitempty"`
	Color       *optional.Int    `json:"color,omitempty"`
	Hoist       *optional.Bool   `json:"hoist,omitempty"`
	Mentionable *optional.Bool   `json:"mentionable,omitempty"`
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

// WithName sets the name of guild a role.
func WithName(name string) Setting {
	return func(s *Settings) {
		s.Name = optional.NewString(name)
	}
}

// WithPermissions sets the permissions of guild a role.
func WithPermissions(perm int) Setting {
	return func(s *Settings) {
		s.Permissions = optional.NewInt(int(perm))
	}
}

// WithColor sets the color of guild a role. It accepts hexadecimal value.
func WithColor(hexCode int) Setting {
	return func(s *Settings) {
		s.Color = optional.NewInt(int(hexCode))
	}
}

// WithHoist sets whether this guild role is hoisted.
func WithHoist(yes bool) Setting {
	return func(s *Settings) {
		s.Hoist = optional.NewBool(yes)
	}
}

// WithMentionable sets whether this guild role is mentionable by others.
func WithMentionable(yes bool) Setting {
	return func(s *Settings) {
		s.Mentionable = optional.NewBool(yes)
	}
}
