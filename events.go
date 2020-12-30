package harmony

import (
	"github.com/skwair/harmony/discord"
	"github.com/skwair/harmony/voice"
)

type handler interface {
	handle(interface{})
}

var intents = map[string]discord.GatewayIntent{
	eventGuildCreate:       discord.GatewayIntentGuild,
	eventGuildUpdate:       discord.GatewayIntentGuild,
	eventGuildRoleCreate:   discord.GatewayIntentGuild,
	eventGuildRoleUpdate:   discord.GatewayIntentGuild,
	eventGuildRoleDelete:   discord.GatewayIntentGuild,
	eventChannelCreate:     discord.GatewayIntentGuild,
	eventChannelUpdate:     discord.GatewayIntentGuild,
	eventChannelDelete:     discord.GatewayIntentGuild,
	eventChannelPinsUpdate: discord.GatewayIntentGuild,

	eventGuildMemberAdd:    discord.GatewayIntentGuildMembers,
	eventGuildMemberUpdate: discord.GatewayIntentGuildMembers,
	eventGuildMemberRemove: discord.GatewayIntentGuildMembers,

	eventGuildBanAdd:    discord.GatewayIntentGuildBans,
	eventGuildBanRemove: discord.GatewayIntentGuildBans,

	eventGuildEmojisUpdate: discord.GatewayIntentGuildEmojis,

	eventGuildIntegrationsUpdate: discord.GatewayIntentGuildIntegrations,

	eventWebhooksUpdate: discord.GatewayIntentGuildWebhooks,

	eventGuildInviteCreate: discord.GatewayIntentGuildInvites,
	eventGuildInviteDelete: discord.GatewayIntentGuildInvites,

	eventVoiceStateUpdate: discord.GatewayIntentGuildVoiceStates,

	eventPresenceUpdate: discord.GatewayIntentGuildPresences,

	eventMessageCreate:     discord.GatewayIntentGuildMessages | discord.GatewayIntentDirectMessages,
	eventMessageUpdate:     discord.GatewayIntentGuildMessages | discord.GatewayIntentDirectMessages,
	eventMessageDelete:     discord.GatewayIntentGuildMessages | discord.GatewayIntentDirectMessages,
	eventMessageDeleteBulk: discord.GatewayIntentGuildMessages | discord.GatewayIntentDirectMessages,

	eventMessageReactionAdd:         discord.GatewayIntentGuildMessageReactions | discord.GatewayIntentDirectMessageReactions,
	eventMessageReactionRemove:      discord.GatewayIntentGuildMessageReactions | discord.GatewayIntentDirectMessageReactions,
	eventMessageReactionRemoveAll:   discord.GatewayIntentGuildMessageReactions | discord.GatewayIntentDirectMessageReactions,
	eventMessageReactionRemoveEmoji: discord.GatewayIntentGuildMessageReactions | discord.GatewayIntentDirectMessageReactions,

	eventTypingStart: discord.GatewayIntentGuildMessageTyping | discord.GatewayIntentDirectMessageTyping,
}

func (c *Client) registerHandler(event string, h handler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if missing := intents[event]; missing&c.intents == 0 {
		c.logger.Warnf("registering handler for event %q without required intent %q", event, missing)
	}

	if h == nil {
		panic("harmony: trying to register a nil event handler")
	}
	c.handlersMu.Lock()
	c.handlers[event] = h
	c.handlersMu.Unlock()

	c.logger.Debugf("registered handler for %s events", event)
}

type readyHandler func(*Ready)

// handle implements the handler interface.
func (h readyHandler) handle(v interface{}) {
	h(v.(*Ready))
}

// OnReady registers the handler function for the "READY" event.
func (c *Client) OnReady(f func(r *Ready)) {
	c.registerHandler(eventReady, readyHandler(f))
}

type channelCreateHandler func(*discord.Channel)

// handle implements the handler interface.
func (h channelCreateHandler) handle(v interface{}) {
	h(v.(*discord.Channel))
}

// OnChannelCreate registers the handler function for the "CHANNEL_CREATE" event.
// This event is fired when a new channel is created, relevant to the current user.
func (c *Client) OnChannelCreate(f func(c *discord.Channel)) {
	c.registerHandler(eventChannelCreate, channelCreateHandler(f))
}

type channelUpdateHandler func(*discord.Channel)

