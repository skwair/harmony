package discord

import (
	"time"

	"github.com/skwair/harmony/voice"
)

// GuildVerificationLevel is the level of verification that applies on a guild.
// Members must meet criteria before they can send messages or initiate
// direct message conversations with other guild members.
// This does not apply for members that have a role assigned to them.
type GuildVerificationLevel int

const (
	// GuildVerificationLevelNone means there is no verification.
	GuildVerificationLevelNone GuildVerificationLevel = iota
	// GuildVerificationLevelLow means a member must have
	// a verified email on their account.
	GuildVerificationLevelLow
	// GuildVerificationLevelMedium means a member must be
	// registered on Discord for longer than 5 minutes.
	GuildVerificationLevelMedium
	// GuildVerificationLevelHigh means a member must be
	// in this server for longer than 10 minutes.
	GuildVerificationLevelHigh
	// GuildVerificationLevelVeryHigh means a member must have
	// a verified phone number.
	GuildVerificationLevelVeryHigh
)

// GuildExplicitContentFilter determines how the explicit content filter
// should behave for a server.
type GuildExplicitContentFilter int

const (
	// GuildExplicitContentFilterDisabled disables the filter.
	GuildExplicitContentFilterDisabled GuildExplicitContentFilter = iota
	// GuildExplicitContentFilterWithoutRole filters messages from
	// members without a role.
	GuildExplicitContentFilterWithoutRole
	// GuildExplicitContentFilterAll filters messages from all members.
	GuildExplicitContentFilterAll
)

// GuildAFKTimeout is the set of allowed values for AFK timeouts.
type GuildAFKTimeout int

// Valid Guild AFK timeouts:
const (
	GuildAFKTimeout1m  GuildAFKTimeout = 60
	GuildAFKTimeout5m  GuildAFKTimeout = 300
	GuildAFKTimeout15m GuildAFKTimeout = 900
	GuildAFKTimeout30m GuildAFKTimeout = 1800
	GuildAFKTimeout1h  GuildAFKTimeout = 3600
)

// GuildDefaultNotificationLevel determines whether members who have not explicitly
// set their notification settings receive a notification for every message
// sent in this server or not.
type GuildDefaultNotificationLevel int

const (
	// GuildDefaultNotificationLevelAll means a notification
	// will be sent for all messages.
	GuildDefaultNotificationLevelAll GuildDefaultNotificationLevel = iota
	// GuildDefaultNotificationLevelMentionOnly means a
	// notification will be sent for mentions only.
	GuildDefaultNotificationLevelMentionOnly
)

// Guild in Discord represents an isolated collection of users and channels,
// and are often referred to as "servers" in the UI.
type Guild struct {
	ID                          string                        `json:"id"`
	Name                        string                        `json:"name"`
	Icon                        string                        `json:"icon"`
	Splash                      string                        `json:"splash"`
	Owner                       bool                          `json:"owner"`
	OwnerID                     string                        `json:"owner_id"`
	Permissions                 int                           `json:"permissions"`
	Region                      string                        `json:"region"`
	AFKChannelID                string                        `json:"afk_channel_id"`
	AFKTimeout                  GuildAFKTimeout               `json:"afk_timeout"`
	VerificationLevel           GuildVerificationLevel        `json:"verification_level"`
	DefaultMessageNotifications GuildDefaultNotificationLevel `json:"default_message_notifications"`
	ExplicitContentFilter       GuildExplicitContentFilter    `json:"explicit_content_filter"`
	Roles                       []Role                        `json:"roles"`
	Emojis                      []Emoji                       `json:"emojis"`
	Features                    []string                      `json:"features"`
	MFALevel                    int                           `json:"mfa_level"`
	ApplicationID               string                        `json:"application_id"`
	WidgetEnabled               bool                          `json:"widget_enabled"`
	WidgetChannelID             string                        `json:"widget_channel_id"`
	SystemChannelID             string                        `json:"system_channel_id"`

	// Following fields are only sent within the GUILD_CREATE event.
	JoinedAt    time.Time     `json:"joined_at"`
	Large       bool          `json:"large"`
	Unavailable bool          `json:"unavailable"`
	MemberCount int           `json:"member_count"`
	VoiceStates []voice.State `json:"voice_states"`
	Members     []GuildMember `json:"members"`
	Channels    []Channel     `json:"channels"`
	Presences   []Presence    `json:"presences"`
}

