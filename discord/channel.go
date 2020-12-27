package discord

// ChannelType describes the type of a channel. Different fields
// are set or not depending on the channel's type.
type ChannelType int

// Supported channel types:
const (
	ChannelTypeGuildText ChannelType = iota
	ChannelTypeDM
	ChannelTypeGuildVoice
	ChannelTypeGroupDM
	ChannelTypeGuildCategory
	ChannelTypeGuildNews
	ChannelTypeGuildStore
)

// ChannelUserRateLimit is the set of allowed values for Channel.RateLimitPerUser.
type ChannelUserRateLimit int

// Valid Channel User rate limits:
const (
	ChannelUserRateLimit5s  ChannelUserRateLimit = 5
	ChannelUserRateLimit10s ChannelUserRateLimit = 10
	ChannelUserRateLimit15s ChannelUserRateLimit = 15
	ChannelUserRateLimit30s ChannelUserRateLimit = 30
	ChannelUserRateLimit1m  ChannelUserRateLimit = 60
	ChannelUserRateLimit2m  ChannelUserRateLimit = 120
	ChannelUserRateLimit5m  ChannelUserRateLimit = 300
	ChannelUserRateLimit10m ChannelUserRateLimit = 600
	ChannelUserRateLimit15m ChannelUserRateLimit = 900
	ChannelUserRateLimit30m ChannelUserRateLimit = 1800
	ChannelUserRateLimit1h  ChannelUserRateLimit = 3600
	ChannelUserRateLimit2h  ChannelUserRateLimit = 7200
	ChannelUserRateLimit6h  ChannelUserRateLimit = 21600
)

// Channel represents a guild or DM channel within Discord.
type Channel struct {
	ID                   string                `json:"id"`
	Type                 ChannelType           `json:"type"`
	GuildID              string                `json:"guild_id"`
	Position             int                   `json:"position"` // Sorting position of the channel.
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	Name                 string                `json:"name"`
	Topic                string                `json:"topic"`
	NSFW                 bool                  `json:"nsfw"`
	LastMessageID        string                `json:"last_message_id"`
	ParentID             string                `json:"parent_id"` // ID of the parent category for a channel (only in guilds).
	LastPinTimestamp     Time                  `json:"last_pin_timestamp"`

	// For voice channels only.
	Bitrate          int                  `json:"bitrate"`
	UserLimit        int                  `json:"user_limit"`
	RateLimitPerUser ChannelUserRateLimit `json:"rate_limit_per_user"`

	// For DMs only.
	Recipients    []User `json:"recipients"`
	Icon          string `json:"icon"`
	OwnerID       string `json:"owner_id"`
	ApplicationID string `json:"application_id"` // Application id of the group DM creator if it is bot-created.
}

// ChannelMention represents a channel mention.
type ChannelMention struct {
	ID      string      `json:"id"`
	GuildID string      `json:"guild_id"`
	Type    ChannelType `json:"type"`
	Name    string      `json:"name"`
}