// handle implements the handler interface.
func (h channelUpdateHandler) handle(v interface{}) {
	h(v.(*discord.Channel))
}

// OnChannelUpdate registers the handler function for the "CHANNEL_UPDATE" event.
// This event is fired when a channel is updated, relevant to the current user.
func (c *Client) OnChannelUpdate(f func(c *discord.Channel)) {
	c.registerHandler(eventChannelUpdate, channelUpdateHandler(f))
}

type channelDeleteHandler func(*discord.Channel)

// handle implements the handler interface.
func (h channelDeleteHandler) handle(v interface{}) {
	h(v.(*discord.Channel))
}

// OnChannelDelete registers the handler function for the "CHANNEL_DELETE" event.
// This event is fired when a channel is deleted, relevant to the current user.
func (c *Client) OnChannelDelete(f func(c *discord.Channel)) {
	c.registerHandler(eventChannelDelete, channelDeleteHandler(f))
}

// ChannelPinsUpdate is Fired when a message is pinned or unpinned in a text channel.
type ChannelPinsUpdate struct {
	ChannelID        string       `json:"channel_id"`
	LastPinTimestamp discord.Time `json:"last_pin_timestamp"`
}

type channelPinsUpdateHandler func(*ChannelPinsUpdate)

// handle implements the handler interface.
func (h channelPinsUpdateHandler) handle(v interface{}) {
	h(v.(*ChannelPinsUpdate))
}

// OnChannelPinsUpdate registers the handler function for the "CHANNEL_PINS_UPDATE" event.
// This event is fired when a message is pinned or unpinned, but not when a pinned message
// is deleted.
func (c *Client) OnChannelPinsUpdate(f func(cpu *ChannelPinsUpdate)) {
	c.registerHandler(eventChannelPinsUpdate, channelPinsUpdateHandler(f))
}

type guildCreateHandler func(*discord.Guild)

// handle implements the handler interface.
func (h guildCreateHandler) handle(v interface{}) {
	h(v.(*discord.Guild))
}

// OnGuildCreate registers the handler function for the "GUILD_CREATE" event.
// This event can be sent in three different scenarios:
// 	1. When a user is initially connecting, to lazily load and backfill information for all unavailable guilds sent in the Ready event.
// 	2. When a Guild becomes available again to the client.
// 	3. When the current user joins a new Guild.
func (c *Client) OnGuildCreate(f func(g *discord.Guild)) {
	c.registerHandler(eventGuildCreate, guildCreateHandler(f))
}

type guildUpdateHandler func(*discord.Guild)

// handle implements the handler interface.
func (h guildUpdateHandler) handle(v interface{}) {
	h(v.(*discord.Guild))
}

// HandleGuildUpdate registers the handler function for the "GUILD_UPDATE" event.
func (c *Client) OnGuildUpdate(f func(g *discord.Guild)) {
	c.registerHandler(eventGuildUpdate, guildUpdateHandler(f))
}

type guildDeleteHandler func(*discord.UnavailableGuild)

// handle implements the handler interface.
func (h guildDeleteHandler) handle(v interface{}) {
	h(v.(*discord.UnavailableGuild))
}

// OnGuildDelete registers the handler function for the "GUILD_DELETE" event.
// This event is fired when a guild becomes unavailable during a guild outage,
// or when the user leaves or is removed from a guild. If the unavailable field
// is not set, the user was removed from the guild.
func (c *Client) OnGuildDelete(f func(g *discord.UnavailableGuild)) {
	c.registerHandler(eventGuildDelete, guildDeleteHandler(f))
}

type GuildBan struct {
	*discord.User
	GuildID string `json:"guild_id"`
}

type guildBanAddHandler func(*GuildBan)

// handle implements the handler interface.
func (h guildBanAddHandler) handle(v interface{}) {
	h(v.(*GuildBan))
}

// OnGuildBanAdd registers the handler function for the "GUILD_BAN_ADD" event.
func (c *Client) OnGuildBanAdd(f func(ban *GuildBan)) {
	c.registerHandler(eventGuildBanAdd, guildBanAddHandler(f))
}

type guildBanRemoveHandler func(*GuildBan)

// handle implements the handler interface.
func (h guildBanRemoveHandler) handle(v interface{}) {
	h(v.(*GuildBan))
}

