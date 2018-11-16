package guild

import "github.com/skwair/discord/optional"

// GuildMemberSettings are the settings of a guild member, all fields are optional
// and only those explicitly set will be modified.
type MemberSettings struct {
	Nick  *optional.String      `json:"nick,omitempty"`
	Roles *optional.StringSlice `json:"roles,omitempty"`
	Mute  *optional.Bool        `json:"mute,omitempty"`
	Deaf  *optional.Bool        `json:"deaf,omitempty"`
	// ID of channel to move user to (if they are connected to voice).
	ChannelID *optional.String `json:"channel_id,omitempty"`
}

// MemberSetting is a function that configures a guild member.
type MemberSetting func(*MemberSettings)

// NewMemberSettings returns new Settings to modify a a guild member.
func NewMemberSettings(opts ...MemberSetting) *MemberSettings {
	s := &MemberSettings{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithNick sets the name of a guild member.
func WithNick(name string) MemberSetting {
	return func(s *MemberSettings) {
		s.Nick = optional.NewString(name)
	}
}

// WithRoles sets the roles of a guild member.
func WithRoles(roleIDs []string) MemberSetting {
	return func(s *MemberSettings) {
		s.Roles = optional.NewStringSlice(roleIDs)
	}
}

// WithMute sets whether a guild member is muted.
func WithMute(yes bool) MemberSetting {
	return func(s *MemberSettings) {
		s.Mute = optional.NewBool(yes)
	}
}

// WithDeaf sets whether a guild member is deafen.
func WithDeaf(yes bool) MemberSetting {
	return func(s *MemberSettings) {
		s.Deaf = optional.NewBool(yes)
	}
}

// WithChannelID sets the channel id of a guild member (if connected to voice).
func WithChannelID(id string) MemberSetting {
	return func(s *MemberSettings) {
		s.ChannelID = optional.NewString(id)
	}
}
