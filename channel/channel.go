package channel

// Type describes the type of a channel. Different fields
// are set or not depending on the channel's type.
type Type int

// Supported channel types:
const (
	TypeGuildText Type = iota
	TypeDM
	TypeGuildVoice
	TypeGroupDM
	TypeGuildCategory
	TypeGuildNews
	TypeGuildStore
)

// Mention represents a channel mention.
type Mention struct {
	ID      string `json:"id"`
	GuildID string `json:"guild_id"`
	Type    Type   `json:"type"`
	Name    string `json:"name"`
}