// OnGuildBanRemove registers the handler function for the "GUILD_BAN_REMOVE" event.
// This event is fired when a guild is updated.
func (c *Client) OnGuildBanRemove(f func(ban *GuildBan)) {
	c.registerHandler(eventGuildBanRemove, guildBanRemoveHandler(f))
}

type GuildEmojis struct {
	Emojis  []discord.Emoji `json:"emojis"`
	GuildID string          `json:"guild_id"`
}

type guildEmojisUpdateHandler func(*GuildEmojis)

// handle implements the handler interface.
func (h guildEmojisUpdateHandler) handle(v interface{}) {
	h(v.(*GuildEmojis))
}

// OnGuildEmojisUpdate registers the handler function for the "GUILD_EMOJIS_UPDATE" event.
// Fired when a guild's emojis have been updated.
func (c *Client) OnGuildEmojisUpdate(f func(emojis *GuildEmojis)) {
	c.registerHandler(eventGuildEmojisUpdate, guildEmojisUpdateHandler(f))
}

type GuildIntegrationUpdate struct {
	GuildID string `json:"guild_id"`
}

type guildIntegrationUpdateHandler func(*GuildIntegrationUpdate)

// handle implements the handler interface.
func (h guildIntegrationUpdateHandler) handle(v interface{}) {
	h(v.(*GuildIntegrationUpdate))
}

// OnGuildIntegrationsUpdate registers the handler function for the "GUILD_INTEGRATIONS_UPDATE" event.
// Fired when a guild integration is updated.
func (c *Client) OnGuildIntegrationsUpdate(f func(u *GuildIntegrationUpdate)) {
	c.registerHandler(eventGuildIntegrationsUpdate, guildIntegrationUpdateHandler(f))
}

type GuildMemberAdd struct {
	*discord.GuildMember
	GuildID string `json:"guild_id"`
}

type guildMemberAddHandler func(*GuildMemberAdd)

// handle implements the handler interface.
func (h guildMemberAddHandler) handle(v interface{}) {
	h(v.(*GuildMemberAdd))
}

// OnGuildMemberAdd registers the handler function for the "GUILD_MEMBER_ADD" event.
// Fired when a new user joins a guild.
func (c *Client) OnGuildMemberAdd(f func(m *GuildMemberAdd)) {
	c.registerHandler(eventGuildMemberAdd, guildMemberAddHandler(f))
}

type GuildMemberRemove struct {
	User    *discord.User `json:"user"`
	GuildID string        `json:"guild_id"`
}

type guildMemberRemoveHandler func(*GuildMemberRemove)

// handle implements the handler interface.
func (h guildMemberRemoveHandler) handle(v interface{}) {
	h(v.(*GuildMemberRemove))
}

// OnGuildMemberRemove registers the handler function for the "GUILD_MEMBER_REMOVE" event.
// Fired when a user is removed from a guild (leave/kick/ban).
func (c *Client) OnGuildMemberRemove(f func(m *GuildMemberRemove)) {
	c.registerHandler(eventGuildMemberRemove, guildMemberRemoveHandler(f))
}

type GuildMemberUpdate struct {
	GuildID string        `json:"guild_id"`
	Roles   []string      `json:"roles"`
	User    *discord.User `json:"user"`
	Nick    string        `json:"nick"`
}

type guildMemberUpdateHandler func(*GuildMemberUpdate)

// handle implements the handler interface.
func (h guildMemberUpdateHandler) handle(v interface{}) {
	h(v.(*GuildMemberUpdate))
}

// OnGuildMemberUpdate registers the handler function for the "GUILD_MEMBER_UPDATE" event.
// Fired when a guild member is updated.
func (c *Client) OnGuildMemberUpdate(f func(m *GuildMemberUpdate)) {
	c.registerHandler(eventGuildMemberUpdate, guildMemberUpdateHandler(f))
}

type GuildMembersChunk struct {
	GuildID string                `json:"guild_id"`
	Members []discord.GuildMember `json:"members"`
}

type guildMembersChunkHandler func(*GuildMembersChunk)

// handle implements the handler interface.
func (h guildMembersChunkHandler) handle(v interface{}) {
	h(v.(*GuildMembersChunk))
}

// OnGuildMembersChunk registers the handler function for the "GUILD_MEMBERS_CHUNK" event.
// Sent in response to Guild Request Members.
func (c *Client) OnGuildMembersChunk(f func(m *GuildMembersChunk)) {
	c.registerHandler(eventGuildMembersChunk, guildMembersChunkHandler(f))
}

