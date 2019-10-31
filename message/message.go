package message

// Type describes the type of a message. Different fields
// are set or not depending on the message's type.
type Type int

// Supported message types:
const (
	TypeDefault Type = iota
	TypeRecipientAdd
	TypeRecipientRemove
	TypeCall
	TypeChannelNameChange
	TypeChannelIconChange
	TypeChannelPinnedMessage
	TypeGuildMemberJoin
	TypeUserPremiumGuildSubscription
	TypeUserPremiumGuildSubscriptionTier1
	TypeUserPremiumGuildSubscriptionTier2
	TypeUserPremiumGuildSubscriptionTier3
	TypeChannelFollowAdd
)

// Flag describes extra features a message can have.
type Flag int

const (
	// This message has been published to subscribed channels (via Channel Following).
	FlagCrossposted Flag = 1 << 0
	// This message originated from a message in another channel (via Channel Following).
	FlagIsCrosspost Flag = 1 << 1
	// Do not include any embeds when serializing this message.
	FlagSuppressEmbeds Flag = 1 << 2
)

// Attachment is a file attached to a message.
type Attachment struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Size     int    `json:"size"`
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

// Reference is a reference to an original message.
type Reference struct {
	MessageID string `json:"message_id"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
}
