package harmony

import (
	"encoding/json"
	"fmt"

	"github.com/skwair/harmony/discord"
	"github.com/skwair/harmony/voice"
)

const (
	eventHello                      = "HELLO"
	eventReady                      = "READY"
	eventResumed                    = "RESUMED"
	eventInvalidSession             = "INVALID_SESSION"
	eventChannelCreate              = "CHANNEL_CREATE"
	eventChannelUpdate              = "CHANNEL_UPDATE"
	eventChannelDelete              = "CHANNEL_DELETE"
	eventChannelPinsUpdate          = "CHANNEL_PINS_UPDATE"
	eventGuildCreate                = "GUILD_CREATE"
	eventGuildUpdate                = "GUILD_UPDATE"
	eventGuildDelete                = "GUILD_DELETE"
	eventGuildBanAdd                = "GUILD_BAN_ADD"
	eventGuildBanRemove             = "GUILD_BAN_REMOVE"
	eventGuildEmojisUpdate          = "GUILD_EMOJIS_UPDATE"
	eventGuildIntegrationsUpdate    = "GUILD_INTEGRATIONS_UPDATE"
	eventGuildMemberAdd             = "GUILD_MEMBER_ADD"
	eventGuildMemberRemove          = "GUILD_MEMBER_REMOVE"
	eventGuildMemberUpdate          = "GUILD_MEMBER_UPDATE"
	eventGuildMembersChunk          = "GUILD_MEMBERS_CHUNK"
	eventGuildRoleCreate            = "GUILD_ROLE_CREATE"
	eventGuildRoleUpdate            = "GUILD_ROLE_UPDATE"
	eventGuildRoleDelete            = "GUILD_ROLE_DELETE"
	eventGuildInviteCreate          = "INVITE_CREATE"
	eventGuildInviteDelete          = "INVITE_DELETE"
	eventMessageCreate              = "MESSAGE_CREATE"
	eventMessageUpdate              = "MESSAGE_UPDATE"
	eventMessageDelete              = "MESSAGE_DELETE"
	eventMessageDeleteBulk          = "MESSAGE_DELETE_BULK"
	eventMessageAck                 = "MESSAGE_ACK"
	eventMessageReactionAdd         = "MESSAGE_REACTION_ADD"
	eventMessageReactionRemove      = "MESSAGE_REACTION_REMOVE"
	eventMessageReactionRemoveAll   = "MESSAGE_REACTION_REMOVE_ALL"
	eventMessageReactionRemoveEmoji = "MESSAGE_REACTION_REMOVE_EMOJI"
	eventPresenceUpdate             = "PRESENCE_UPDATE"
	eventTypingStart                = "TYPING_START"
	eventUserUpdate                 = "USER_UPDATE"
	eventVoiceStateUpdate           = "VOICE_STATE_UPDATE"
	eventVoiceServerUpdate          = "VOICE_SERVER_UPDATE"
	eventWebhooksUpdate             = "WEBHOOKS_UPDATE"
)

// NOTE: consider using a map[string]sync.Pool to cache event objects.