type GuildRole struct {
	GuildID string        `json:"guild_id"`
	Role    *discord.Role `json:"role"`
}

type guildRoleCreateHandler func(*GuildRole)

// handle implements the handler interface.
func (h guildRoleCreateHandler) handle(v interface{}) {
	h(v.(*GuildRole))
}

// OnGuildRoleCreate registers the handler function for the "GUILD_ROLE_CREATE" event.
// Fired when a guild role is created.
func (c *Client) OnGuildRoleCreate(f func(r *GuildRole)) {
	c.registerHandler(eventGuildRoleCreate, guildRoleCreateHandler(f))
}

type guildRoleUpdateHandler func(*GuildRole)

// handle implements the handler interface.
func (h guildRoleUpdateHandler) handle(v interface{}) {
	h(v.(*GuildRole))
}

// OnGuildRoleUpdate registers the handler function for the "GUILD_ROLE_UPDATE" event.
// Fired when a guild role is updated.
func (c *Client) OnGuildRoleUpdate(f func(r *GuildRole)) {
	c.registerHandler(eventGuildRoleUpdate, guildRoleUpdateHandler(f))
}

type GuildRoleDelete struct {
	GuildID string `json:"guild_id"`
	RoleID  string `json:"role_id"`
}

type guildRoleDeleteHandler func(*GuildRoleDelete)

// handle implements the handler interface.
func (h guildRoleDeleteHandler) handle(v interface{}) {
	h(v.(*GuildRoleDelete))
}

// OnGuildRoleDelete registers the handler function for the "GUILD_ROLE_DELETE" event.
// Fired when a guild role is deleted.
func (c *Client) OnGuildRoleDelete(f func(r *GuildRoleDelete)) {
	c.registerHandler(eventGuildRoleDelete, guildRoleDeleteHandler(f))
}

type GuildInviteCreate struct {
	ChannelID      string        `json:"channel_id"`
	Code           string        `json:"code"`
	CreatedAt      discord.Time  `json:"created_at"`
	GuildID        string        `json:"guild_id"`
	Inviter        *discord.User `json:"inviter"`
	MaxAge         int           `json:"max_age"`
	MaxUses        int           `json:"max_uses"`
	TargetUser     *discord.User `json:"target_user"`
	TargetUserType int           `json:"target_user_type"`
	Temporary      bool          `json:"temporary"`
	Uses           int           `json:"uses"`
}

type guildInviteCreateHandler func(*GuildInviteCreate)

// handle implements the handler interface.
func (h guildInviteCreateHandler) handle(v interface{}) {
	h(v.(*GuildInviteCreate))
}

// OnGuildInviteCreate registers the handler function for the "GUILD_ROLE_DELETE" event.
// Fired when a guild role is deleted.
func (c *Client) OnGuildInviteCreate(f func(i *GuildInviteCreate)) {
	c.registerHandler(eventGuildInviteCreate, guildInviteCreateHandler(f))
}

type GuildInviteDelete struct {
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
	Code      string `json:"code"`
}

type guildInviteDeleteHandler func(*GuildInviteDelete)

// handle implements the handler interface.
func (h guildInviteDeleteHandler) handle(v interface{}) {
	h(v.(*GuildInviteDelete))
}

// OnGuildInviteDelete registers the handler function for the "GUILD_ROLE_DELETE" event.
// Fired when a guild role is deleted.
func (c *Client) OnGuildInviteDelete(f func(i *GuildInviteDelete)) {
	c.registerHandler(eventGuildInviteDelete, guildInviteDeleteHandler(f))
}

type messageCreateHandler func(*discord.Message)

// handle implements the handler interface.
func (h messageCreateHandler) handle(v interface{}) {
	h(v.(*discord.Message))
}

// OnMessageCreate registers the handler function for the "MESSAGE_CREATE" event.
// Fired when a message is created.
func (c *Client) OnMessageCreate(f func(m *discord.Message)) {
	c.registerHandler(eventMessageCreate, messageCreateHandler(f))
}

type messageUpdateHandler func(*discord.Message)

// handle implements the handler interface.
func (h messageUpdateHandler) handle(v interface{}) {
	h(v.(*discord.Message))
}

