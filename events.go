package harmony

import "time"

type handler interface {
	handle(interface{})
}

func (c *Client) registerHandler(event string, h handler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if h == nil {
		panic("harmony: nil handler")
	}
	c.handlersMu.Lock()
	c.handlers[event] = h
	c.handlersMu.Unlock()
}

type readyHandler func(*Ready)

func (h readyHandler) handle(v interface{}) {
	h(v.(*Ready))
}

// HandleReady registers the handler function for the "READY" event.
func (c *Client) HandleReady(f func(r *Ready)) {
	c.registerHandler(eventReady, readyHandler(f))
}

type channelCreateHandler func(*Channel)

func (h channelCreateHandler) handle(v interface{}) {
	h(v.(*Channel))
}

// HandleChannelCreate registers the handler function for the "CHANNEL_CREATE" event.
// This event is fired when a new channel is created, relevant to the current user.
func (c *Client) HandleChannelCreate(f func(c *Channel)) {
	c.registerHandler(eventChannelCreate, channelCreateHandler(f))
}

type channelUpdateHandler func(*Channel)

func (h channelUpdateHandler) handle(v interface{}) {
	h(v.(*Channel))
}

// HandleChannelUpdate registers the handler function for the "CHANNEL_UPDATE" event.
// This event is fired when a channel is updated, relevant to the current user.
func (c *Client) HandleChannelUpdate(f func(c *Channel)) {
	c.registerHandler(eventChannelUpdate, channelUpdateHandler(f))
}

type channelDeleteHandler func(*Channel)

func (h channelDeleteHandler) handle(v interface{}) {
	h(v.(*Channel))
}

// HandleChannelDelete registers the handler function for the "CHANNEL_DELETE" event.
// This event is fired when a channel is deleted, relevant to the current user.
func (c *Client) HandleChannelDelete(f func(c *Channel)) {
	c.registerHandler(eventChannelDelete, channelDeleteHandler(f))
}

// ChannelPinsUpdate is Fired when a message is pinned or unpinned in a text channel.
type ChannelPinsUpdate struct {
	ChannelID        string    `json:"channel_id"`
	LastPinTimestamp time.Time `json:"last_pin_timestamp"`
}

type channelPinsUpdateHandler func(*ChannelPinsUpdate)

func (h channelPinsUpdateHandler) handle(v interface{}) {
	h(v.(*ChannelPinsUpdate))
}

// HandleChannelPinsUpdate registers the handler function for the "CHANNEL_PINS_UPDATE" event.
// This event is fired when a message is pinned or unpinned, but not when a pinned message
// is deleted.
func (c *Client) HandleChannelPinsUpdate(f func(cpu *ChannelPinsUpdate)) {
	c.registerHandler(eventChannelPinsUpdate, channelPinsUpdateHandler(f))
}

type guildCreateHandler func(*Guild)

func (h guildCreateHandler) handle(v interface{}) {
	h(v.(*Guild))
}

// HandleGuildCreate registers the handler function for the "GUILD_CREATE" event.
// This event can be sent in three different scenarios :
// 	1. When a user is initially connecting, to lazily load and backfill information for all unavailable guilds sent in the Ready event.
// 	2. When a Guild becomes available again to the client.
// 	3. When the current user joins a new Guild.
func (c *Client) HandleGuildCreate(f func(g *Guild)) {
	c.registerHandler(eventGuildCreate, guildCreateHandler(f))
}

type guildUpdateHandler func(*Guild)

func (h guildUpdateHandler) handle(v interface{}) {
	h(v.(*Guild))
}

// HandleGuildUpdate registers the handler function for the "GUILD_UPDATE" event.
func (c *Client) HandleGuildUpdate(f func(g *Guild)) {
	c.registerHandler(eventGuildUpdate, guildUpdateHandler(f))
}

type guildDeleteHandler func(*UnavailableGuild)

func (h guildDeleteHandler) handle(v interface{}) {
	h(v.(*UnavailableGuild))
}

// HandleGuildDelete registers the handler function for the "GUILD_DELETE" event.
// This event is fired when a guild becomes unavailable during a guild outage,
// or when the user leaves or is removed from a guild. If the unavailable field
// is not set, the user was removed from the guild.
func (c *Client) HandleGuildDelete(f func(g *UnavailableGuild)) {
	c.registerHandler(eventGuildDelete, guildDeleteHandler(f))
}

type GuildBan struct {
	*User
	GuildID string `json:"guild_id"`
}

type guildBanAddHandler func(*GuildBan)

func (h guildBanAddHandler) handle(v interface{}) {
	h(v.(*GuildBan))
}

