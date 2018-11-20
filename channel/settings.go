package channel

import (
	"github.com/skwair/harmony/optional"
	"github.com/skwair/harmony/permission"
)

// Type describes the type of a channel. Different fields
// are set or not depending on the channel's type.
type Type int

// Supported channel types :
const (
	TypeGuildText Type = iota
	TypeDM
	TypeGuildVoice
	TypeGroupDM
	TypeGuildCategory
)

// Settings describes a channel creation.
type Settings struct {
	Name      *optional.String `json:"name,omitempty"` // 2-100 characters.
	Type      *optional.Int    `json:"type,omitempty"`
	Topic     *optional.String `json:"topic"` // 0-1000 characters.
	Bitrate   *optional.Int    `json:"bitrate,omitempty"`
	UserLimit *optional.Int    `json:"user_limit,omitempty"`
	// RateLimitPerUser is the amount of seconds a user has to wait before sending
	// another message (0-120); bots, as well as users with the permission
	// 'manage_messages' or 'manage_channel', are unaffected.
	RateLimitPerUser *optional.Int `json:"rate_limit_per_user"`
	// Sorting position of the channel.
	Position    *optional.Int          `json:"position"`
	Permissions []permission.Overwrite `json:"permission_overwrites,omitempty"`
	ParentID    *optional.String       `json:"parent_id,omitempty"`
	NSFW        *optional.Bool         `json:"nsfw,omitempty"`
}

// Setting is a function that configures a channel.
type Setting func(*Settings)

// NewSettings returns new Settings to modify a a channel.
func NewSettings(opts ...Setting) *Settings {
	s := &Settings{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithName sets the name of a channel.
func WithName(name string) Setting {
	return func(s *Settings) {
		s.Name = optional.NewString(name)
	}
}

// WithType sets the type of a channel. Can only be used when creating channels.
func WithType(typ Type) Setting {
	return func(s *Settings) {
		s.Type = optional.NewInt(int(typ))
	}
}

// WithTopic sets the topic of a channel (text only).
func WithTopic(topic string) Setting {
	return func(s *Settings) {
		s.Topic = optional.NewString(topic)
	}
}

// WithBitrate sets the bit rate of a channel (audio only).
func WithBitrate(bitrate int) Setting {
	return func(s *Settings) {
		s.Bitrate = optional.NewInt(bitrate)
	}
}

// WithUserLimit sets the user limit of a channel (audio only).
func WithUserLimit(limit int) Setting {
	return func(s *Settings) {
		s.UserLimit = optional.NewInt(limit)
	}
}

// WithRateLimitPerUser sets the rate limit per user (text only).
func WithRateLimitPerUser(rateLimit int) Setting {
	return func(s *Settings) {
		s.RateLimitPerUser = optional.NewInt(rateLimit)
	}
}

// WithPosition sets the position of a channel.
func WithPosition(pos int) Setting {
	return func(s *Settings) {
		s.Position = optional.NewInt(pos)
	}
}

// WithPermissions sets the permission overwrites of a channel.
// Pass an empty array to remove all permission overwrites.
func WithPermissions(perms []permission.Overwrite) Setting {
	return func(s *Settings) {
		s.Permissions = perms
	}
}

// WithParent sets the parent ID channel of a channel.
func WithParent(id string) Setting {
	return func(s *Settings) {
		s.ParentID = optional.NewString(id)
	}
}

// WithNSFW sets whether a channel is not safe for work.
func WithNSFW(yes bool) Setting {
	return func(s *Settings) {
		s.NSFW = optional.NewBool(yes)
	}
}