// OnMessageUpdate registers the handler function for the "MESSAGE_UPDATE" event.
// Fired when a message is updated. Unlike creates, message updates may contain only
// a subset of the full message object payload (but will always contain an id and channel_id).
func (c *Client) OnMessageUpdate(f func(m *discord.Message)) {
	c.registerHandler(eventMessageUpdate, messageUpdateHandler(f))
}

type MessageDelete struct {
	ChannelID string `json:"channel_id"`
	MessageID string `json:"id"`
}

type messageDeleteHandler func(*MessageDelete)

// handle implements the handler interface.
func (h messageDeleteHandler) handle(v interface{}) {
	h(v.(*MessageDelete))
}

// OnMessageDelete registers the handler function for the "MESSAGE_DELETE" event.
// Fired when a message is deleted.
func (c *Client) OnMessageDelete(f func(m *MessageDelete)) {
	c.registerHandler(eventMessageDelete, messageDeleteHandler(f))
}

type MessageDeleteBulk struct {
	GuildID   string   `json:"guild_id"`
	ChannelID string   `json:"channel_id"`
	IDs       []string `json:"ids"`
}

type messageDeleteBulkHandler func(*MessageDeleteBulk)

// handle implements the handler interface.
func (h messageDeleteBulkHandler) handle(v interface{}) {
	h(v.(*MessageDeleteBulk))
}

// OnMessageDeleteBulk registers the handler function for the "MESSAGE_DELETE_BULK" event.
// Fired when multiple messages are deleted at once.
func (c *Client) OnMessageDeleteBulk(f func(mdb *MessageDeleteBulk)) {
	c.registerHandler(eventMessageDeleteBulk, messageDeleteBulkHandler(f))
}

type MessageAck struct {
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
}

type messageAckHandler func(*MessageAck)

// handle implements the handler interface.
func (h messageAckHandler) handle(v interface{}) {
	h(v.(*MessageAck))
}

// OnMessageAck registers the handler function for the "MESSAGE_ACK" event.
func (c *Client) OnMessageAck(f func(ack *MessageAck)) {
	c.registerHandler(eventMessageAck, messageAckHandler(f))
}

type MessageReaction struct {
	UserID    string         `json:"user_id"`
	GuildID   string         `json:"guild_id"`
	ChannelID string         `json:"channel_id"`
	MessageID string         `json:"message_id"`
	Emoji     *discord.Emoji `json:"emoji"`
}

type messageReactionAddHandler func(*MessageReaction)

// handle implements the handler interface.
func (h messageReactionAddHandler) handle(v interface{}) {
	h(v.(*MessageReaction))
}

// OnMessageReactionAdd registers the handler function for the "MESSAGE_REACTION_ADD" event.
// Fired when a user adds a reaction to a message.
func (c *Client) OnMessageReactionAdd(f func(r *MessageReaction)) {
	c.registerHandler(eventMessageReactionAdd, messageReactionAddHandler(f))
}

type messageReactionRemoveHandler func(*MessageReaction)

// handle implements the handler interface.
func (h messageReactionRemoveHandler) handle(v interface{}) {
	h(v.(*MessageReaction))
}

// OnMessageReactionRemove registers the handler function for the "MESSAGE_REACTION_REMOVE" event.
// Fired when a user removes a reaction from a message.
func (c *Client) OnMessageReactionRemove(f func(r *MessageReaction)) {
	c.registerHandler(eventMessageReactionRemove, messageReactionRemoveHandler(f))
}

type MessageReactionRemoveAll struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
}

type messageReactionRemoveAllHandler func(*MessageReactionRemoveAll)

// handle implements the handler interface.
func (h messageReactionRemoveAllHandler) handle(v interface{}) {
	h(v.(*MessageReactionRemoveAll))
}

// OnMessageReactionRemoveAll registers the handler function for the "MESSAGE_REACTION_REMOVE_ALL" event.
// Fired when a user explicitly removes all reactions from a message.
func (c *Client) OnMessageReactionRemoveAll(f func(r *MessageReactionRemoveAll)) {
	c.registerHandler(eventMessageReactionRemoveAll, messageReactionRemoveAllHandler(f))
}

type MessageReactionRemoveEmoji struct {
	GuildID   string         `json:"guild_id"`
	ChannelID string         `json:"channel_id"`
	MessageID string         `json:"message_id"`
	Emoji     *discord.Emoji `json:"emoji"`
}

