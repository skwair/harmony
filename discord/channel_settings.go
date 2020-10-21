package discord

import "github.com/skwair/harmony/optional"

// ChannelSettings describes a channel creation or update.
type ChannelSettings struct {
	Name      *optional.String `json:"name,omitempty"` // 2-100 characters.
	Type      *optional.Int    `json:"type,omitempty"`
	Topic     *optional.String `json:"topic,omitempty"` // 0-1000 characters.
	Bitrate   *optional.Int    `json:"bitrate,omitempty"`
	UserLimit *optional.Int    `json:"user_limit,omitempty"`
	// RateLimitPerUser is the amount of seconds a user has to wait before sending
	// another message (0-120); bots, as well as users with the permission
	// 'manage_messages' or 'manage_channel', are unaffected.
	RateLimitPerUser *optional.Int `json:"rate_limit_per_user,omitempty"`
	// Sorting position of the channel.
	Position    *optional.Int         `json:"position,omitempty"`
	Permissions []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID    *optional.String      `json:"parent_id,omitempty"`
	NSFW        *optional.Bool        `json:"nsfw,omitempty"`
}

// ChannelSetting is a function that configures a channel.
type ChannelSetting func(*ChannelSettings)

// NewChannelSettings returns new ChannelSettings to modify a a channel.
func NewChannelSettings(opts ...ChannelSetting) *ChannelSettings {
	s := &ChannelSettings{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithChannelName sets the name of a channel.
func WithChannelName(name string) ChannelSetting {
	return func(s *ChannelSettings) {
		s.Name = optional.NewString(name)
	}
}

// WithChannelType sets the name of a channel.
func WithChannelType(t ChannelType) ChannelSetting {
	return func(s *ChannelSettings) {
		s.Type = optional.NewInt(int(t))
	}
}

// WithChannelTopic sets the topic of a channel (text only).
func WithChannelTopic(topic string) ChannelSetting {
	return func(s *ChannelSettings) {
		s.Topic = optional.NewString(topic)
	}
}

// WithChannelBitrate sets the bit rate of a channel (audio only).
// Must be a value between 8 and 96 for regular, non-premium channels
// and can go up to 256 for Tier 2 guild and 384 for Tier 3 guilds.
func WithChannelBitrate(bitrate int) ChannelSetting {
	return func(s *ChannelSettings) {
		s.Bitrate = optional.NewInt(bitrate)
	}
}

// WithChannelUserLimit sets the user limit of a channel (audio only).
func WithChannelUserLimit(limit int) ChannelSetting {
	return func(s *ChannelSettings) {
		s.UserLimit = optional.NewInt(limit)
	}
}

// WithChannelRateLimitPerUser sets the rate limit per user (text only).
func WithChannelRateLimitPerUser(rateLimit ChannelUserRateLimit) ChannelSetting {
	return func(s *ChannelSettings) {
		s.RateLimitPerUser = optional.NewInt(int(rateLimit))
	}
}

// WithChannelPosition sets the position of a channel.
func WithChannelPosition(pos int) ChannelSetting {
	return func(s *ChannelSettings) {
		s.Position = optional.NewInt(pos)
	}
}

// WithChannelPermissions sets the permission overwrites of a channel.
// Pass an empty array to remove all permission overwrites.
func WithChannelPermissions(perms []PermissionOverwrite) ChannelSetting {
	return func(s *ChannelSettings) {
		s.Permissions = perms
	}
}

// WithChannelParent sets the parent ID channel of a channel.
func WithChannelParent(id string) ChannelSetting {
	return func(s *ChannelSettings) {
		s.ParentID = optional.NewString(id)
	}
}

// WithChannelNSFW sets whether a channel is not safe for work.
func WithChannelNSFW(yes bool) ChannelSetting {
	return func(s *ChannelSettings) {
		s.NSFW = optional.NewBool(yes)
	}
}
