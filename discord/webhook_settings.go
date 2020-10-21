package discord

import "github.com/skwair/harmony/optional"

// WebhookSettings describes a webhook's settings.
type WebhookSettings struct {
	Name *optional.String `json:"name,omitempty"`
	// Avatar is a data URI scheme that support JPG, GIF, and PNG formats, see
	// https://discord.com/developers/docs/resources/user#avatar-data
	// for more information.
	Avatar    *optional.String `json:"avatar,omitempty"`
	ChannelID *optional.String `json:"channel_id,omitempty"`
}

// WebhookSetting is a function that configures a webhook.
type WebhookSetting func(*WebhookSettings)

// NewWebhookSettings returns new WebhookSettings to modify a a webhook.
func NewWebhookSettings(opts ...WebhookSetting) *WebhookSettings {
	s := &WebhookSettings{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithWebhookName sets the name of a webhook.
func WithWebhookName(name string) WebhookSetting {
	return func(s *WebhookSettings) {
		s.Name = optional.NewString(name)
	}
}

// WithWebhookAvatar sets the avatar of a webhook.
func WithWebhookAvatar(name string) WebhookSetting {
	return func(s *WebhookSettings) {
		s.Avatar = optional.NewString(name)
	}
}

// WithWebhookChannel sets the channel ID of a webhook.
func WithWebhookChannel(id string) WebhookSetting {
	return func(s *WebhookSettings) {
		s.ChannelID = optional.NewString(id)
	}
}