// dispatch dispatches events to user handlers, updating the State
// if it is enabled.
func (c *Client) dispatch(typ string, data json.RawMessage) error {
	var err error
	switch typ {
	case eventHello:
	case eventReady:
		c.connected.Store(true)
		var r Ready
		if err = json.Unmarshal(data, &r); err != nil {
			return fmt.Errorf("unmarshal ready event: %w", err)
		}
		c.handle(eventReady, &r)
	case eventResumed:
		c.connected.Store(true)
	case eventInvalidSession:

	case eventChannelCreate:
		var ch discord.Channel
		if err = json.Unmarshal(data, &ch); err != nil {
			return fmt.Errorf("unmarshal channel create event: %w", err)
		}
		if c.withStateTracking {
			c.State.updateChannel(&ch)
		}
		c.handle(eventChannelCreate, &ch)
	case eventChannelUpdate:
		var ch discord.Channel
		if err = json.Unmarshal(data, &ch); err != nil {
			return fmt.Errorf("unmarshal channel update event: %w", err)
		}
		if c.withStateTracking {
			c.State.updateChannel(&ch)
		}
		c.handle(eventChannelUpdate, &ch)
	case eventChannelDelete:
		var ch discord.Channel
		if err = json.Unmarshal(data, &ch); err != nil {
			return fmt.Errorf("unmarshal channel delete event: %w", err)
		}
		if c.withStateTracking {
			c.State.removeChannel(&ch)
		}
		c.handle(eventChannelDelete, &ch)

	case eventChannelPinsUpdate:
		var pins ChannelPinsUpdate
		if err = json.Unmarshal(data, &pins); err != nil {
			return fmt.Errorf("unmarshal channel pins update event: %w", err)
		}
		if c.withStateTracking {
			c.State.updatePins(&pins)
		}
		c.handle(eventChannelPinsUpdate, &pins)

	case eventGuildCreate:
		var g discord.Guild
		fmt.Println(string(data))
		if err = json.Unmarshal(data, &g); err != nil {
			return fmt.Errorf("unmarshal guild create event: %w", err)
		}
		if c.withStateTracking {
			c.State.updateGuild(&g)
		}
		c.handle(eventGuildCreate, &g)
	case eventGuildUpdate:
		var g discord.Guild
		if err = json.Unmarshal(data, &g); err != nil {
			return fmt.Errorf("unmarshal guild update event: %w", err)
		}
		if c.withStateTracking {
			c.State.updateGuild(&g)
		}
		c.handle(eventGuildUpdate, &g)
	case eventGuildDelete:
		var g discord.UnavailableGuild
		if err = json.Unmarshal(data, &g); err != nil {
			return fmt.Errorf("unmarshal guild delete event: %w", err)
		}
		if c.withStateTracking {
			c.State.removeGuild(&g)
		}
		c.handle(eventGuildDelete, &g)

	case eventGuildBanAdd:
		var ban GuildBan
		if err = json.Unmarshal(data, &ban); err != nil {
			return fmt.Errorf("unmarshal ban add event: %w", err)
		}
		c.handle(eventGuildBanAdd, &ban)
	case eventGuildBanRemove:
		var ban GuildBan
		if err = json.Unmarshal(data, &ban); err != nil {
			return fmt.Errorf("unmarshal ban remove event: %w", err)
		}
		c.handle(eventGuildBanRemove, &ban)

	case eventGuildEmojisUpdate:
		var ge GuildEmojis
		if err = json.Unmarshal(data, &ge); err != nil {
			return fmt.Errorf("unmarshal guild emojis update event: %w", err)
		}
		if c.withStateTracking {
			c.State.updateGuildEmojis(ge.GuildID, ge.Emojis)
		}
		c.handle(eventGuildEmojisUpdate, &ge)

	case eventGuildIntegrationsUpdate:
		var giu GuildIntegrationUpdate
		if err = json.Unmarshal(data, &giu); err != nil {
			return fmt.Errorf("unmarshal guild integrations update event: %w", err)
		}
		c.handle(eventGuildIntegrationsUpdate, &giu)

	case eventGuildMemberAdd:
		var m GuildMemberAdd
		if err = json.Unmarshal(data, &m); err != nil {
			return fmt.Errorf("unmarshal guild member add event: %w", err)
		}
		if c.withStateTracking {
			c.State.guildMemberAdd(&m)
		}
		c.handle(eventGuildMemberAdd, &m)
	case eventGuildMemberRemove:
		var m GuildMemberRemove
		if err = json.Unmarshal(data, &m); err != nil {
			return fmt.Errorf("unmarshal guild member remove event: %w", err)
		}
		if c.withStateTracking {
			c.State.guildMemberRemove(&m)
		}
		c.handle(eventGuildMemberRemove, &m)
	case eventGuildMemberUpdate:
		var m GuildMemberUpdate
		if err = json.Unmarshal(data, &m); err != nil {
			return fmt.Errorf("unmarshal guild member update event: %w", err)
		}
		if c.withStateTracking {
			c.State.guildMemberUpdate(&m)
		}
		c.handle(eventGuildMemberUpdate, &m)

	case eventGuildMembersChunk:
		var chunk GuildMembersChunk
		if err = json.Unmarshal(data, &chunk); err != nil {
			return fmt.Errorf("unmarshal guild members chunk event: %w", err)
		}
		if c.withStateTracking {
			for _, m := range chunk.Members {
				c.State.guildMemberUpdate(&GuildMemberUpdate{
					GuildID: chunk.GuildID,
					Roles:   m.Roles,
					User:    m.User,
					Nick:    m.Nick,
				})
			}
		}
		c.handle(eventGuildMembersChunk, &chunk)

	case eventGuildRoleCreate:
		var gr GuildRole
		if err = json.Unmarshal(data, &gr); err != nil {
			return fmt.Errorf("unmarshal guild role create event: %w", err)
		}
		if c.withStateTracking {
			c.State.guildRoleCreate(&gr)
		}
		c.handle(eventGuildRoleCreate, &gr)
	case eventGuildRoleUpdate:
		var gr GuildRole
		if err = json.Unmarshal(data, &gr); err != nil {
			return fmt.Errorf("unmarshal guild role update event: %w", err)
		}
		if c.withStateTracking {
			c.State.guildRoleUpdate(&gr)
		}
		c.handle(eventGuildRoleUpdate, &gr)
	case eventGuildRoleDelete:
		var gr GuildRoleDelete
		if err = json.Unmarshal(data, &gr); err != nil {
			return fmt.Errorf("unmarshal guild role delete event: %w", err)
		}
		if c.withStateTracking {
			c.State.guildRoleRemove(&gr)
		}
		c.handle(eventGuildRoleDelete, &gr)
	case eventGuildInviteCreate:
		var gic GuildInviteCreate
		if err = json.Unmarshal(data, &gic); err != nil {
			return fmt.Errorf("unmarshal guild invite create event: %w", err)
		}
		c.handle(eventGuildInviteCreate, &gic)
	case eventGuildInviteDelete:
		var gid GuildInviteDelete
		if err = json.Unmarshal(data, &gid); err != nil {
			return fmt.Errorf("unmarshal guild invite delete event: %w", err)
		}
		c.handle(eventGuildInviteDelete, &gid)

	case eventMessageCreate:
		var msg discord.Message
		if err = json.Unmarshal(data, &msg); err != nil {
			return fmt.Errorf("unmarshal message create event: %w", err)
		}
		c.handle(eventMessageCreate, &msg)
	case eventMessageUpdate:
		var msg discord.Message
		if err = json.Unmarshal(data, &msg); err != nil {
			return fmt.Errorf("unmarshal message update event: %w", err)
		}
		c.handle(eventMessageUpdate, &msg)
	case eventMessageDelete:
		var md MessageDelete
		if err = json.Unmarshal(data, &md); err != nil {
			return fmt.Errorf("unmarshal message delete event: %w", err)
		}
		c.handle(eventMessageDelete, &md)
	case eventMessageDeleteBulk:
		var md MessageDeleteBulk
		if err = json.Unmarshal(data, &md); err != nil {
			return fmt.Errorf("unmarshal message delete bulk event: %w", err)
		}
		c.handle(eventMessageDeleteBulk, &md)
	case eventMessageAck:
		var ma MessageAck
		if err = json.Unmarshal(data, &ma); err != nil {
			return fmt.Errorf("unmarshal message ack event: %w", err)
		}
		c.handle(eventMessageAck, &ma)

	case eventMessageReactionAdd:
		var mr MessageReaction
		if err = json.Unmarshal(data, &mr); err != nil {
			return fmt.Errorf("unmarshal message reaction add event: %w", err)
		}
		c.handle(eventMessageReactionAdd, &mr)
	case eventMessageReactionRemove:
		var mr MessageReaction
		if err = json.Unmarshal(data, &mr); err != nil {
			return fmt.Errorf("unmarshal message reaction remove event: %w", err)
		}
		c.handle(eventMessageReactionRemove, &mr)
	case eventMessageReactionRemoveAll:
		var mr MessageReactionRemoveAll
		if err = json.Unmarshal(data, &mr); err != nil {
			return fmt.Errorf("unmarshal message reaction remove all event: %w", err)
		}
		c.handle(eventMessageReactionRemoveAll, &mr)
	case eventMessageReactionRemoveEmoji:
		var m MessageReactionRemoveEmoji
		if err = json.Unmarshal(data, &m); err != nil {
			return fmt.Errorf("unmarshal message reaction remove emoji event: %w", err)
		}
		c.handle(eventMessageReactionRemoveEmoji, &m)

	case eventPresenceUpdate:
		var p discord.Presence
		if err = json.Unmarshal(data, &p); err != nil {
			return fmt.Errorf("unmarshal presence update event: %w", err)
		}
		if c.withStateTracking {
			c.State.updatePresence(&p)
		}
		c.handle(eventPresenceUpdate, &p)

	case eventTypingStart:
		var ts TypingStart
		if err = json.Unmarshal(data, &ts); err != nil {
			return fmt.Errorf("unmarshal typing start event: %w", err)
		}
		c.handle(eventTypingStart, &ts)

	case eventUserUpdate:
		var u discord.User
		if err = json.Unmarshal(data, &u); err != nil {
			return fmt.Errorf("unmarshal user update event: %w", err)
		}
		if c.withStateTracking {
			c.State.updateUser(&u)
		}
		c.handle(eventUserUpdate, &u)

	case eventVoiceStateUpdate:
		var vs voice.StateUpdate
		if err = json.Unmarshal(data, &vs); err != nil {
			return fmt.Errorf("unmarshal voice state update event: %w", err)
		}

		// If this update concerns a voice connection managed
		// by the Client and it's not a channel leave, make
		// sure to update its state so it stays coherent.
		// Failing to do so would make this connection try to
		// reconnect to a wrong channel or with a wrong state
		// (deafen/muted) if it had to reconnect.
		if vs.UserID == c.userID && vs.ChannelID != nil {
			conn, ok := c.voiceConnections[vs.GuildID]
			if ok {
				conn.SetState(&vs.State)
			}
		}

		if c.withStateTracking {
			c.State.updateGuildVoiceStates(&vs)
		}
		c.handle(eventVoiceStateUpdate, &vs)
	case eventVoiceServerUpdate:
		var vs voice.ServerUpdate
		if err = json.Unmarshal(data, &vs); err != nil {
			return fmt.Errorf("unmarshal voice server update event: %w", err)
		}

		// If this update concerns a voice connection managed
		// by the Client, make sure to update it accordingly
		// so it can connect to the new voice server.
		if conn, ok := c.voiceConnections[vs.GuildID]; ok {
			go func() {
				if err = conn.UpdateServer(&vs); err != nil {
					c.logger.Errorf("could not update voice server (guild=%q): %v", vs.GuildID, err)
					return
				}
				c.logger.Debugf("successfully update voice server (guild=%q)", vs.GuildID)
			}()
		}

		c.handle(eventVoiceServerUpdate, &vs)

	case eventWebhooksUpdate:
		var wu WebhooksUpdate
		if err = json.Unmarshal(data, &wu); err != nil {
			return fmt.Errorf("unmarshal webhooks update event: %w", err)
		}
		c.handle(eventWebhooksUpdate, &wu)

	default:
		c.logger.Infof("unrecognized event %s: %s", typ, string(data))
		return nil
	}
	return nil
}

// handle calls the registered user event handler for the given event,
// if there is one.
func (c *Client) handle(event string, d interface{}) {
	c.handlersMu.RLock()
	h, ok := c.handlers[event]
	c.handlersMu.RUnlock()
	if ok {
		// Call the registered handler in its own goroutine
		// so it does not block the dispatcher and events
		// can continue to be treated as we receive them.
		go h.handle(d)
	}
}
