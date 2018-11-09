package discord

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/skwair/discord/internal/endpoint"
)

// ChannelType describes the type of a channel. Different fields
// are set or not depending on the channel's type.
type ChannelType int

// Supported channel types :
const (
	GuildText ChannelType = iota
	DM
	GuildVoice
	GroupDM
	GuildCategory
)

// Channel represents a guild or DM channel within Discord.
type Channel struct {
	ID                   string                `json:"id,omitempty"`
	Type                 ChannelType           `json:"type,omitempty"`
	GuildID              string                `json:"guild_id,omitempty"`
	Position             int                   `json:"position,omitempty"` // Sorting position of the channel.
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	Name                 string                `json:"name,omitempty"`
	Topic                string                `json:"topic,omitempty"`
	NSFW                 bool                  `json:"nsfw,omitempty"`
	LastMessageID        string                `json:"last_message_id,omitempty"`

	// For voice channels.
	Bitrate   int `json:"bitrate,omitempty"`
	UserLimit int `json:"user_limit,omitempty"`

	// For DMs.
	Recipients    []User `json:"recipients,omitempty"`
	Icon          string `json:"icon,omitempty"`
	OwnerID       string `json:"owner_id,omitempty"`
	ApplicationID string `json:"application_id,omitempty"` // Application id of the group DM creator if it is bot-created.

	ParentID         string    `json:"parent_id,omitempty"` // ID of the parent category for a channel.
	LastPinTimestamp time.Time `json:"last_pin_timestamp,omitempty"`
}

// PermissionOverwrite describes a specific permission that overwrites
// server-wide permissions.
type PermissionOverwrite struct {
	ID    string
	Type  string // Either "role" or "member".
	Allow int
	Deny  int
}

// GetChannel returns the channel object for the given id.
func (c *Client) GetChannel(id string) (*Channel, error) {
	e := endpoint.GetChannel(id)
	resp, err := c.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var channel Channel
	if err = json.NewDecoder(resp.Body).Decode(&channel); err != nil {
		return nil, err
	}
	return &channel, nil
}