type messageReactionRemoveEmojiHandler func(*MessageReactionRemoveEmoji)

// handle implements the handler interface.
func (h messageReactionRemoveEmojiHandler) handle(v interface{}) {
	h(v.(*MessageReactionRemoveEmoji))
}

// HandleMessageReactionRemoveEmoji registers the handler function for the "MESSAGE_REACTION_REMOVE_ALL" event.
// Fired when a user explicitly removes all reactions from a message.
func (c *Client) OnMessageReactionRemoveEmoji(f func(r *MessageReactionRemoveEmoji)) {
	c.registerHandler(eventMessageReactionRemoveEmoji, messageReactionRemoveEmojiHandler(f))
}

type presenceUpdateHandler func(*discord.Presence)

// handle implements the handler interface.
func (h presenceUpdateHandler) handle(v interface{}) {
	h(v.(*discord.Presence))
}

// OnPresenceUpdate registers the handler function for the "PRESENCE_UPDATE" event.
// This event is fired when a user's presence is updated for a guild.
// The user object within this event can be partial, the only field which must be sent
// is the id field, everything else is optional. Along with this limitation, no fields
// are required, and the types of the fields are not validated. Your client should expect
// any combination of fields and types within this event.
func (c *Client) OnPresenceUpdate(f func(p *discord.Presence)) {
	c.registerHandler(eventPresenceUpdate, presenceUpdateHandler(f))
}

type TypingStart struct {
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
	UserID    string `json:"user_id"`
	Timestamp int64  `json:"timestamp"`
}

type typingStartHandler func(*TypingStart)

// handle implements the handler interface.
func (h typingStartHandler) handle(v interface{}) {
	h(v.(*TypingStart))
}

// OnTypingStart registers the handler function for the "TYPING_START" event.
// Fired when a user starts typing in a channel.
func (c *Client) OnTypingStart(f func(ts *TypingStart)) {
	c.registerHandler(eventTypingStart, typingStartHandler(f))
}

type userUpdateHandler func(*discord.User)

// handle implements the handler interface.
func (h userUpdateHandler) handle(v interface{}) {
	h(v.(*discord.User))
}

// OnUserUpdate registers the handler function for the "USER_UPDATE" event.
// Fired when properties about the user change.
func (c *Client) OnUserUpdate(f func(u *discord.User)) {
	c.registerHandler(eventUserUpdate, userUpdateHandler(f))
}

type voiceStateUpdateHandler func(*voice.StateUpdate)

// handle implements the handler interface.
func (h voiceStateUpdateHandler) handle(v interface{}) {
	h(v.(*voice.StateUpdate))
}

// OnVoiceStateUpdate registers the handler function for the "VOICE_STATE_UPDATE" event.
// Fired when someone joins/leaves/moves voice channels.
func (c *Client) OnVoiceStateUpdate(f func(update *voice.StateUpdate)) {
	c.registerHandler(eventVoiceStateUpdate, voiceStateUpdateHandler(f))
}

type voiceServerUpdateHandler func(*voice.ServerUpdate)

// handle implements the handler interface.
func (h voiceServerUpdateHandler) handle(v interface{}) {
	h(v.(*voice.ServerUpdate))
}

// OnVoiceServerUpdate registers the handler function for the "VOICE_SERVER_UPDATE" event.
// Fired when a guild's voice server is updated. This is Fired when initially connecting to voice,
// and when the current voice instance fails over to a new server.
func (c *Client) OnVoiceServerUpdate(f func(update *voice.ServerUpdate)) {
	c.registerHandler(eventVoiceServerUpdate, voiceServerUpdateHandler(f))
}

type WebhooksUpdate struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
}

type webhooksUpdateHandler func(*WebhooksUpdate)

// handle implements the handler interface.
func (h webhooksUpdateHandler) handle(v interface{}) {
	h(v.(*WebhooksUpdate))
}

// OnWebhooksUpdate registers the handler function for the "WEBHOOKS_UPDATE" event.
// Fired when a guild channel's webhook is created, updated, or deleted.
func (c *Client) OnWebhooksUpdate(f func(wu *WebhooksUpdate)) {
	c.registerHandler(eventWebhooksUpdate, webhooksUpdateHandler(f))
}