// HandleGuildBanAdd registers the handler function for the "GUILD_BAN_ADD" event.
func (c *Client) HandleGuildBanAdd(f func(ban *GuildBan)) {
	c.registerHandler(eventGuildBanAdd, guildBanAddHandler(f))
}

type guildBanRemoveHandler func(*GuildBan)

func (h guildBanRemoveHandler) handle(v interface{}) {
	h(v.(*GuildBan))
}

// HandleGuildBanRemove registers the handler function for the "GUILD_BAN_REMOVE" event.
// This event is fired when a guild is updated.
func (c *Client) HandleGuildBanRemove(f func(ban *GuildBan)) {
	c.registerHandler(eventGuildBanRemove, guildBanRemoveHandler(f))
}

type GuildEmojis struct {
	Emojis  []Emoji `json:"emojis"`
	GuildID string  `json:"guild_id"`
}

type guildEmojisUpdateHandler func(*GuildEmojis)

func (h guildEmojisUpdateHandler) handle(v interface{}) {
	h(v.(*GuildEmojis))
}

// HandleGuildEmojisUpdate registers the handler function for the "GUILD_EMOJIS_UPDATE" event.
// Fired when a guild's emojis have been updated.
func (c *Client) HandleGuildEmojisUpdate(f func(emojis *GuildEmojis)) {
	c.registerHandler(eventGuildEmojisUpdate, guildEmojisUpdateHandler(f))
}

type guildIntegrationUpdateHandler func(string)

func (h guildIntegrationUpdateHandler) handle(v interface{}) {
	h(v.(string))
}

// HandleGuildIntegrationsUpdate registers the handler function for the "GUILD_INTEGRATIONS_UPDATE" event.
// Fired when a guild integration is updated.
func (c *Client) HandleGuildIntegrationsUpdate(f func(guildID string)) {
	c.registerHandler(eventGuildIntegrationsUpdate, guildIntegrationUpdateHandler(f))
}

type GuildMemberAdd struct {
	*GuildMember
	GuildID string `json:"guild_id"`
}

type guildMemberAddHandler func(*GuildMemberAdd)

func (h guildMemberAddHandler) handle(v interface{}) {
	h(v.(*GuildMemberAdd))
}

// HandleGuildMemberAdd registers the handler function for the "GUILD_MEMBER_ADD" event.
// Fired when a new user joins a guild.
func (c *Client) HandleGuildMemberAdd(f func(m *GuildMemberAdd)) {
	c.registerHandler(eventGuildMemberAdd, guildMemberAddHandler(f))
}

type GuildMemberRemove struct {
	User    *User  `json:"user"`
	GuildID string `json:"guild_id"`
}

type guildMemberRemoveHandler func(*GuildMemberRemove)

func (h guildMemberRemoveHandler) handle(v interface{}) {
	h(v.(*GuildMemberRemove))
}

// HandleGuildMemberRemove registers the handler function for the "GUILD_MEMBER_REMOVE" event.
// Fired when a user is removed from a guild (leave/kick/ban).
func (c *Client) HandleGuildMemberRemove(f func(m *GuildMemberRemove)) {
	c.registerHandler(eventGuildMemberRemove, guildMemberRemoveHandler(f))
}

type GuildMemberUpdate struct {
	GuildID string   `json:"guild_id"`
	Roles   []string `json:"roles"`
	User    *User    `json:"user"`
	Nick    string   `json:"nick"`
}

type guildMemberUpdateHandler func(*GuildMemberUpdate)

func (h guildMemberUpdateHandler) handle(v interface{}) {
	h(v.(*GuildMemberUpdate))
}

// HandleGuildMemberUpdate registers the handler function for the "GUILD_MEMBER_UPDATE" event.
// Fired when a guild member is updated.
func (c *Client) HandleGuildMemberUpdate(f func(m *GuildMemberUpdate)) {
	c.registerHandler(eventGuildMemberUpdate, guildMemberUpdateHandler(f))
}

type GuildMembersChunk struct {
	GuildID string        `json:"guild_id"`
	Members []GuildMember `json:"members"`
}

type guildMembersChunkHandler func(*GuildMembersChunk)

func (h guildMembersChunkHandler) handle(v interface{}) {
	h(v.(*GuildMembersChunk))
}

// HandleGuildMembersChunk registers the handler function for the "GUILD_MEMBERS_CHUNK" event.
// Sent in response to Guild Request Members.
func (c *Client) HandleGuildMembersChunk(f func(m *GuildMembersChunk)) {
	c.registerHandler(eventGuildMembersChunk, guildMembersChunkHandler(f))
}

type GuildRole struct {
	GuildID string `json:"guild_id"`
	Role    *Role  `json:"role"`
}

type guildRoleCreateHandler func(*GuildRole)

