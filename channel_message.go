package harmony

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/skwair/harmony/channel"
	"github.com/skwair/harmony/embed"
	"github.com/skwair/harmony/internal/endpoint"
	"github.com/skwair/harmony/message"
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
	Author *User `json:"author"`
	// Guild member info of the author that sent the message.
	// Only set for MESSAGE_CREATE and MESSAGE_UPDATE Gateway
	// events.
	Member          *GuildMember `json:"member"`
	Content         string       `json:"content"`
	Timestamp       time.Time    `json:"timestamp"`
	EditedTimestamp time.Time    `json:"edited_timestamp"`
	TTS             bool         `json:"tts"`
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
	MentionChannels []channel.Mention    `json:"mention_channels"`
	Attachments     []message.Attachment `json:"attachments"` // Any attached files.
	Embeds          []embed.Embed        `json:"embeds"`      // Any embedded content.
	Reactions       []Reaction           `json:"reactions"`
	Nonce           string               `json:"nonce"` // Used for validating a message was sent.
	Pinned          bool                 `json:"pinned"`
	WebhookID       string               `json:"webhook_id"`
	Type            message.Type         `json:"type"`

	// Sent with Rich Presence-related chat embeds.
	Activity         *MessageActivity    `json:"activity"`
	Application      *MessageApplication `json:"application"`
	MessageReference *message.Reference  `json:"message_reference"`
	Flags            message.Flag        `json:"flags"`
}