// PartialGuild is a subset of the Guild object, returned by the Discord API
// when fetching current user's guilds.
type PartialGuild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Owner       bool   `json:"owner"`
	Permissions int    `json:"permissions"`
}

// UnavailableGuild is a Guild that is not available, either because there is a
// guild outage or because the connected user was removed from this guild.
type UnavailableGuild struct {
	ID          string `json:"id"`
	Unavailable *bool  `json:"unavailable"` // If not set, the connected user was removed from this Guild.
}

// GuildMember represents a User in a Guild.
// The field User won't be set in objects attached to MESSAGE_CREATE and MESSAGE_UPDATE gateway events.
type GuildMember struct {
	User         *User     `json:"user"`
	Nick         string    `json:"nick"`
	Roles        []string  `json:"roles"` // Role IDs.
	JoinedAt     time.Time `json:"joined_at"`
	PremiumSince time.Time `json:"premium_since"`
	Deaf         bool      `json:"deaf"`
	Mute         bool      `json:"mute"`
}

// PermissionsIn returns the permissions of the Guild member in the given Guild and channel.
func (m *GuildMember) PermissionsIn(g *Guild, ch *Channel) (permissions int) {
	base := computeBasePermissions(g, m)
	return computeOverwrites(ch, m, base)
}

// HasRole returns whether this member has the given role.
// Note that this method does not try to fetch this member latest roles, it instead looks
// in the roles it already had when this member object was created.
func (m *GuildMember) HasRole(id string) bool {
	for _, roleID := range m.Roles {
		if roleID == id {
			return true
		}
	}
	return false
}

// Role represents a set of permissions attached to a group of users.
// Roles have unique names, colors, and can be "pinned" to the side bar,
// causing their members to be listed separately. Roles are unique per guild,
// and can have separate permission profiles for the global context (guild)
// and channel context.
type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Color       int    `json:"color"`    // Integer representation of hexadecimal color code.
	Hoist       bool   `json:"hoist"`    // Whether this role is pinned in the user listing.
	Position    int    `json:"position"` // Integer	position of this role.
	Permissions int    `json:"permissions"`
	Managed     bool   `json:"managed"` // Whether this role is managed by an integration.
	Mentionable bool   `json:"mentionable"`
}

// Emoji represents a Discord emoji (both standard and custom).
type Emoji struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Roles         []Role `json:"roles"`
	User          *User  `json:"user"` // The user that created this emoji.
	RequireColons bool   `json:"require_colons"`
	Managed       bool   `json:"managed"`
	Animated      bool   `json:"animated"`
	// Whether this emoji can be used, may be false due to loss of Server Boosts.
	Available bool `json:"available"`
}

// Presence is a user's current state on a guild.
// This event is sent when a user's presence is updated for a guild.
type Presence struct {
	User    *User     `json:"user"`
	Roles   []string  `json:"roles"` // Array of IDs.
	Game    *Activity `json:"game"`
	GuildID string    `json:"guild_id"`
	Status  string    `json:"status"` // Either "idle", "dnd", "online", or "offline".
}

// VoiceRegion represents a voice region a guild can use or is using for its voice channels.
type VoiceRegion struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Whether this is a vip-only server.
	VIP bool `json:"vip"`
	// Whether this is a single server that is closest to the current user's client.
	Optimal bool `json:"optimal"`
	// Whether this is a deprecated voice region (avoid switching to these.
	Deprecated bool `json:"deprecated"`
	// Whether this is a custom voice region (used for events/etc).
	Custom bool `json:"custom"`
}

// Ban represents a Guild ban.
type Ban struct {
	Reason string
	User   *User
}

type GuildIntegration struct {
	ID                string                  `json:"id"`
	Name              string                  `json:"name"`
	Type              string                  `json:"type"`
	Enabled           bool                    `json:"enabled"`
	Syncing           bool                    `json:"syncing"`
	RoleID            string                  `json:"role_id"`
	ExpireBehavior    int                     `json:"expire_behavior"`
	ExpireGracePeriod int                     `json:"expire_grace_period"`
	User              User                    `json:"user"`
	Account           GuildIntegrationAccount `json:"account"`
	SyncedAt          time.Time               `json:"synced_at"`
}

type GuildIntegrationAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GuildWidget struct {
	Enabled   bool   `json:"enabled"`
	ChannelID string `json:"channel_id"`
}
