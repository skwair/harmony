package discord

// MessageType describes the type of a message. Different fields
// are set or not depending on the message's type.
type MessageType int

// Supported message types:
const (
	MessageTypeDefault                           MessageType = 0
	MessageTypeRecipientAdd                      MessageType = 1
	MessageTypeRecipientRemove                   MessageType = 2
	MessageTypeCall                              MessageType = 3
	MessageTypeChannelNameChange                 MessageType = 4
	MessageTypeChannelIconChange                 MessageType = 5
	MessageTypeChannelPinnedMessage              MessageType = 6
	MessageTypeGuildMemberJoin                   MessageType = 7
	MessageTypeUserPremiumGuildSubscription      MessageType = 8
	MessageTypeUserPremiumGuildSubscriptionTier1 MessageType = 9
	MessageTypeUserPremiumGuildSubscriptionTier2 MessageType = 10
	MessageTypeUserPremiumGuildSubscriptionTier3 MessageType = 11
	MessageTypeChannelFollowAdd                  MessageType = 12
	MessageTypeGuildDiscoveryDisqualified        MessageType = 14
	MessageTypeGuildDiscoveryRequaligied         MessageType = 15
)

// MessageFlag describes extra features a message can have.
type MessageFlag int

const (
	// This message has been published to subscribed channels (via Channel Following).
	MessageFlagCrossposted MessageFlag = 1 << 0
	// This message originated from a message in another channel (via Channel Following).
	MessageFlagIsCrosspost MessageFlag = 1 << 1
	// Do not include any embeds when serializing this message.
	MessageFlagSuppressEmbeds MessageFlag = 1 << 2
)

// Message represents a message sent in a channel within Discord.
// The author object follows the structure of the user object, but is
// only a valid user in the case where the message is generated by a
// user or bot user. If the message is generated by a webhook, the
// author object corresponds to the webhook's id, username, and avatar.
// You can tell if a message is generated by a webhook by checking for
// the webhook_id on the message object.
type Message struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
	// Author of this message. Can be a user or a webhook.
	// Check the WebhookID field to know.
	Author User `json:"author"`
	// Guild member info of the author that sent the message.
	// Only set for MESSAGE_CREATE and MESSAGE_UPDATE Gateway
	// events.
	Member          GuildMember `json:"member"`
	Content         string      `json:"content"`
	Timestamp       Time        `json:"timestamp"`
	EditedTimestamp Time        `json:"edited_timestamp"`
	TTS             bool        `json:"tts"`
	// MentionEveryone is set to true if '@everyone' or '@here'
	// is set in the message's content.
	MentionEveryone bool `json:"mention_everyone"`
	// Mentions contains an array of users that where mentioned
	// in the message's content.
	Mentions []User `json:"mentions"`
	// MentionRoles contains an array of IDs of te roles that
	// were mentioned in the message's content.
	MentionRoles []string `json:"mention_roles"`
	// Not all channel mentions in a message will appear in mention_channels.
	// Only textual channels that are visible to everyone in a public guild
	// will ever be included. Only crossposted messages (via Channel Following)
	// currently include mention_channels at all. If no mentions in the message
	// meet these requirements, this field will not be sent.
	MentionChannels []ChannelMention    `json:"mention_channels"`
	Attachments     []MessageAttachment `json:"attachments"` // Any attached files.
	Embeds          []MessageEmbed      `json:"embeds"`      // Any embedded content.
	Reactions       []MessageReaction   `json:"reactions"`
	Nonce           string              `json:"nonce"` // Used for validating a message was sent.
	Pinned          bool                `json:"pinned"`
	WebhookID       string              `json:"webhook_id"`
	Type            MessageType         `json:"type"`

	// Sent with Rich Presence-related chat embeds.
	Activity         MessageActivity    `json:"activity"`
	Application      MessageApplication `json:"application"`
	MessageReference MessageReference   `json:"message_reference"`
	Flags            MessageFlag        `json:"flags"`
}

// MessageAttachment is a file attached to a message.
type MessageAttachment struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Size     int    `json:"size"`
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

// MessageReference is a reference to an original message.
type MessageReference struct {
	MessageID string `json:"message_id"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
}

// MessageReaction is a reaction on a Discord message.
type MessageReaction struct {
	Count int   `json:"count"`
	Me    bool  `json:"me"`
	Emoji Emoji `json:"emoji"`
}

type MessageActivityType int

const (
	MessageActivityTypeJoin MessageActivityType = iota
	MessageActivityTypeSpectate
	MessageActivityTypeListen
	MessageActivityTypeJoinRequest
)

type MessageActivity struct {
	Type    MessageActivityType
	PartyID string
}

type MessageApplication struct {
	ID          string
	CoverImage  string
	Description string
	Icon        string
	Name        string
}