func (h guildRoleCreateHandler) handle(v interface{}) {
	h(v.(*GuildRole))
}

// HandleGuildRoleCreate registers the handler function for the "GUILD_ROLE_CREATE" event.
// Fired when a guild role is created.
func (c *Client) HandleGuildRoleCreate(f func(r *GuildRole)) {
	c.registerHandler(eventGuildRoleCreate, guildRoleCreateHandler(f))
}

type guildRoleUpdateHandler func(*GuildRole)

func (h guildRoleUpdateHandler) handle(v interface{}) {
	h(v.(*GuildRole))
}

// HandleGuildRoleUpdate registers the handler function for the "GUILD_ROLE_UPDATE" event.
// Fired when a guild role is updated.
func (c *Client) HandleGuildRoleUpdate(f func(r *GuildRole)) {
	c.registerHandler(eventGuildRoleUpdate, guildRoleUpdateHandler(f))
}

type GuildRoleDelete struct {
	GuildID string `json:"guild_id"`
	RoleID  string `json:"role_id"`
}

type guildRoleDeleteHandler func(*GuildRoleDelete)

func (h guildRoleDeleteHandler) handle(v interface{}) {
	h(v.(*GuildRoleDelete))
}

// HandleGuildRoleDelete registers the handler function for the "GUILD_ROLE_DELETE" event.
// Fired when a guild role is deleted.
func (c *Client) HandleGuildRoleDelete(f func(r *GuildRoleDelete)) {
	c.registerHandler(eventGuildRoleDelete, guildRoleDeleteHandler(f))
}

type messageCreateHandler func(*Message)

func (h messageCreateHandler) handle(v interface{}) {
	h(v.(*Message))
}

// HandleMessageCreate registers the handler function for the "MESSAGE_CREATE" event.
// Fired when a message is created.
func (c *Client) HandleMessageCreate(f func(m *Message)) {
	c.registerHandler(eventMessageCreate, messageCreateHandler(f))
}

type messageUpdateHandler func(*Message)

func (h messageUpdateHandler) handle(v interface{}) {
	h(v.(*Message))
}

// HandleMessageUpdate registers the handler function for the "MESSAGE_UPDATE" event.
// Fired when a message is updated. Unlike creates, message updates may contain only
// a subset of the full message object payload (but will always contain an id and channel_id).
func (c *Client) HandleMessageUpdate(f func(m *Message)) {
	c.registerHandler(eventMessageUpdate, messageUpdateHandler(f))
}

type MessageDelete struct {
	ChannelID string `json:"channel_id"`
	MessageID string `json:"id"`
}

type messageDeleteHandler func(*MessageDelete)

func (h messageDeleteHandler) handle(v interface{}) {
	h(v.(*MessageDelete))
}

// HandleMessageDelete registers the handler function for the "MESSAGE_DELETE" event.
// Fired when a message is deleted.
func (c *Client) HandleMessageDelete(f func(m *MessageDelete)) {
	c.registerHandler(eventMessageDelete, messageDeleteHandler(f))
}

type MessageDeleteBulk struct {
	GuildID   string   `json:"guild_id"`
	ChannelID string   `json:"channel_id"`
	IDs       []string `json:"ids"`
}

type messageDeleteBulkHandler func(*MessageDeleteBulk)

func (h messageDeleteBulkHandler) handle(v interface{}) {
	h(v.(*MessageDeleteBulk))
}

// HandleMessageDeleteBulk registers the handler function for the "MESSAGE_DELETE_BULK" event.
// Fired when multiple messages are deleted at once.
func (c *Client) HandleMessageDeleteBulk(f func(mdb *MessageDeleteBulk)) {
	c.registerHandler(eventMessageDeleteBulk, messageDeleteBulkHandler(f))
}

type MessageAck struct {
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
}

type messageAckHandler func(*MessageAck)

func (h messageAckHandler) handle(v interface{}) {
	h(v.(*MessageAck))
}

// HandleMessageAck registers the handler function for the "MESSAGE_ACK" event.
func (c *Client) HandleMessageAck(f func(ack *MessageAck)) {
	c.registerHandler(eventMessageAck, messageAckHandler(f))
}

type MessageReaction struct {
	UserID    string `json:"user_id"`
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
	Emoji     *Emoji `json:"emoji"`
}

type messageReactionAddHandler func(*MessageReaction)

func (h messageReactionAddHandler) handle(v interface{}) {
	h(v.(*MessageReaction))
}

// HandleMessageReactionAdd registers the handler function for the "MESSAGE_REACTION_ADD" event.
// Fired when a user adds a reaction to a message.
func (c *Client) HandleMessageReactionAdd(f func(r *MessageReaction)) {
	c.registerHandler(eventMessageReactionAdd, messageReactionAddHandler(f))
}

