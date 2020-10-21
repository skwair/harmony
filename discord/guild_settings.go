package discord

import "github.com/skwair/harmony/optional"

// GuildSettings are the settings of a guild, all fields are optional and only those
// explicitly set will be modified.
type GuildSettings struct {
	Name                        *optional.String `json:"name,omitempty"`
	Region                      *optional.String `json:"region,omitempty"`
	VerificationLevel           *optional.Int    `json:"verification_level,omitempty"`
	DefaultMessageNotifications *optional.Int    `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *optional.Int    `json:"explicit_content_filter,omitempty"`
	AFKChannelID                *optional.String `json:"afk_channel_id,omitempty"`
	AFKTimeout                  *optional.Int    `json:"afk_timeout,omitempty"`
	Icon                        *optional.String `json:"icon,omitempty"`
	OwnerID                     *optional.String `json:"owner_id,omitempty"`
	Splash                      *optional.String `json:"splash,omitempty"`
	SystemChannelID             *optional.String `json:"system_channel_id,omitempty"`
}

// GuildSetting is a function that configures a guild.
type GuildSetting func(*GuildSettings)

// NewGuildSettings returns new GuildSettings to modify a a guild.
func NewGuildSettings(opts ...GuildSetting) *GuildSettings {
	s := &GuildSettings{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithGuildName sets the name of a guild.
func WithGuildName(name string) GuildSetting {
	return func(s *GuildSettings) {
		s.Name = optional.NewString(name)
	}
}

// WithGuildRegion sets the region of a guild.
func WithGuildRegion(region string) GuildSetting {
	return func(s *GuildSettings) {
		s.Region = optional.NewString(region)
	}
}

// WithGuildVerificationLevel sets the verification level of a guild.
func WithGuildVerificationLevel(lvl GuildVerificationLevel) GuildSetting {
	return func(s *GuildSettings) {
		s.VerificationLevel = optional.NewInt(int(lvl))
	}
}

// WithGuildDefaultMessageNotifications sets the default notification level of a guild.
func WithGuildDefaultMessageNotifications(lvl GuildDefaultNotificationLevel) GuildSetting {
	return func(s *GuildSettings) {
		s.DefaultMessageNotifications = optional.NewInt(int(lvl))
	}
}

// WithGuildExplicitContentFilter sets the explicit content filter of a guild.
func WithGuildExplicitContentFilter(lvl GuildExplicitContentFilter) GuildSetting {
	return func(s *GuildSettings) {
		s.ExplicitContentFilter = optional.NewInt(int(lvl))
	}
}

// WithGuildAFKChannel sets the AFK channel ID of a guild.
// An empty id will disable the AFK channel.
func WithGuildAFKChannel(id string) GuildSetting {
	return func(s *GuildSettings) {
		if id == "" {
			s.AFKChannelID = optional.NewNilString()
		} else {
			s.AFKChannelID = optional.NewString(id)
		}
	}
}

// WithGuildAFKTimeout sets the AFK timeout of a guild.
func WithGuildAFKTimeout(t GuildAFKTimeout) GuildSetting {
	return func(s *GuildSettings) {
		s.AFKTimeout = optional.NewInt(int(t))
	}
}

// WithGuildIcon sets the Guild icon which is a base64 encoded 128x128 jpeg image.
func WithGuildIcon(icon string) GuildSetting {
	return func(s *GuildSettings) {
		s.Icon = optional.NewString(icon)
	}
}

// WithGuildOwner sets the owner ID of a guild (must be the guild owner to for this to have effect).
func WithGuildOwner(id string) GuildSetting {
	return func(s *GuildSettings) {
		s.OwnerID = optional.NewString(id)
	}
}

// WithGuildSplash sets the Guild splash (VIP only) which is a base64 encoded 128x128 image.
func WithGuildSplash(splash string) GuildSetting {
	return func(s *GuildSettings) {
		s.Splash = optional.NewString(splash)
	}
}

// WithGuildSystemChannel sets the id of the channel to which system messages are sent.
func WithGuildSystemChannel(id string) GuildSetting {
	return func(s *GuildSettings) {
		s.SystemChannelID = optional.NewString(id)
	}
}

// GuildMemberSettings are the settings of a guild member, all fields are optional
// and only those explicitly set will be modified.
type GuildMemberSettings struct {
	Nick  *optional.String      `json:"nick,omitempty"`
	Roles *optional.StringSlice `json:"roles,omitempty"`
	Mute  *optional.Bool        `json:"mute,omitempty"`
	Deaf  *optional.Bool        `json:"deaf,omitempty"`
	// ID of channel to move user to (if they are connected to voice).
	ChannelID *optional.String `json:"channel_id,omitempty"`
}

// GuildMemberSetting is a function that configures a guild member.
type GuildMemberSetting func(*GuildMemberSettings)

// NewGuildMemberSettings returns new GuildMemberSetting to modify a a guild member.
func NewGuildMemberSettings(opts ...GuildMemberSetting) *GuildMemberSettings {
	s := &GuildMemberSettings{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithGuildMemberNick sets the name of a guild member.
func WithGuildMemberNick(name string) GuildMemberSetting {
	return func(s *GuildMemberSettings) {
		s.Nick = optional.NewString(name)
	}
}

// WithGuildMemberRoles sets the roles of a guild member.
func WithGuildMemberRoles(roleIDs []string) GuildMemberSetting {
	return func(s *GuildMemberSettings) {
		s.Roles = optional.NewStringSlice(roleIDs)
	}
}

// WithGuildMemberMute sets whether a guild member is muted.
func WithGuildMemberMute(yes bool) GuildMemberSetting {
	return func(s *GuildMemberSettings) {
		s.Mute = optional.NewBool(yes)
	}
}

// WithGuildMemberDeaf sets whether a guild member is deafen.
func WithGuildMemberDeaf(yes bool) GuildMemberSetting {
	return func(s *GuildMemberSettings) {
		s.Deaf = optional.NewBool(yes)
	}
}

// WithGuildMemberChannelID sets the channel id of a guild member (if connected to voice).
func WithGuildMemberChannelID(id string) GuildMemberSetting {
	return func(s *GuildMemberSettings) {
		s.ChannelID = optional.NewString(id)
	}
}

// RoleSettings describes how to modify a guild role, all fields are optional
// and only those explicitly set will be modified.
type RoleSettings struct {
	Name        *optional.String `json:"name,omitempty"`
	Permissions *optional.Int    `json:"permissions,omitempty"`
	Color       *optional.Int    `json:"color,omitempty"`
	Hoist       *optional.Bool   `json:"hoist,omitempty"`
	Mentionable *optional.Bool   `json:"mentionable,omitempty"`
}

// RoleSetting is a function that configures a guild role.
type RoleSetting func(*RoleSettings)

// NewRoleSettings returns new RoleSettings to modify a a guild role.
func NewRoleSettings(opts ...RoleSetting) *RoleSettings {
	s := &RoleSettings{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithRoleName sets the name of guild a role.
func WithRoleName(name string) RoleSetting {
	return func(s *RoleSettings) {
		s.Name = optional.NewString(name)
	}
}

// WithRolePermissions sets the permissions of guild a role.
func WithRolePermissions(perm int) RoleSetting {
	return func(s *RoleSettings) {
		s.Permissions = optional.NewInt(perm)
	}
}

// WithRoleColor sets the color of guild a role. It accepts hexadecimal value.
func WithRoleColor(hexCode int) RoleSetting {
	return func(s *RoleSettings) {
		s.Color = optional.NewInt(hexCode)
	}
}

// WithRoleHoist sets whether this guild role is hoisted.
func WithRoleHoist(yes bool) RoleSetting {
	return func(s *RoleSettings) {
		s.Hoist = optional.NewBool(yes)
	}
}

// WithRoleMentionable sets whether this guild role is mentionable by others.
func WithRoleMentionable(yes bool) RoleSetting {
	return func(s *RoleSettings) {
		s.Mentionable = optional.NewBool(yes)
	}
}

// GuildIntegrationSettings describes a guild integration's settings.
type GuildIntegrationSettings struct {
	ExpireBehavior    *optional.Int  `json:"expire_behavior,omitempty"`
	ExpireGracePeriod *optional.Int  `json:"expire_grace_period,omitempty"`
	EnableEmoticons   *optional.Bool `json:"enable_emoticons,omitempty"`
}

// GuildIntegrationSetting is a function that configures a guild integration.
type GuildIntegrationSetting func(*GuildIntegrationSettings)

// NewGuildIntegrationSettings returns new GuildIntegrationSettings to modify a a guild integration.
func NewGuildIntegrationSettings(opts ...GuildIntegrationSetting) *GuildIntegrationSettings {
	s := &GuildIntegrationSettings{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithGuildIntegrationExpireBehavior sets the behavior when an integration subscription lapses.
func WithGuildIntegrationExpireBehavior(v int) GuildIntegrationSetting {
	return func(s *GuildIntegrationSettings) {
		s.ExpireBehavior = optional.NewInt(v)
	}
}

// WithGuildIntegrationExpireGracePeriod sets the period (in seconds) where the integration
// will ignore lapsed subscriptions.
func WithGuildIntegrationExpireGracePeriod(v int) GuildIntegrationSetting {
	return func(s *GuildIntegrationSettings) {
		s.ExpireGracePeriod = optional.NewInt(v)
	}
}

// WithGuildIntegrationEnableEmoticons sets whether emoticons should be synced for this
// integration (twitch only currently).
func WithGuildIntegrationEnableEmoticons(yes bool) GuildIntegrationSetting {
	return func(s *GuildIntegrationSettings) {
		s.EnableEmoticons = optional.NewBool(yes)
	}
}
