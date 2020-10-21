package discord

import (
	"encoding/json"

	"github.com/skwair/harmony/optional"
)

// WebhookParameters are the parameters available when executing a
// webhook with ExecWebhook.
type WebhookParameters struct {
	Content   *optional.String `json:"content,omitempty"`
	Username  *optional.String `json:"username,omitempty"`
	AvatarURL *optional.String `json:"avatar_url,omitempty"`
	TTS       *optional.Bool   `json:"tts,omitempty"`
	Embeds    []MessageEmbed   `json:"embeds,omitempty"`
	Files     []File           `json:"-"`
}

// Bytes implements the rest.MultipartPayload interface so WebhookParameters can be used as
// a payload with the rest.MultipartFromFiles function.
func (p *WebhookParameters) Bytes() ([]byte, error) {
	return json.Marshal(p)
}

// WebhookParameter is a function that sets webhook parameters.
type WebhookParameter func(*WebhookParameters)

// NewWebhookParameters returns new WebhookParameters to modify a a webhook.
func NewWebhookParameters(opts ...WebhookParameter) *WebhookParameters {
	s := &WebhookParameters{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithWebhookContent sets the content of a webhook.
func WithWebhookContent(content string) WebhookParameter {
	return func(s *WebhookParameters) {
		s.Content = optional.NewString(content)
	}
}

// WithWebhookUsername sets the content of a webhook.
func WithWebhookUsername(content string) WebhookParameter {
	return func(s *WebhookParameters) {
		s.Username = optional.NewString(content)
	}
}

// WithWebhookAvatarURL sets the content of a webhook.
func WithWebhookAvatarURL(content string) WebhookParameter {
	return func(s *WebhookParameters) {
		s.AvatarURL = optional.NewString(content)
	}
}

// WithWebhookTTS sets the content of a webhook.
func WithWebhookTTS(yes bool) WebhookParameter {
	return func(s *WebhookParameters) {
		s.TTS = optional.NewBool(yes)
	}
}

// WithWebhookEmbeds sets the content of a webhook.
func WithWebhookEmbeds(embeds []MessageEmbed) WebhookParameter {
	return func(s *WebhookParameters) {
		s.Embeds = embeds
	}
}

// WithWebhookFiles sets the content of a webhook.
func WithWebhookFiles(files []File) WebhookParameter {
	return func(s *WebhookParameters) {
		s.Files = files
	}
}
