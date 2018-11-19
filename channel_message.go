package discord

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/skwair/discord/embed"
	"github.com/skwair/discord/internal/endpoint"
)

// MessageType describes the type of a message. Different fields
// are set or not depending on the message's type.
type MessageType int

// supported message types :
const (
	Default MessageType = iota
	RecipientAdd
	RecipientRemove
	Call
	ChannelNameChange
	ChannelIconChange
	ChannelPinnedMessage
	GuildMemberJoin
)

// Message represents a message sent in a channel within Discord.
// The author object follows the structure of the user object, but is
// only a valid user in the case where the message is generated by a
// user or bot user. If the message is generated by a webhook, the
// author object corresponds to the webhook's id, username, and avatar.
// You can tell if a message is generated by a webhook by checking for
// the webhook_id on the message object.
type Message struct {
	ID              string        `json:"id"`
	ChannelID       string        `json:"channel_id"`
	GuildID         string        `json:"guild_id"`
	Author          *User         `json:"author"`
	Content         string        `json:"content"`
	Timestamp       time.Time     `json:"timestamp"`
	EditedTimestamp time.Time     `json:"edited_timestamp"`
	TTS             bool          `json:"tts"`
	MentionEveryone bool          `json:"mention_everyone"`
	Mentions        []User        `json:"mentions"`
	MentionRoles    []string      `json:"mention_roles"` // Role IDs
	Attachments     []Attachment  `json:"attachments"`   // Any attached files.
	Embeds          []embed.Embed `json:"embeds"`        // Any embedded content.
	Reactions       []Reaction    `json:"reactions"`
	Nonce           string        `json:"nonce"` // Used for validating a message was sent.
	Pinned          bool          `json:"pinned"`
	WebhookID       string        `json:"webhook_id"`
	Type            MessageType   `json:"type"`

	// Sent with Rich Presence-related chat embeds.
	Activity    *MessageActivity    `json:"activity"`
	Application *MessageApplication `json:"application"`
}

// Attachment is file attached to a message.
type Attachment struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Size     int    `json:"size"`
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

