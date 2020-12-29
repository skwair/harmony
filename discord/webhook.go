package discord

// Webhook is a low-effort way to post messages to channels in Discord.
// It do not require a bot user or authentication to use.
type Webhook struct {
	ID        string `json:"id"`
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	User      *User  `json:"user"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	Token     string `json:"token"`
}
