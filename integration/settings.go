package integration

import "github.com/skwair/harmony/optional"

// Settings describes a guild integration's settings.
type Settings struct {
	// The behavior when an integration subscription lapses.
	ExpireBehavior *optional.Int
	// Period (in seconds) where the integration will ignore lapsed subscriptions.
	ExpireGracePeriod *optional.Int
	// Whether emoticons should be synced for this integration (twitch only currently).
	EnableEmoticons *optional.Bool
}

// Setting is a function that configures an integration.
type Setting func(*Settings)

// NewSettings returns new Settings to modify an integration.
func NewSettings(opts ...Setting) *Settings {
	c := &Settings{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithExpireBehavior sets the expiration behaviour of an integration.
func WithExpireBehavior(t int) Setting {
	return func(s *Settings) {
		s.ExpireBehavior = optional.NewInt(t)
	}
}

// WithExpireGracePeriod sets the expire grace period of an integration.
func WithExpireGracePeriod(t int) Setting {
	return func(s *Settings) {
		s.ExpireGracePeriod = optional.NewInt(t)
	}
}

// WithEnableEmoticons sets whether emoticons are enabled for an integration.
func WithEnableEmoticons(yes bool) Setting {
	return func(s *Settings) {
		s.EnableEmoticons = optional.NewBool(yes)
	}
}
