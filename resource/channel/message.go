package channel

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/skwair/harmony/discord"
	"github.com/skwair/harmony/internal/endpoint"
	"github.com/skwair/harmony/internal/rest"
)

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
func (r *Resource) Messages(ctx context.Context, query string, limit int) ([]discord.Message, error) {
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
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var msgs []discord.Message
	if err = json.NewDecoder(resp.Body).Decode(&msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}

// Message returns a specific message in the channel. If operating on a guild channel,
// this endpoints requires the 'READ_MESSAGE_HISTORY' permission to be present on the current user.
func (r *Resource) Message(ctx context.Context, id string) (*discord.Message, error) {
	e := endpoint.GetChannelMessage(r.channelID, id)
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var msg discord.Message
	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// DeleteMessage is like DeleteMessageWithReason but with no particular reason.
func (r *Resource) DeleteMessage(ctx context.Context, messageID string) error {
	return r.DeleteMessageWithReason(ctx, messageID, "")
}

// DeleteMessageWithReason deletes a message. If operating on a guild channel and trying to delete a
// message that was not sent by the current user, this endpoint requires the 'MANAGE_MESSAGES'
// permission. Fires a Message Delete Gateway event.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) DeleteMessageWithReason(ctx context.Context, messageID, reason string) error {
	e := endpoint.DeleteMessage(r.channelID, messageID)
	resp, err := r.client.DoWithHeader(ctx, e, nil, rest.ReasonHeader(reason))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return discord.NewAPIError(resp)
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
func (r *Resource) DeleteMessageBulk(ctx context.Context, messageIDs []string) error {
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
	resp, err := r.client.Do(ctx, e, rest.JSONPayload(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return discord.NewAPIError(resp)
	}
	return nil
}

// MessageOption allows to customize the content of a message.
type MessageOption func(*createMessage)

// WithMessageContent sets the content of a message, up to 2000 characters
func WithMessageContent(text string) MessageOption {
	return func(m *createMessage) {
		m.Content = text
	}
}

// WithMessageEmbed sets the embed of a message.
// See embed sub package for more information about embeds.
func WithMessageEmbed(e *discord.MessageEmbed) MessageOption {
	return func(m *createMessage) {
		m.Embed = e
	}
}

// WithMessageFiles attach files to a message.
func WithMessageFiles(files ...*discord.File) MessageOption {
	return func(m *createMessage) {
		for _, f := range files {
			m.files = append(m.files, *f)
		}
	}
}

// WithMessageTTS enables text to speech for a message.
func WithMessageTTS() MessageOption {
	return func(m *createMessage) {
		m.TTS = true
	}
}

// WithMessageNonce sets the nonce of a message.
// The nonce will be returned in the result and also transmitted to other clients.
func WithMessageNonce(n string) MessageOption {
	return func(m *createMessage) {
		m.Nonce = n
	}
}

// Send sends a message to the channel. If operating on a guild channel,
// this endpoint requires the 'SEND_MESSAGES' permission to be present on the
// current user. If the option WithMessageTTS is set, the 'SEND_TTS_MESSAGES' permission is
// required for the message to be spoken. Returns the message sent.
// Fires a Message Create Gateway event.
// Before using this endpoint, you must connect to the gateway at least once.
func (r *Resource) Send(ctx context.Context, opts ...MessageOption) (*discord.Message, error) {
	var msg createMessage

	for _, opt := range opts {
		opt(&msg)
	}

	if msg.Content == "" && msg.Embed == nil && len(msg.files) == 0 {
		return nil, discord.ErrInvalidMessageSend
	}

	return r.sendMessage(ctx, r.channelID, &msg)
}

// SendMessage is a shorthand for Send(ctx, WithMessageContent(text)).
func (r *Resource) SendMessage(ctx context.Context, text string) (*discord.Message, error) {
	return r.Send(ctx, WithMessageContent(text))
}

// createMessage describes a message creation.
type createMessage struct {
	Content string                `json:"content,omitempty"` // Up to 2000 characters.
	Nonce   string                `json:"nonce,omitempty"`
	TTS     bool                  `json:"tts,omitempty"`
	Embed   *discord.MessageEmbed `json:"embed,omitempty"`

	files []discord.File
}

// Bytes implements the rest.MultipartPayload interface so createMessage can be used as
// a payload with the rest.MultipartFromFiles function.
func (cm *createMessage) Bytes() ([]byte, error) {
	return json.Marshal(cm)
}

func (r *Resource) sendMessage(ctx context.Context, channelID string, msg *createMessage) (*discord.Message, error) {
	if msg.Embed != nil && msg.Embed.Type == "" {
		msg.Embed.Type = "rich"
	}

	var payload *rest.Payload
	if len(msg.files) > 0 {
		b, contentType, err := rest.MultipartFromFiles(msg, msg.files...)
		if err != nil {
			return nil, err
		}
		payload = rest.CustomPayload(b, contentType)
	} else {
		b, err := json.Marshal(msg)
		if err != nil {
			return nil, err
		}
		payload = rest.JSONPayload(b)
	}

	e := endpoint.CreateMessage(channelID)
	resp, err := r.client.Do(ctx, e, payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var m discord.Message
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

type editMessage struct {
	Content string                `json:"content,omitempty"`
	Embed   *discord.MessageEmbed `json:"embed,omitempty"`
}

// EditMessage edits a previously sent message. You can only edit messages that have
// been sent by the current user. Fires a Message Update Gateway event. See EditEmbed
// if you need to edit some emended content.
func (r *Resource) EditMessage(ctx context.Context, messageID, content string) (*discord.Message, error) {
	return r.editMessage(ctx, r.channelID, messageID, &editMessage{content, nil})
}

// EditEmbed is like EditMessage but with embedded content support.
func (r *Resource) EditEmbed(ctx context.Context, messageID, content string, embed *discord.MessageEmbed) (*discord.Message, error) {
	return r.editMessage(ctx, r.channelID, messageID, &editMessage{content, embed})
}

func (r *Resource) editMessage(ctx context.Context, channelID, messageID string, edit *editMessage) (*discord.Message, error) {
	b, err := json.Marshal(edit)
	if err != nil {
		return nil, err
	}

	e := endpoint.EditMessage(channelID, messageID)
	resp, err := r.client.Do(ctx, e, rest.JSONPayload(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var msg discord.Message
	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// Crosspost a message in a News Channel to following channels.
func (r *Resource) CrossPostMessage(ctx context.Context, messageID string) (*discord.Message, error) {
	e := endpoint.CrossPostMessage(r.channelID, messageID)

	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var msg discord.Message
	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// Reactions returns a list of users that reacted to a message with the given emoji.
// limit is the number of users to return and can be set to any value ranging from 1 to 100.
// If set to 0, it defaults to 25. If more than 100 users reacted with the given emoji,
// the before and after parameters can be used to fetch more users.
func (r *Resource) Reactions(ctx context.Context, messageID, emoji string, limit int, before, after string) ([]discord.User, error) {
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
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var users []discord.User
	if err = json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}
	return users, nil
}

// AddReaction adds a reaction to a message in the channel. This endpoint requires
// the 'READ_MESSAGE_HISTORY' permission to be present on the current user. Additionally,
// if nobody else has reacted to the message using this emoji, this endpoint requires
// the 'ADD_REACTIONS' permission to be present on the current user.
func (r *Resource) AddReaction(ctx context.Context, messageID, emoji string) error {
	e := endpoint.CreateReaction(r.channelID, messageID, url.PathEscape(emoji))
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return discord.NewAPIError(resp)
	}
	return nil
}

// RemoveReaction removes a reaction the current user has made for the message.
func (r *Resource) RemoveReaction(ctx context.Context, messageID, emoji string) error {
	e := endpoint.DeleteOwnReaction(r.channelID, messageID, url.PathEscape(emoji))
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return discord.NewAPIError(resp)
	}
	return nil
}

// RemoveUserReaction removes another user's reaction. This endpoint requires the
// 'MANAGE_MESSAGES' permission to be present on the current user.
func (r *Resource) RemoveUserReaction(ctx context.Context, messageID, userID, emoji string) error {
	e := endpoint.DeleteUserReaction(r.channelID, messageID, userID, url.PathEscape(emoji))
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return discord.NewAPIError(resp)
	}
	return nil
}

// RemoveAllReactions removes all reactions on a message. This endpoint requires the
// 'MANAGE_MESSAGES' permission to be present on the current user.
func (r *Resource) RemoveAllReactions(ctx context.Context, messageID string) error {
	e := endpoint.DeleteAllReactions(r.channelID, messageID)
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return discord.NewAPIError(resp)
	}
	return nil
}

// RemoveAllReactionsForEmoji removes all reactions for the given emoji on a message. This endpoint requires
// the 'MANAGE_MESSAGES' permission to be present on the current user.
func (r *Resource) RemoveAllReactionsForEmoji(ctx context.Context, messageID, emoji string) error {
	e := endpoint.DeleteAllReactionsForEmoji(r.channelID, messageID, url.PathEscape(emoji))
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return discord.NewAPIError(resp)
	}
	return nil
}

// Pins returns all pinned messages in the channel as an array of messages.
func (r *Resource) Pins(ctx context.Context) ([]discord.Message, error) {
	e := endpoint.GetPinnedMessages(r.channelID)
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var messages []discord.Message
	if err = json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, err
	}
	return messages, nil
}

// PinMessage pins a message in the channel. Requires the 'MANAGE_MESSAGES' permission.
func (r *Resource) PinMessage(ctx context.Context, id string) error {
	e := endpoint.AddPinnedChannelMessage(r.channelID, id)
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return discord.NewAPIError(resp)
	}
	return nil
}

// UnpinMessage deletes a pinned message in the channel. Requires the
// 'MANAGE_MESSAGES' permission.
func (r *Resource) UnpinMessage(ctx context.Context, id string) error {
	e := endpoint.DeletePinnedChannelMessage(r.channelID, id)
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return discord.NewAPIError(resp)
	}
	return nil
}
