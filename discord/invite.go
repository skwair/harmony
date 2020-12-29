package discord

// Invite represents a code that when used, adds a user to a guild or group DM channel.
type Invite struct {
	Code                     string  `json:"code"`
	Guild                    Guild   `json:"guild"` // Nil if this invite is for a group DM channel.
	Channel                  Channel `json:"channel"`
	ApproximatePresenceCount int     `json:"approximate_presence_count"`
	ApproximateMemberCount   int     `json:"approximate_member_count"`

	InviteMetadata
}

// InviteMetadata contains additional information about an Invite.
type InviteMetadata struct {
	Inviter   User `json:"inviter"`
	Uses      int  `json:"uses"`
	MaxUses   int  `json:"max_uses"`
	MaxAge    int  `json:"max_age"`
	Temporary bool `json:"temporary"`
	CreatedAt Time `json:"created_at"`
	Revoked   bool `json:"revoked"`
}
