package harmony

import (
	"encoding/json"
	"fmt"

	"github.com/skwair/harmony/voice"
)

const (
	eventHello                    = "HELLO"
	eventReady                    = "READY"
	eventResumed                  = "RESUMED"
	eventInvalidSession           = "INVALID_SESSION"
	eventChannelCreate            = "CHANNEL_CREATE"
	eventChannelUpdate            = "CHANNEL_UPDATE"
	eventChannelDelete            = "CHANNEL_DELETE"
	eventChannelPinsUpdate        = "CHANNEL_PINS_UPDATE"
	eventGuildCreate              = "GUILD_CREATE"
	eventGuildUpdate              = "GUILD_UPDATE"
	eventGuildDelete              = "GUILD_DELETE"
	eventGuildBanAdd              = "GUILD_BAN_ADD"
	eventGuildBanRemove           = "GUILD_BAN_REMOVE"
	eventGuildEmojisUpdate        = "GUILD_EMOJIS_UPDATE"
	eventGuildIntegrationsUpdate  = "GUILD_INTEGRATIONS_UPDATE"
	eventGuildMemberAdd           = "GUILD_MEMBER_ADD"
	eventGuildMemberRemove        = "GUILD_MEMBER_REMOVE"
	eventGuildMemberUpdate        = "GUILD_MEMBER_UPDATE"
	eventGuildMembersChunk        = "GUILD_MEMBERS_CHUNK"
	eventGuildRoleCreate          = "GUILD_ROLE_CREATE"
	eventGuildRoleUpdate          = "GUILD_ROLE_UPDATE"
	eventGuildRoleDelete          = "GUILD_ROLE_DELETE"
	eventMessageCreate            = "MESSAGE_CREATE"
	eventMessageUpdate            = "MESSAGE_UPDATE"
	eventMessageDelete            = "MESSAGE_DELETE"
	eventMessageDeleteBulk        = "MESSAGE_DELETE_BULK"
	eventMessageAck               = "MESSAGE_ACK"
	eventMessageReactionAdd       = "MESSAGE_REACTION_ADD"
	eventMessageReactionRemove    = "MESSAGE_REACTION_REMOVE"
	eventMessageReactionRemoveAll = "MESSAGE_REACTION_REMOVE_ALL"
	eventPresenceUpdate           = "PRESENCE_UPDATE"
	eventTypingStart              = "TYPING_START"
	eventUserUpdate               = "USER_UPDATE"
	eventVoiceStateUpdate         = "VOICE_STATE_UPDATE"
	eventVoiceServerUpdate        = "VOICE_SERVER_UPDATE"
	eventWebhooksUpdate           = "WEBHOOKS_UPDATE"
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
			return err
		}
		c.handle(eventReady, &r)
	case eventResumed:
		c.connected.Store(true)
	case eventInvalidSession:

	case eventChannelCreate:
		var ch Channel
		if err = json.Unmarshal(data, &ch); err != nil {
			return err
		}
		if c.withStateTracking {
			c.State.updateChannel(&ch)
		}
		c.handle(eventChannelCreate, &ch)
	case eventChannelUpdate:
		var ch Channel
		if err = json.Unmarshal(data, &ch); err != nil {
			return err
		}
		if c.withStateTracking {
			c.State.updateChannel(&ch)
		}
		c.handle(eventChannelUpdate, &ch)
	case eventChannelDelete:
		var ch Channel
		if err = json.Unmarshal(data, &ch); err != nil {
			return err
		}
		if c.withStateTracking {
			c.State.removeChannel(&ch)
		}
		c.handle(eventChannelDelete, &ch)

	case eventChannelPinsUpdate:
		var pins ChannelPinsUpdate
		if err = json.Unmarshal(data, &pins); err != nil {
			return err
		}
		if c.withStateTracking {
			c.State.updatePins(&pins)
		}
		c.handle(eventChannelPinsUpdate, &pins)

	case eventGuildCreate:
		var g Guild
		if err = json.Unmarshal(data, &g); err != nil {
			return err
		}
		if c.withStateTracking {
			c.State.updateGuild(&g)
		}
		c.handle(eventGuildCreate, &g)
	case eventGuildUpdate:
		var g Guild
		if err = json.Unmarshal(data, &g); err != nil {
			return err
		}
		if c.withStateTracking {
			c.State.updateGuild(&g)
		}
		c.handle(eventGuildUpdate, &g)
	case eventGuildDelete:
		var g UnavailableGuild
		if err = json.Unmarshal(data, &g); err != nil {
			return err
		}
		if c.withStateTracking {
			c.State.removeGuild(&g)
		}
		c.handle(eventGuildDelete, &g)

	case eventGuildBanAdd:
		var ban GuildBan
		if err = json.Unmarshal(data, &ban); err != nil {
			return err
		}
		c.handle(eventGuildBanAdd, &ban)
	case eventGuildBanRemove:
		var ban GuildBan
		if err = json.Unmarshal(data, &ban); err != nil {
			return err
		}
		c.handle(eventGuildBanRemove, &ban)

	case eventGuildEmojisUpdate:
		var ge GuildEmojis
		if err = json.Unmarshal(data, &ge); err != nil {
			return err
		}
		if c.withStateTracking {
			c.State.updateGuildEmojis(ge.GuildID, ge.Emojis)
		}
		c.handle(eventGuildEmojisUpdate, &ge)

	case eventGuildIntegrationsUpdate:
		var guildID string
		if err = json.Unmarshal(data, &guildID); err != nil {
			return err
		}
		c.handle(eventGuildIntegrationsUpdate, &guildID)

	case eventGuildMemberAdd:
		var m GuildMemberAdd
		if err = json.Unmarshal(data, &m); err != nil {
			return err
		}
		if c.withStateTracking {
			c.State.guildMemberAdd(&m)
		}
		c.handle(eventGuildMemberAdd, &m)
	case eventGuildMemberRemove:
		var m GuildMemberRemove
		if err = json.Unmarshal(data, &m); err != nil {
			return err
		}
		if c.withStateTracking {
			c.State.guildMemberRemove(&m)
		}
		c.handle(eventGuildMemberRemove, &m)
	case eventGuildMemberUpdate:
		var m GuildMemberUpdate
		if err = json.Unmarshal(data, &m); err != nil {
			return err
		}
		if c.withStateTracking {
			c.State.guildMemberUpdate(&m)
		}
		c.handle(eventGuildMemberUpdate, &m)

	case eventGuildMembersChunk:
		var chunk GuildMembersChunk
		if err = json.Unmarshal(data, &chunk); err != nil {
			return err
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
			return err
		}
		if c.withStateTracking {
			c.State.guildRoleCreate(&gr)
		}
		c.handle(eventGuildRoleCreate, &gr)
	case eventGuildRoleUpdate:
		var gr GuildRole
		if err = json.Unmarshal(data, &gr); err != nil {
			return err
		}
		if c.withStateTracking {
			c.State.guildRoleUpdate(&gr)
		}
		c.handle(eventGuildRoleUpdate, &gr)
	case eventGuildRoleDelete:
		var gr GuildRoleDelete
		if err = json.Unmarshal(data, &gr); err != nil {
			return err
		}
		if c.withStateTracking {
			c.State.guildRoleRemove(&gr)
		}
		c.handle(eventGuildRoleDelete, &gr)

	case eventMessageCreate:
		var msg Message
		if err = json.Unmarshal(data, &msg); err != nil {
			return err
		}
		c.handle(eventMessageCreate, &msg)
	case eventMessageUpdate:
		var msg Message
		if err = json.Unmarshal(data, &msg); err != nil {
			return err
		}
		c.handle(eventMessageUpdate, &msg)
	case eventMessageDelete:
		var md MessageDelete
		if err = json.Unmarshal(data, &md); err != nil {
			return err
		}
		c.handle(eventMessageDelete, &md)
	case eventMessageDeleteBulk:
		var md MessageDeleteBulk
		if err = json.Unmarshal(data, &md); err != nil {
			return err
		}
		c.handle(eventMessageDeleteBulk, &md)
	case eventMessageAck:
		var ma MessageAck
		if err = json.Unmarshal(data, &ma); err != nil {
			return err
		}
		c.handle(eventMessageAck, &ma)

	case eventMessageReactionAdd:
		var mr MessageReaction
		if err = json.Unmarshal(data, &mr); err != nil {
			return err
		}
		c.handle(eventMessageReactionAdd, &mr)
	case eventMessageReactionRemove:
		var mr MessageReaction
		if err = json.Unmarshal(data, &mr); err != nil {
			return err
		}
		c.handle(eventMessageReactionRemove, &mr)
	case eventMessageReactionRemoveAll:
		var mr MessageReactionRemoveAll
		if err = json.Unmarshal(data, &mr); err != nil {
			return err
		}
		c.handle(eventMessageReactionRemoveAll, &mr)

	case eventPresenceUpdate:
		var p Presence
		if err = json.Unmarshal(data, &p); err != nil {
			return err
		}
		if c.withStateTracking {
			c.State.updatePresence(&p)
		}
		c.handle(eventPresenceUpdate, &p)

	case eventTypingStart:
		var ts TypingStart
		if err = json.Unmarshal(data, &ts); err != nil {
			return err
		}
		c.handle(eventTypingStart, &ts)

	case eventUserUpdate:
		var u User
		if err = json.Unmarshal(data, &u); err != nil {
			return err
		}
		if c.withStateTracking {
			c.State.updateUser(&u)
		}
		c.handle(eventUserUpdate, &u)

	case eventVoiceStateUpdate:
		var vs voice.StateUpdate
		if err = json.Unmarshal(data, &vs); err != nil {
			return err
		}

		// If this update concerns a voice connections managed
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
			return err
		}

		if conn, ok := c.voiceConnections[vs.GuildID]; ok {
			go func() {
				if err := conn.UpdateServer(&vs); err != nil {
					fmt.Println("ERROR:", err)
				}
				fmt.Println("--------------> successfully updated voice server! <--------------")
			}()
		}

		c.handle(eventVoiceServerUpdate, &vs)

	case eventWebhooksUpdate:
		var wu WebhooksUpdate
		if err = json.Unmarshal(data, &wu); err != nil {
			return err
		}
		c.handle(eventWebhooksUpdate, &wu)

	default:
		c.logger.Infof("unrecognized event %s: %s", typ, string(data))
		return nil
	}
	return err
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