// Messages returns messages in the channel. If operating on a guild channel, this
// endpoint requires the 'VIEW_CHANNEL' permission to be present on the current user. If the
// current user is missing the 'READ_MESSAGE_HISTORY' permission in the channel then this will
// return no messages (since they cannot read the message history).
// The query parameter is a message ID prefixed with one of the following character :
//	- '>' for fetching messages after
//	- '<' for fetching messages before
//	- '~' for fetching messages around
// For example, to retrieve 50 messages around (25 before, 25 after) a message having the
// ID 221588207995121520, set query to "~221588207995121520".
// Limit is a positive integer between 1 and 100 that default to 50 if set to 0.
func (r *ChannelResource) Messages(query string, limit int) ([]Message, error) {
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
	resp, err := r.client.doReq(http.MethodGet, e, nil)
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
func (r *ChannelResource) Message(id string) (*Message, error) {
	e := endpoint.GetChannelMessage(r.channelID, id)
	resp, err := r.client.doReq(http.MethodGet, e, nil)
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

// DeleteMessage deletes a message. If operating on a guild channel and trying to delete a
// message that was not sent by the current user, this endpoint requires the 'MANAGE_MESSAGES'
// permission. Fires a Message Delete Gateway event.
func (r *ChannelResource) DeleteMessage(messageID string) error {
	e := endpoint.DeleteMessage(r.channelID, messageID)
	resp, err := r.client.doReq(http.MethodDelete, e, nil)
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
func (r *ChannelResource) DeleteMessageBulk(messageIDs []string) error {
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
	resp, err := r.client.doReq(http.MethodPost, e, b)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// createMessage describes a message creation.
type createMessage struct {
	Content string       `json:"content,omitempty"` // Up to 2000 characters.
	Nonce   string       `json:"nonce,omitempty"`
	TTS     bool         `json:"tts,omitempty"`
	Embed   *embed.Embed `json:"embed,omitempty"`
}

// json implements the multipartPayload interface so createMessage can be used as
// a payload with the multipartFromFiles method.
func (cm *createMessage) json() ([]byte, error) {
	return json.Marshal(cm)
}

func (c *Client) sendMessage(channelID string, msg *createMessage) (*Message, error) {
	if msg.Embed != nil && msg.Embed.Type == "" {
		msg.Embed.Type = "rich"
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	e := endpoint.CreateMessage(channelID)
	resp, err := c.doReq(http.MethodPost, e, b)
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

// SendMessage is like SendMessageWithOptions with an empty nonce and text to speech disabled.
func (r *ChannelResource) SendMessage(text string) (*Message, error) {
	return r.SendMessageWithOptions(text, "", false)
}

// SendMessageWithOptions posts a message to the channel. If operating on a guild channel,
// this endpoint requires the 'SEND_MESSAGES' permission to be present on the
// current user. If the tts field is set to true, the 'SEND_TTS_MESSAGES' permission is
// required for the message to be spoken. Returns the message sent.
// Fires a Message Create Gateway event.
// Before using this endpoint, you must connect to the gateway at least once.
// The nonce will be returned in the result and also transmitted to other clients.
// You can set it to empty if you do not need this feature.
func (r *ChannelResource) SendMessageWithOptions(text, nonce string, tts bool) (*Message, error) {
	return r.client.sendMessage(r.channelID, &createMessage{
		Content: text,
		Nonce:   nonce,
		TTS:     tts,
	})
}

// SendEmbed is like SendEmbedWithOptions with no text, an empty nonce and text to speech disabled.
func (r *ChannelResource) SendEmbed(embed *embed.Embed) (*Message, error) {
	return r.SendEmbedWithOptions(embed, "", "", false)
}

// SendEmbedWithOptions sends some embedded rich content attached to a message to the channel.
// See SendMessageWithOptions for required permissions and the embed sub package for more information
// about embeds.
func (r *ChannelResource) SendEmbedWithOptions(embed *embed.Embed, text, nonce string, tts bool) (*Message, error) {
	return r.client.sendMessage(r.channelID, &createMessage{
		Content: text,
		Nonce:   nonce,
		TTS:     tts,
		Embed:   embed,
	})
}

// File is a file along with its name. It is used to send files
// to channels with SendFiles.
type File struct {
	Name   string
	Reader io.Reader
}

// SendFiles is like SendFilesWithOptions with no text, an empty nonce, no
// embed and text to speech disabled.
func (r *ChannelResource) SendFiles(files ...File) (*Message, error) {
	return r.SendFilesWithOptions("", "", nil, false, files...)
}

// SendFilesWithOptions sends some attached files with an optional text and/or embedded rich
// content to the channel.
// See SendMessageWithOptions for required permissions and the embed sub package for more information
// about embeds.
func (r *ChannelResource) SendFilesWithOptions(text, nonce string, embed *embed.Embed, tts bool, files ...File) (*Message, error) {
	if len(files) < 1 {
		return nil, ErrNoFileProvided
	}

	cm := &createMessage{
		Content: text,
		Embed:   embed,
		Nonce:   nonce,
		TTS:     tts,
	}
	b, h, err := multipartFromFiles(cm, files...)
	if err != nil {
		return nil, err
	}

	e := endpoint.CreateMessage(r.channelID)
	resp, err := r.client.doReqWithHeader(http.MethodPost, e, b, h)
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
func (r *ChannelResource) EditMessage(messageID, content string) (*Message, error) {
	return r.client.editMessage(r.channelID, messageID, &editMessage{content, nil})
}

// EditEmbed is like EditMessage but with embedded content support.
func (r *ChannelResource) EditEmbed(messageID, content string, embed *embed.Embed) (*Message, error) {
	return r.client.editMessage(r.channelID, messageID, &editMessage{content, embed})
}

func (c *Client) editMessage(channelID, messageID string, edit *editMessage) (*Message, error) {
	b, err := json.Marshal(edit)
	if err != nil {
		return nil, err
	}

	e := endpoint.EditMessage(channelID, messageID)
	resp, err := c.doReq(http.MethodPatch, e, b)
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

// GetReactions returns a list of users that reacted to a message with the given emoji.
// limit is the number of users to return and can be set to any value ranging from 1 to 100.
// If set to 0, it defaults to 25. If more than 100 users reacted with the given emoji,
// the before and after parameters can be used to fetch more users.
func (r *ChannelResource) GetReactions(messageID, emoji string, limit int, before, after string) ([]User, error) {
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
	resp, err := r.client.doReq(http.MethodGet, e, nil)
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

// CreateReaction adds a reaction to a message in the channel. This endpoint requires
// the 'READ_MESSAGE_HISTORY' permission to be present on the current user. Additionally,
// if nobody else has reacted to the message using this emoji, this endpoint requires
// the'ADD_REACTIONS' permission to be present on the current user.
func (r *ChannelResource) CreateReaction(messageID, emoji string) error {
	e := endpoint.CreateReaction(r.channelID, messageID, url.PathEscape(emoji))
	resp, err := r.client.doReq(http.MethodPut, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// DeleteReaction deletes a reaction the current user has made for the message.
func (r *ChannelResource) DeleteReaction(messageID, emoji string) error {
	e := endpoint.DeleteOwnReaction(r.channelID, messageID, url.PathEscape(emoji))
	resp, err := r.client.doReq(http.MethodDelete, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// DeleteUserReaction deletes another user's reaction. This endpoint requires the
// 'MANAGE_MESSAGES' permission to be present on the current user.
func (r *ChannelResource) DeleteUserReaction(messageID, userID, emoji string) error {
	e := endpoint.DeleteUserReaction(r.channelID, messageID, userID, url.PathEscape(emoji))
	resp, err := r.client.doReq(http.MethodDelete, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// DeleteAllReactions deletes all reactions on a message. This endpoint requires the
// 'MANAGE_MESSAGES' permission to be present on the current user.
func (r *ChannelResource) DeleteAllReactions(messageID string) error {
	e := endpoint.DeleteAllReactions(r.channelID, messageID)
	resp, err := r.client.doReq(http.MethodDelete, e, nil)
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
func (r *ChannelResource) Pins() ([]Message, error) {
	e := endpoint.GetPinnedMessages(r.channelID)
	resp, err := r.client.doReq(http.MethodGet, e, nil)
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
func (r *ChannelResource) PinMessage(id string) error {
	e := endpoint.AddPinnedChannelMessage(r.channelID, id)
	resp, err := r.client.doReq(http.MethodPut, e, nil)
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
func (r *ChannelResource) UnpinMessage(id string) error {
	e := endpoint.DeletePinnedChannelMessage(r.channelID, id)
	resp, err := r.client.doReq(http.MethodDelete, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}