// Messages returns messages in the channel. If operating on a guild channel, this
// endpoint requires the 'VIEW_CHANNEL' permission to be present on the current user. If the
// current user is missing the 'READ_MESSAGE_HISTORY' permission in the channel then this will
// return no messages (since they cannot read the message history).
// The query parameter is a message ID prefixed with one of the following character:
//	- '>' for fetching messages after
//	- '<' for fetching messages before
//	- '~' for fetching messages around
// For example, to retrieve 50 messages around (25 before, 25 after) a message having the
// ID 221588207995121520, set query to "~221588207995121520".
// Limit is a positive integer between 1 and 100 that defaults to 50 if set to 0.
func (r *ChannelResource) Messages(ctx context.Context, query string, limit int) ([]Message, error) {
	if query == "" {
		return nil, errors.New("empty query")
	}

	q := url.Values{}
	switch query[0] {
	case '>':
		q.Add("after", query[1:])
	case '<':
		q.Add("before", query[1:])
	case '~':
		q.Add("around", query[1:])
	default:
		return nil, errors.New("lll-formatted query: prefix the message ID with '>' (after), '<' (before) or '~' (around)")
	}

	if limit > 0 {
		if limit > 100 {
			limit = 100
		}
		q.Set("limit", strconv.Itoa(limit))
	}

	e := endpoint.GetChannelMessages(r.channelID, q.Encode())
	resp, err := r.client.doReq(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var msgs []Message
	if err = json.NewDecoder(resp.Body).Decode(&msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}

// Message returns a specific message in the channel. If operating on a guild channel,
// this endpoints requires the 'READ_MESSAGE_HISTORY' permission to be present on the current user.
func (r *ChannelResource) Message(ctx context.Context, id string) (*Message, error) {
	e := endpoint.GetChannelMessage(r.channelID, id)
	resp, err := r.client.doReq(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var msg Message
	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// DeleteMessage is like DeleteMessageWithReason but with no particular reason.
func (r *ChannelResource) DeleteMessage(ctx context.Context, messageID string) error {
	return r.DeleteMessageWithReason(ctx, messageID, "")
}

// DeleteMessageWithReason deletes a message. If operating on a guild channel and trying to delete a
// message that was not sent by the current user, this endpoint requires the 'MANAGE_MESSAGES'
// permission. Fires a Message Delete Gateway event.
// The given reason will be set in the audit log entry for this action.
func (r *ChannelResource) DeleteMessageWithReason(ctx context.Context, messageID, reason string) error {
	e := endpoint.DeleteMessage(r.channelID, messageID)
	resp, err := r.client.doReqWithHeader(ctx, e, nil, reasonHeader(reason))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// DeleteMessageBulk deletes multiple messages in a single request. This endpoint can only be
// used on guild channels and requires the 'MANAGE_MESSAGES' permission. Fires multiple
// Message Delete Gateway events.
// Any message IDs given that do not exist or are invalid will count towards the minimum and
// maximum message count (currently 2 and 100 respectively). Additionally, duplicated IDs will
// only be counted once.
// This endpoint will not delete messages older than 2 weeks, and will fail if any message
// provided is older than that.
func (r *ChannelResource) DeleteMessageBulk(ctx context.Context, messageIDs []string) error {
	st := struct {
		Messages []string `json:"messages"`
	}{
		Messages: messageIDs,
	}
	b, err := json.Marshal(st)
	if err != nil {
		return err
	}

	e := endpoint.BulkDeleteMessage(r.channelID)
	resp, err := r.client.doReq(ctx, e, jsonPayload(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// MessageOption allows to customize the content of a message.
type MessageOption func(*createMessage)

// WithContent sets the content of a message, up to 2000 characters
func WithContent(text string) MessageOption {
	return MessageOption(func(m *createMessage) {
		m.Content = text
	})
}

// WithEmbed sets the embed of a message.
// See embed sub package for more information about embeds.
func WithEmbed(e *embed.Embed) MessageOption {
	return MessageOption(func(m *createMessage) {
		m.Embed = e
	})
}

// WithFiles attach files to a message.
func WithFiles(files ...*File) MessageOption {
	return MessageOption(func(m *createMessage) {
		for _, f := range files {
			m.files = append(m.files, *f)
		}
	})
}

// WithTTS enables text to speech for a message.
func WithTTS() MessageOption {
	return MessageOption(func(m *createMessage) {
		m.TTS = true
	})
}

// WithNonce sets the nonce of a message.
// The nonce will be returned in the result and also transmitted to other clients.
func WithNonce(n string) MessageOption {
	return MessageOption(func(m *createMessage) {
		m.Nonce = n
	})
}

// Send sends a message to the channel. If operating on a guild channel,
// this endpoint requires the 'SEND_MESSAGES' permission to be present on the
// current user. If the option WithTTS is set, the 'SEND_TTS_MESSAGES' permission is
// required for the message to be spoken. Returns the message sent.
// Fires a Message Create Gateway event.
// Before using this endpoint, you must connect to the gateway at least once.
func (r *ChannelResource) Send(ctx context.Context, opts ...MessageOption) (*Message, error) {
	var msg createMessage

	for _, opt := range opts {
		opt(&msg)
	}

	if msg.Content == "" && msg.Embed == nil && len(msg.files) == 0 {
		return nil, ErrInvalidSend
	}

	return r.client.sendMessage(ctx, r.channelID, &msg)
}

// SendMessage is a shorthand for Send(ctx, WithContent(text)).
func (r *ChannelResource) SendMessage(ctx context.Context, text string) (*Message, error) {
	return r.Send(ctx, WithContent(text))
}

// createMessage describes a message creation.
type createMessage struct {
	Content string       `json:"content,omitempty"` // Up to 2000 characters.
	Nonce   string       `json:"nonce,omitempty"`
	TTS     bool         `json:"tts,omitempty"`
	Embed   *embed.Embed `json:"embed,omitempty"`

	files []File
}

// json implements the multipartPayload interface so createMessage can be used as
// a payload with the multipartFromFiles method.
func (cm *createMessage) json() ([]byte, error) {
	return json.Marshal(cm)
}

func (c *Client) sendMessage(ctx context.Context, channelID string, msg *createMessage) (*Message, error) {
	if msg.Embed != nil && msg.Embed.Type == "" {
		msg.Embed.Type = "rich"
	}

	var payload *requestPayload
	if len(msg.files) > 0 {
		b, contentType, err := multipartFromFiles(msg, msg.files...)
		if err != nil {
			return nil, err
		}
		payload = customPayload(b, contentType)
	} else {
		b, err := json.Marshal(msg)
		if err != nil {
			return nil, err
		}
		payload = jsonPayload(b)
	}

	e := endpoint.CreateMessage(channelID)
	resp, err := c.doReq(ctx, e, payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var m Message
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

type editMessage struct {
	Content string       `json:"content,omitempty"`
	Embed   *embed.Embed `json:"embed,omitempty"`
}

// EditMessage edits a previously sent message. You can only edit messages that have
// been sent by the current user. Fires a Message Update Gateway event. See EditEmbed
// if you need to edit some emended content.
func (r *ChannelResource) EditMessage(ctx context.Context, messageID, content string) (*Message, error) {
	return r.client.editMessage(ctx, r.channelID, messageID, &editMessage{content, nil})
}

// EditEmbed is like EditMessage but with embedded content support.
func (r *ChannelResource) EditEmbed(ctx context.Context, messageID, content string, embed *embed.Embed) (*Message, error) {
	return r.client.editMessage(ctx, r.channelID, messageID, &editMessage{content, embed})
}

func (c *Client) editMessage(ctx context.Context, channelID, messageID string, edit *editMessage) (*Message, error) {
	b, err := json.Marshal(edit)
	if err != nil {
		return nil, err
	}

	e := endpoint.EditMessage(channelID, messageID)
	resp, err := c.doReq(ctx, e, jsonPayload(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var msg Message
	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// Reaction is a reaction on a Discord message.
type Reaction struct {
	Count int    `json:"count"`
	Me    bool   `json:"me"`
	Emoji *Emoji `json:"emoji"`
}

// Reactions returns a list of users that reacted to a message with the given emoji.
// limit is the number of users to return and can be set to any value ranging from 1 to 100.
// If set to 0, it defaults to 25. If more than 100 users reacted with the given emoji,
// the before and after parameters can be used to fetch more users.
func (r *ChannelResource) Reactions(ctx context.Context, messageID, emoji string, limit int, before, after string) ([]User, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if before != "" {
		q.Set("before", before)
	}
	if after != "" {
		q.Set("after", after)
	}

	e := endpoint.GetReactions(r.channelID, messageID, url.PathEscape(emoji), q.Encode())
	resp, err := r.client.doReq(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var users []User
	if err = json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}
	return users, nil
}

// AddReaction adds a reaction to a message in the channel. This endpoint requires
// the 'READ_MESSAGE_HISTORY' permission to be present on the current user. Additionally,
// if nobody else has reacted to the message using this emoji, this endpoint requires
// the 'ADD_REACTIONS' permission to be present on the current user.
func (r *ChannelResource) AddReaction(ctx context.Context, messageID, emoji string) error {
	e := endpoint.CreateReaction(r.channelID, messageID, url.PathEscape(emoji))
	resp, err := r.client.doReq(ctx, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// RemoveReaction removes a reaction the current user has made for the message.
func (r *ChannelResource) RemoveReaction(ctx context.Context, messageID, emoji string) error {
	e := endpoint.DeleteOwnReaction(r.channelID, messageID, url.PathEscape(emoji))
	resp, err := r.client.doReq(ctx, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// RemoveUserReaction removes another user's reaction. This endpoint requires the
// 'MANAGE_MESSAGES' permission to be present on the current user.
func (r *ChannelResource) RemoveUserReaction(ctx context.Context, messageID, userID, emoji string) error {
	e := endpoint.DeleteUserReaction(r.channelID, messageID, userID, url.PathEscape(emoji))
	resp, err := r.client.doReq(ctx, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// RemoveAllReactions removes all reactions on a message. This endpoint requires the
// 'MANAGE_MESSAGES' permission to be present on the current user.
func (r *ChannelResource) RemoveAllReactions(ctx context.Context, messageID string) error {
	e := endpoint.DeleteAllReactions(r.channelID, messageID)
	resp, err := r.client.doReq(ctx, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// Pins returns all pinned messages in the channel as an array of messages.
func (r *ChannelResource) Pins(ctx context.Context) ([]Message, error) {
	e := endpoint.GetPinnedMessages(r.channelID)
	resp, err := r.client.doReq(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var messages []Message
	if err = json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, err
	}
	return messages, nil
}

// PinMessage pins a message in the channel. Requires the 'MANAGE_MESSAGES' permission.
func (r *ChannelResource) PinMessage(ctx context.Context, id string) error {
	e := endpoint.AddPinnedChannelMessage(r.channelID, id)
	resp, err := r.client.doReq(ctx, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// UnpinMessage deletes a pinned message in the channel. Requires the
// 'MANAGE_MESSAGES' permission.
func (r *ChannelResource) UnpinMessage(ctx context.Context, id string) error {
	e := endpoint.DeletePinnedChannelMessage(r.channelID, id)
	resp, err := r.client.doReq(ctx, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}