type messageReactionRemoveHandler func(*MessageReaction)

func (h messageReactionRemoveHandler) handle(v interface{}) {
	h(v.(*MessageReaction))
}

// HandleMessageReactionRemove registers the handler function for the "MESSAGE_REACTION_REMOVE" event.
// Fired when a user removes a reaction from a message.
func (c *Client) HandleMessageReactionRemove(f func(r *MessageReaction)) {
	c.registerHandler(eventMessageReactionRemove, messageReactionRemoveHandler(f))
}

type MessageReactionRemoveAll struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
}

type messageReactionRemoveAllHandler func(*MessageReactionRemoveAll)

func (h messageReactionRemoveAllHandler) handle(v interface{}) {
	h(v.(*MessageReactionRemoveAll))
}

// HandleMessageReactionRemoveAll registers the handler function for the "MESSAGE_REACTION_REMOVE_ALL" event.
// Fired when a user explicitly removes all reactions from a message.
func (c *Client) HandleMessageReactionRemoveAll(f func(r *MessageReactionRemoveAll)) {
	c.registerHandler(eventMessageReactionRemoveAll, messageReactionRemoveAllHandler(f))
}

type presenceUpdateHandler func(*Presence)

func (h presenceUpdateHandler) handle(v interface{}) {
	h(v.(*Presence))
}

// HandlePresenceUpdate registers the handler function for the "PRESENCE_UPDATE" event.
// This event is fired when a user's presence is updated for a guild.
// The user object within this event can be partial, the only field which must be sent
// is the id field, everything else is optional. Along with this limitation, no fields
// are required, and the types of the fields are not validated. Your client should expect
// any combination of fields and types within this event.
func (c *Client) HandlePresenceUpdate(f func(p *Presence)) {
	c.registerHandler(eventPresenceUpdate, presenceUpdateHandler(f))
}

type TypingStart struct {
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
	UserID    string `json:"user_id"`
	Timestamp int64  `json:"timestamp"`
}

type typingStartHandler func(*TypingStart)

func (h typingStartHandler) handle(v interface{}) {
	h(v.(*TypingStart))
}

// HandleTypingStart registers the handler function for the "TYPING_START" event.
// Fired when a user starts typing in a channel.
func (c *Client) HandleTypingStart(f func(ts *TypingStart)) {
	c.registerHandler(eventTypingStart, typingStartHandler(f))
}

type userUpdateHandler func(*User)

func (h userUpdateHandler) handle(v interface{}) {
	h(v.(*User))
}

// HandleUserUpdate registers the handler function for the "USER_UPDATE" event.
// Fired when properties about the user change.
func (c *Client) HandleUserUpdate(f func(u *User)) {
	c.registerHandler(eventUserUpdate, userUpdateHandler(f))
}

type voiceStateUpdateHandler func(*VoiceState)

func (h voiceStateUpdateHandler) handle(v interface{}) {
	h(v.(*VoiceState))
}

// HandleVoiceStateUpdate registers the handler function for the "VOICE_STATE_UPDATE" event.
// Fired when someone joins/leaves/moves voice channels.
func (c *Client) HandleVoiceStateUpdate(f func(vs *VoiceState)) {
	c.registerHandler(eventVoiceStateUpdate, voiceStateUpdateHandler(f))
}

type VoiceServerUpdate struct {
	Token    string `json:"token"`
	GuildID  string `json:"guild_id"`
	Endpoint string `json:"endpoint"`
}

type voiceServerUpdateHandler func(*VoiceServerUpdate)

func (h voiceServerUpdateHandler) handle(v interface{}) {
	h(v.(*VoiceServerUpdate))
}

// HandleVoiceServerUpdate registers the handler function for the "VOICE_SERVER_UPDATE" event.
// Fired when a guild's voice server is updated. This is Fired when initially connecting to voice,
// and when the current voice instance fails over to a new server.
func (c *Client) HandleVoiceServerUpdate(f func(vs *VoiceServerUpdate)) {
	c.registerHandler(eventVoiceServerUpdate, voiceServerUpdateHandler(f))
}

type WebhooksUpdate struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
}

type webhooksUpdateHandler func(*WebhooksUpdate)

func (h webhooksUpdateHandler) handle(v interface{}) {
	h(v.(*WebhooksUpdate))
}

// HandleWebhooksUpdate registers the handler function for the "WEBHOOKS_UPDATE" event.
// Fired when a guild channel's webhook is created, updated, or deleted.
func (c *Client) HandleWebhooksUpdate(f func(wu *WebhooksUpdate)) {
	c.registerHandler(eventWebhooksUpdate, webhooksUpdateHandler(f))
}
