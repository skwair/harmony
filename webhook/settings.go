package webhook

import "github.com/skwair/discord/optional"

// Settings describes a webhook's settings.
type Settings struct {
	Name *optional.String `json:"name,omitempty"`
	// Avatar is a data URI scheme that support JPG, GIF, and PNG formats, see
	// https://discordapp.com/developers/docs/resources/user#avatar-data
	// for more information.
	Avatar    *optional.String `json:"avatar,omitempty"`
	ChannelID *optional.String `json:"channel_id,omitempty"`
}

// Setting is a function that configures a webhook.
type Setting func(*Settings)

// NewSettings returns new Settings to modify a webhook.
func NewSettings(opts ...Setting) *Settings {
	c := &Settings{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithName sets the name of a webhook.
func WithName(name string) Setting {
	return func(s *Settings) {
		s.Name = optional.NewString(name)
	}
}

// WithAvatar sets the avatar of a webhook.
func WithAvatar(uri string) Setting {
	return func(s *Settings) {
		s.Avatar = optional.NewString(uri)
	}
}

// WithChannelID sets the channel ID of a webhook.
func WithChannelID(id string) Setting {
	return func(s *Settings) {
		s.ChannelID = optional.NewString(id)
	}
}