// ChannelSettings describes a channel's settings.
type ChannelSettings struct {
	// Available for all channels :
	Name                 string                `json:"name,omitempty"`
	Position             int                   `json:"position,omitempty"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	// Available for text channels :
	Topic string `json:"topic,omitempty"`
	NSFW  bool   `json:"nsfw,omitempty"`
	// Available for audio channels :
	Bitrate   int `json:"bitrate,omitempty"`
	UserLimit int `json:"user_limit,omitempty"`
	// Available for text and audio channels :
	ParentID string `json:"parent_id,omitempty"`
}

// ModifyChannel updates a channel's settings given its ID and some new settings.
// Requires the 'MANAGE_CHANNELS' permission for the guild. Fires a Channel Update
// Gateway event. If modifying a category, individual Channel Update events will
// fire for each child channel that also changes.
func (c *Client) ModifyChannel(id string, s *ChannelSettings) (*Channel, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	e := endpoint.ModifyChannel(id)
	resp, err := c.doReq(http.MethodPatch, e, b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var ch Channel
	if err = json.NewDecoder(resp.Body).Decode(&ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

// DeleteChannel deletes a channel, or close a private message. Requires the 'MANAGE_CHANNELS'
// permission for the guild. Deleting a category does not delete its child channels; they will
// have their parent_id removed and a Channel Update Gateway event will fire for each of them.
// Returns a channel object on success. Fires a Channel Delete Gateway event.
func (c *Client) DeleteChannel(id string) error {
	e := endpoint.DeleteChannel(id)
	resp, err := c.doReq(http.MethodDelete, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return apiError(resp)
	}
	return nil
}

// GetChannelMessages returns the messages for a channel. If operating on a guild channel, this
// endpoint requires the 'VIEW_CHANNEL' permission to be present on the current user. If the
// current user is missing the 'READ_MESSAGE_HISTORY' permission in the channel then this will
// return no messages (since they cannot read the message history).
// The before, after, and around keys are mutually exclusive, only one may be passed at a time.
func (c *Client) GetChannelMessages(id string, limit int, around, before, after string) ([]Message, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	// NOTE: consider using a syntax like ><~messageID and have this method
	// take only 3 parameters :
	// - >messageID for 'after'
	// - <messageID for 'before'
	// - ~messageID for 'around'
	// - messageID or any other character could default to 'before' or
	// return an error.
	if around != "" {
		q.Set("around", around)
	}
	if before != "" {
		q.Set("before", before)
	}
	if after != "" {
		q.Set("after", after)
	}

	e := endpoint.GetChannelMessages(id, q.Encode())
	resp, err := c.doReq(http.MethodGet, e, nil)
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

// GetChannelMessage returns a specific message in the channel. If operating on a guild channel,
// this endpoints requires the 'READ_MESSAGE_HISTORY' permission to be present on the current user.
func (c *Client) GetChannelMessage(channelID, messageID string) (*Message, error) {
	e := endpoint.GetChannelMessage(channelID, messageID)
	resp, err := c.doReq(http.MethodGet, e, nil)
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

// UpdateChannelPermissions updates the channel permission overwrites for a user or role in a channel.
// typ is "member" if targetID is a user or "role" if it is a role.
// If the channel permission overwrites do not not exist, they are created.
// Only usable for guild channels. Requires the 'MANAGE_ROLES' permission.
func (c *Client) UpdateChannelPermissions(channelID, targetID string, allow, deny int, typ string) error {
	s := struct {
		ID    string `json:"id,omitempty"`
		Allow int    `json:"allow,omitempty"`
		Deny  int    `json:"deny,omitempty"`
		Type  string `json:"type,omitempty"`
	}{
		ID:    targetID,
		Allow: allow,
		Deny:  deny,
		Type:  typ,
	}

	b, err := json.Marshal(s)
	if err != nil {
		return err
	}

	e := endpoint.UpdateChannelPermissions(channelID, targetID)
	resp, err := c.doReq(http.MethodPut, e, b)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// DeleteChannelPermission deletes a channel permission overwrite for a user or role in a
// channel. Only usable for guild channels. Requires the 'MANAGE_ROLES' permission.
func (c *Client) DeleteChannelPermission(channelID, targetID string) error {
	e := endpoint.DeleteChannelPermission(channelID, targetID)
	resp, err := c.doReq(http.MethodDelete, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// GetChannelInvites returns a list of invite objects (with invite metadata) for the channel.
// Only usable for guild channels. Requires the 'MANAGE_CHANNELS' permission.
func (c *Client) GetChannelInvites(channelID string) ([]Invite, error) {
	e := endpoint.GetChannelInvites(channelID)
	resp, err := c.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var invites []Invite
	if err = json.NewDecoder(resp.Body).Decode(&invites); err != nil {
		return nil, err
	}
	return invites, nil
}

// CreateInvite allows to specify Invite settings when creating one.
type CreateInvite struct {
	MaxAge    int  `json:"max_age,omitempty"`
	MaxUses   int  `json:"max_uses,omitempty"`
	Temporary bool `json:"temporary,omitempty"` // Whether this invite only grants temporary membership.
	Unique    bool `json:"unique,omitempty"`
}

// CreateChannelInvite creates a new invite object for the channel. Only usable for guild channels.
// Requires the CREATE_INSTANT_INVITE permission.
func (c *Client) CreateChannelInvite(channelID string, i *CreateInvite) (*Invite, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	e := endpoint.CreateChannelInvite(channelID)
	resp, err := c.doReq(http.MethodPost, e, b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, apiError(resp)
	}

	var invite Invite
	if err = json.NewDecoder(resp.Body).Decode(&invite); err != nil {
		return nil, err
	}
	return &invite, nil
}

// AddRecipient adds a recipient to an existing Group DM or to a
// DM channel, creating a new Group DM channel.
// Groups have a limit of 10 recipients, including the current user.
func (c *Client) AddRecipient(channelID, recipientID string) error {
	e := endpoint.AddRecipient(channelID, recipientID)
	resp, err := c.doReq(http.MethodPut, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// RemoveRecipient removes a recipient from a Group DM
func (c *Client) RemoveRecipient(channelID, recipientID string) error {
	e := endpoint.RemoveRecipient(channelID, recipientID)
	resp, err := c.doReq(http.MethodDelete, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// TriggerTyping triggers a typing indicator for the specified channel.
// Generally bots should not implement this route. However, if a bot is
// responding to a command and expects the computation to take a few
// seconds, this endpoint may be called to let the user know that the
// bot is processing their message. Fires a Typing Start Gateway event.
func (c *Client) TriggerTyping(channelID string) error {
	e := endpoint.TriggerTyping(channelID)
	resp, err := c.doReq(http.MethodPost, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}
