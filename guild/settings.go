package guild

import (
	"github.com/skwair/harmony/optional"
)

// Settings are the settings of a guild, all fields are optional and only those
// explicitly set will be modified.
type Settings struct {
	Name                        *optional.String `json:"name,omitempty"`
	Region                      *optional.String `json:"region,omitempty"`
	VerificationLevel           *optional.Int    `json:"verification_level,omitempty"`
	DefaultMessageNotifications *optional.Int    `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *optional.Int    `json:"explicit_content_filter,omitempty"`
	AfkChannelID                *optional.String `json:"afk_channel_id,omitempty"`
	AfkTimeout                  *optional.Int    `json:"afk_timeout,omitempty"`
	Icon                        *optional.String `json:"icon,omitempty"`
	OwnerID                     *optional.String `json:"owner_id,omitempty"`
	Splash                      *optional.String `json:"splash,omitempty"`
	SystemChannelID             *optional.String `json:"system_channel_id,omitempty"`
}

// Setting is a function that configures a guild.
type Setting func(*Settings)

// NewSettings returns new Settings to modify a a guild.
func NewSettings(opts ...Setting) *Settings {
	s := &Settings{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithName sets the name of a guild.
func WithName(name string) Setting {
	return func(s *Settings) {
		s.Name = optional.NewString(name)
	}
}

// WithRegion sets the region of a guild.
func WithRegion(region string) Setting {
	return func(s *Settings) {
		s.Region = optional.NewString(region)
	}
}

// VerificationLevel is the level of verification that applies on a server.
// Members must meet criteria before they can send messages or initiate
// direct message conversations with other guild members.
// This does not apply for members that have a role assigned to them.
type VerificationLevel int

const (
	// VerificationLevelNone means there is no verification.
	VerificationLevelNone VerificationLevel = iota
	// VerificationLevelLow means a member must have
	// a verified email on their account.
	VerificationLevelLow
	// VerificationLevelMedium means a member must be
	// registered on Discord for longer than 5 minutes.
	VerificationLevelMedium
	// VerificationLevelHigh means a member must be
	// in this server for longer than 10 minutes.
	VerificationLevelHigh
	// VerificationLevelVeryHigh means a member must have
	// a verified phone number.
	VerificationLevelVeryHigh
)

// WithVerificationLevel sets the verification level of a guild.
func WithVerificationLevel(lvl VerificationLevel) Setting {
	return func(s *Settings) {
		if lvl < 0 {
			lvl = 0
		}
		if lvl > 4 {
			lvl = 4
		}

		s.VerificationLevel = optional.NewInt(int(lvl))
	}
}

// DefaultNotificationLevel determines whether members who have not explicitly
// set their notification settings receive a notification for every message
// sent in this server or not.
type DefaultNotificationLevel int

const (
	// DefaultNotificationLevelAll means a notification
	// will be sent for all messages.
	DefaultNotificationLevelAll DefaultNotificationLevel = iota
	// DefaultNotificationLevelMentionOnly means a
	// notification will be sent for mentions only.
	DefaultNotificationLevelMentionOnly
)

// WithDefaultMessageNotifications sets the default notification level of a guild.
func WithDefaultMessageNotifications(lvl DefaultNotificationLevel) Setting {
	return func(s *Settings) {
		if lvl < 0 {
			lvl = 0
		}
		if lvl > 1 {
			lvl = 1
		}

		s.DefaultMessageNotifications = optional.NewInt(int(lvl))
	}
}

// ExplicitContentFilter determines how the explicit content filter
// should behave for a server.
type ExplicitContentFilter int

const (
	// ExplicitContentFilterDisabled disables the filter.
	ExplicitContentFilterDisabled ExplicitContentFilter = iota
	// ExplicitContentFilterWithoutRole filters messages from
	// members without a role.
	ExplicitContentFilterWithoutRole
	// ExplicitContentFilterAll filters messages from all members.
	ExplicitContentFilterAll
)

// WithExplicitContentFilter sets the explicit content filter of a guild.
func WithExplicitContentFilter(lvl ExplicitContentFilter) Setting {
	return func(s *Settings) {
		if lvl < 0 {
			lvl = 0
		}
		if lvl > 2 {
			lvl = 2
		}

		s.ExplicitContentFilter = optional.NewInt(int(lvl))
	}
}

// WithAfkChannel sets the AFK channel ID of a guild.
func WithAfkChannel(id string) Setting {
	return func(s *Settings) {
		s.AfkChannelID = optional.NewString(id)
	}
}

// WithAfkTimeout sets the AFK timeout of a guild.
func WithAfkTimeout(sec int) Setting {
	return func(s *Settings) {
		s.AfkTimeout = optional.NewInt(sec)
	}
}

// WithIcon sets the Guild icon which is a base64 encoded 128x128 jpeg image.
func WithIcon(icon string) Setting {
	return func(s *Settings) {
		s.Icon = optional.NewString(icon)
	}
}

// WithOwner sets the owner ID of a guild (must be the guild owner to for this to have effect).
func WithOwner(id string) Setting {
	return func(s *Settings) {
		s.OwnerID = optional.NewString(id)
	}
}

// WithSplash sets the Guild splash (VIP only) which is a base64 encoded 128x128 image.
func WithSplash(splash string) Setting {
	return func(s *Settings) {
		s.Splash = optional.NewString(splash)
	}
}

// WithSystemChannel sets the id of the channel to which system messages are sent.
func WithSystemChannel(id string) Setting {
	return func(s *Settings) {
		s.SystemChannelID = optional.NewString(id)
	}
}
