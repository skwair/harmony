package discord

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/skwair/discord/internal/endpoint"
	"github.com/skwair/discord/optional"
)

// Guild in Discord represents an isolated collection of users and channels,
// and are often referred to as "servers" in the UI.
type Guild struct {
	ID                          string   `json:"id"`
	Name                        string   `json:"name,omitempty"`
	Icon                        *string  `json:"icon,omitempty"`
	Splash                      *string  `json:"splash,omitempty"`
	Owner                       bool     `json:"owner,omitempty"`
	OwnerID                     string   `json:"owner_id,omitempty"`
	Permissions                 int      `json:"permissions,omitempty"`
	Region                      string   `json:"region,omitempty"`
	AFKChannelID                *string  `json:"afk_channel_id,omitempty"`
	AFKTimeout                  int      `json:"afk_timeout,omitempty"`
	EmbedEnabled                bool     `json:"embed_enabled,omitempty"`
	EmbedChannelID              string   `json:"embed_channel_id,omitempty"`
	VerificationLevel           int      `json:"verification_level,omitempty"`
	DefaultMessageNotifications int      `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       int      `json:"explicit_content_filter,omitempty"`
	Roles                       []Role   `json:"roles,omitempty"`
	Emojis                      []Emoji  `json:"emojis,omitempty"`
	Features                    []string `json:"features,omitempty"`
	MFALevel                    int      `json:"mfa_level,omitempty"`
	ApplicationID               *string  `json:"application_id,omitempty"`
	WidgetEnabled               bool     `json:"widget_enabled,omitempty"`
	WidgetChannelID             string   `json:"widget_channel_id,omitempty"`
	SystemChannelID             *string  `json:"system_channel_id,omitempty"`

	// Following fields are only sent within the GUILD_CREATE event.
	JoinedAt    time.Time     `json:"joined_at,omitempty"`
	Large       bool          `json:"large,omitempty"`
	Unavailable bool          `json:"unavailable,omitempty"`
	MemberCount int           `json:"member_count,omitempty"`
	VoiceStates []VoiceState  `json:"voice_states,omitempty"`
	Members     []GuildMember `json:"members,omitempty"`
	Channels    []Channel     `json:"channels,omitempty"`
	Presences   []Presence    `json:"presences,omitempty"`
}

// Presence is a user's current state on a guild.
// This event is sent when a user's presence is updated for a guild.
type Presence struct {
	User    *User     `json:"user,omitempty"`
	Roles   []string  `json:"roles,omitempty"` // Array of IDs.
	Game    *Activity `json:"game,omitempty"`
	GuildID string    `json:"guild_id,omitempty"`
	Status  string    `json:"status,omitempty"` // Either "idle", "dnd", "online", or "offline".
}

// PartialGuild is a subset of the Guild object, returned by the Discord API
// when fetching current user's guilds.
type PartialGuild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Owner       bool   `json:"owner"`
	Permissions int    `json:"permissions"`
}

// UnavailableGuild is a Guild that is not available, either because there is a
// guild outage or because the connected user was removed from this guild.
type UnavailableGuild struct {
	ID          string `json:"id"`
	Unavailable *bool  `json:"unavailable"` // If not set, the connected user was removed from this Guild.
}

// GuildMember represents a User in a Guild.
type GuildMember struct {
	User     *User     `json:"user,omitempty"`
	Nick     string    `json:"nick,omitempty"`
	Roles    []string  `json:"roles,omitempty"` // Role IDs.
	JoinedAt time.Time `json:"joined_at,omitempty"`
	Deaf     bool      `json:"deaf,omitempty"`
	Mute     bool      `json:"mute,omitempty"`
}

// CreateGuild creates a new guild with the given name.
// Returns a guild object on success. Fires a Guild Create Gateway event.
func (c *Client) CreateGuild(name string) (*Guild, error) {
	s := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	e := endpoint.CreateGuild()
	resp, err := c.doReq(http.MethodPost, e, b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, apiError(resp)
	}

	var guild Guild
	if err = json.NewDecoder(resp.Body).Decode(&guild); err != nil {
		return nil, err
	}
	return &guild, nil
}

// NOTE: maybe expose a CreateGuildWithParams method that allows to create a guild
// with custom settings without requiring a ModifyGuild call after CreateGuild.

// GetGuild returns the guild object for the given id.
func (c *Client) GetGuild(id string) (*Guild, error) {
	e := endpoint.GetGuild(id)
	resp, err := c.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var guild Guild
	if err = json.NewDecoder(resp.Body).Decode(&guild); err != nil {
		return nil, err
	}
	return &guild, nil
}

// GuildSettings are the settings of a guild, all fields are optional and only those
// explicitly set will be modified.
type GuildSettings struct {
	Name                        *optional.String `json:"name,omitempty"`
	Region                      *optional.String `json:"region,omitempty"`
	VerificationLevel           *optional.Int    `json:"verification_level,omitempty"`
	DefaultMessageNotifications *optional.Int    `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *optional.Int    `json:"explicit_content_filter,omitempty"`
	AfkChannelID                *optional.String `json:"afk_channel_id,omitempty"`
	AfkTimeout                  *optional.Int    `json:"afk_timeout,omitempty"`
	Icon                        *optional.String `json:"icon,omitempty"`
	OwnerID                     *optional.String `json:"owner_id,omitempty"`
	Splash                      *optional.String `json:"splash,omitempty"`
	SystemChannelID             *optional.String `json:"system_channel_id,omitempty"`
}

// ModifyGuild modifies a guild's settings. Requires the 'MANAGE_GUILD' permission.
// Returns the updated guild object on success. Fires a Guild Update Gateway event.
func (c *Client) ModifyGuild(guildID string, s *GuildSettings) (*Guild, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	e := endpoint.ModifyGuild(guildID)
	resp, err := c.doReq(http.MethodPatch, e, b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var guild Guild
	if err = json.NewDecoder(resp.Body).Decode(&guild); err != nil {
		return nil, err
	}
	return &guild, nil
}

// DeleteGuild deletes a guild permanently. User must be owner.
// Fires a Guild Delete Gateway event.
func (c *Client) DeleteGuild(id string) error {
	e := endpoint.DeleteGuild(id)
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

// GetGuildChannels returns the list of channels in the given guild.
func (c *Client) GetGuildChannels(guildID string) ([]Channel, error) {
	e := endpoint.GetGuildChannels(guildID)
	resp, err := c.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var channels []Channel
	if err = json.NewDecoder(resp.Body).Decode(&channels); err != nil {
		return nil, err
	}
	return channels, nil
}

// createChannel describes a channel creation.
type createChannel struct {
	Name        string                `json:"name,omitempty"`
	Type        ChannelType           `json:"type,omitempty"`
	Bitrate     int                   `json:"bitrate,omitempty"`
	UserLimit   int                   `json:"user_limit,omitempty"`
	Permissions []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID    string                `json:"parent_id,omitempty"`
	NSFW        bool                  `json:"nsfw,omitempty"`
}

func (c *Client) createGuildChannel(guildID string, ch *createChannel) (*Channel, error) {
	b, err := json.Marshal(ch)
	if err != nil {
		return nil, err
	}

	e := endpoint.CreateGuildChannel(guildID)
	resp, err := c.doReq(http.MethodPost, e, b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, apiError(resp)
	}

	var channel Channel
	if err = json.NewDecoder(resp.Body).Decode(&channel); err != nil {
		return nil, err
	}
	return &channel, nil
}

// CreateTextChannel creates a text channel in the given guild.
// Requires the 'MANAGE_CHANNELS' permission. Fires a Channel Create Gateway event.
func (c *Client) CreateTextChannel(
	guildID string,
	name string,
	permissions []PermissionOverwrite,
	parentID string,
	nsfw bool) (*Channel, error) {
	return c.createGuildChannel(guildID, &createChannel{
		Name:        name,
		Permissions: permissions,
		Type:        GuildText,
		ParentID:    parentID,
		NSFW:        nsfw,
	})
}

// CreateVoiceChannel creates a voice channel in the given guild.
// Requires the 'MANAGE_CHANNELS' permission. Fires a Channel Create Gateway event.
func (c *Client) CreateVoiceChannel(
	guildID string,
	name string,
	permissions []PermissionOverwrite,
	parentID string,
	bitrate int,
	userLimit int) (*Channel, error) {
	return c.createGuildChannel(guildID, &createChannel{
		Name:        name,
		Permissions: permissions,
		Type:        GuildVoice,
		ParentID:    parentID,
		Bitrate:     bitrate,
		UserLimit:   userLimit,
	})
}

// CreateChannelCategory creates a new channel category in the given guild.
// Requires the 'MANAGE_CHANNELS' permission. Fires a Channel Create Gateway event.
func (c *Client) CreateChannelCategory(
	guildID string,
	name string,
	permissions []PermissionOverwrite) (*Channel, error) {
	return c.createGuildChannel(guildID, &createChannel{
		Name:        name,
		Permissions: permissions,
		Type:        GuildCategory,
	})
}

// ChannelPosition is a pair of channel ID with its position.
type ChannelPosition struct {
	ID       string `json:"id"`
	Position int    `json:"position"`
}

// ModifyChannelPositions modifies the positions of a set of channel for the given guild.
// Requires 'MANAGE_CHANNELS' permission. Fires multiple Channel Update Gateway events.
//
// Only channels to be modified are required, with the minimum being a swap between at
// least two channels.
func (c *Client) ModifyChannelPositions(guildID string, positions []ChannelPosition) error {
	b, err := json.Marshal(positions)
	if err != nil {
		return err
	}

	e := endpoint.ModifyChannelPositions(guildID)
	resp, err := c.doReq(http.MethodPatch, e, b)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// GetGuildMember returns a single guild member.
func (c *Client) GetGuildMember(guildID, userID string) (*GuildMember, error) {
	e := endpoint.GetGuildMember(guildID, userID)
	resp, err := c.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var m GuildMember
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

// GetGuildMembers returns a list of at most limit guild members, starting at after.
// limit must be between 1 and 1000 and will be set to those values if higher/lower.
// after is the ID of the guild member you want to get the list from, leave it
// empty to start from the beginning.
func (c *Client) GetGuildMembers(guildID string, limit int, after string) ([]GuildMember, error) {
	if limit < 1 {
		limit = 1
	}
	if limit > 1000 {
		limit = 1000
	}

	q := url.Values{}
	q.Set("limit", strconv.Itoa(limit))
	if after != "" {
		q.Set("after", after)
	}

	e := endpoint.GetGuildMembers(guildID, q.Encode())
	resp, err := c.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var members []GuildMember
	if err = json.NewDecoder(resp.Body).Decode(&members); err != nil {
		return nil, err
	}
	return members, nil
}

// AddGuildMember adds a user to the given guild, provided you have a valid oauth2 access
// token for the user with the guilds.join scope. Fires a Guild Member Add Gateway event.
// Requires the bot to have the CREATE_INSTANT_INVITE permission.
func (c *Client) AddGuildMember(guildID, userID, accessToken string, s *GuildMemberSettings) (*GuildMember, error) {
	st := struct {
		AccessToken string `json:"access_token,omitempty"`
		*GuildMemberSettings
	}{
		AccessToken:         accessToken,
		GuildMemberSettings: s,
	}
	b, err := json.Marshal(st)
	if err != nil {
		return nil, err
	}

	e := endpoint.AddGuildMember(guildID, userID)
	resp, err := c.doReq(http.MethodPatch, e, b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var member GuildMember
	if err = json.NewDecoder(resp.Body).Decode(&member); err != nil {
		return nil, err
	}
	return &member, nil
}

// RemoveGuildMember removes the given user from the given guild. Requires
// 'KICK_MEMBERS' permission. Fires a Guild Member Remove Gateway event.
func (c *Client) RemoveGuildMember(guildID, userID string) error {
	e := endpoint.RemoveGuildMember(guildID, userID)
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

// GuildMemberSettings are the settings of a guild member, all fields are optional
// and only those explicitly set will be modified.
type GuildMemberSettings struct {
	Nick  *optional.String      `json:"nick,omitempty"`
	Roles *optional.StringSlice `json:"roles,omitempty"`
	Mute  *optional.Bool        `json:"mute,omitempty"`
	Deaf  *optional.Bool        `json:"deaf,omitempty"`
	// ID of channel to move user to (if they are connected to voice).
	ChannelID *optional.String `json:"channel_id,omitempty"`
}

// ModifyGuildMember modifies attributes of a guild member. Fires a Guild Member
// Update Gateway event.
func (c *Client) ModifyGuildMember(guildID, userID string, s *GuildMemberSettings) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}

	e := endpoint.ModifyGuildMember(guildID, userID)
	resp, err := c.doReq(http.MethodPatch, e, b)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// ModifyCurrentUserNick modifies the nickname of the current user in the given
// guild. It returns the nickname on success. Requires the 'CHANGE_NICKNAME'
// permission. Fires a Guild Member Update Gateway event.
func (c *Client) ModifyCurrentUserNick(guildID, nick string) (string, error) {
	s := struct {
		Nick string `json:"nick"`
	}{
		Nick: nick,
	}
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	e := endpoint.ModifyCurrentUserNick(guildID)
	resp, err := c.doReq(http.MethodPatch, e, b)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", apiError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return "", err
	}
	return s.Nick, nil
}

// Ban represents a Guild ban.
type Ban struct {
	Reason string
	User   *User
}

// GetGuildBans returns a list of ban objects for the users banned from this guild.
// Requires the 'BAN_MEMBERS' permission.
func (c *Client) GetGuildBans(guildID string) ([]Ban, error) {
	e := endpoint.GetGuildBans(guildID)
	resp, err := c.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var bans []Ban
	if err = json.NewDecoder(resp.Body).Decode(&bans); err != nil {
		return nil, err
	}
	return bans, nil
}

// Ban is a shorthand to ban a user with no reason and without
// deleting his messages. Requires the 'BAN_MEMBERS' permission.
// For more control, use the BanWithReason method.
func (c *Client) Ban(guildID, userID string) error {
	return c.BanWithReason(guildID, userID, 0, "")
}

// BanWithReason creates a guild ban, and optionally delete previous messages
// sent by the banned user. Requires the 'BAN_MEMBERS' permission.
// Parameter delMsgDays is the number of days to delete messages for (0-7).
// Fires a Guild Ban Add Gateway event.
func (c *Client) BanWithReason(guildID, userID string, delMsgDays int, reason string) error {
	q := url.Values{}
	if reason != "" {
		q.Set("reason", reason)
	}
	if delMsgDays > 0 {
		q.Set("delete-message-days", strconv.Itoa(delMsgDays))
	}

	e := endpoint.BanWithReason(guildID, userID, q.Encode())
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

// Unban removes the ban for a user. Requires the 'BAN_MEMBERS' permissions.
// Fires a Guild Ban Remove Gateway event.
func (c *Client) Unban(guildID, userID string) error {
	e := endpoint.Unban(guildID, userID)
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

// GetPruneCount returns the number of members that would be removed in a prune
// operation. Requires the 'KICK_MEMBERS' permission.
func (c *Client) GetPruneCount(guildID string, days int) (int, error) {
	if days < 1 {
		days = 1
	}

	q := url.Values{}
	q.Set("days", strconv.Itoa(days))
	e := endpoint.GetPruneCount(guildID, q.Encode())
	resp, err := c.doReq(http.MethodGet, e, nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, apiError(resp)
	}

	s := struct {
		Pruned int `json:"pruned"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return 0, err
	}
	return s.Pruned, nil
}

// BeginPrune begins a prune operation. Requires the 'KICK_MEMBERS' permission.
// Returns the number of members that were removed in the prune operation.
// Fires multiple Guild Member Remove Gateway events.
func (c *Client) BeginPrune(guildID string, days int) (int, error) {
	if days < 1 {
		days = 1
	}

	q := url.Values{}
	q.Set("days", strconv.Itoa(days))
	e := endpoint.BeginPrune(guildID, q.Encode())
	resp, err := c.doReq(http.MethodPost, e, nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, apiError(resp)
	}

	s := struct {
		Pruned int `json:"pruned"`
	}{100}
	if err = json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return 0, err
	}
	return s.Pruned, nil
}

// GetGuildVoiceRegions returns a list of available voice regions for the given guild.
// Unlike the similar GetVoiceRegions, this returns VIP servers when the guild is VIP-enabled.
func (c *Client) GetGuildVoiceRegions(guildID string) ([]VoiceRegion, error) {
	e := endpoint.GetGuildVoiceRegions(guildID)
	resp, err := c.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var regions []VoiceRegion
	if err = json.NewDecoder(resp.Body).Decode(&regions); err != nil {
		return nil, err
	}
	return regions, nil
}

// GetGuildInvites returns the list of invites (with invite metadata) for the given guild.
// Requires the 'MANAGE_GUILD' permission.
func (c *Client) GetGuildInvites(guildID string) ([]Invite, error) {
	e := endpoint.GetGuildInvites(guildID)
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

// GetGuildIntegrations returns the list of integrations for the given guild.
// Requires the 'MANAGE_GUILD' permission.
func (c *Client) GetGuildIntegrations(guildID string) ([]Integration, error) {
	e := endpoint.GetGuildIntegrations(guildID)
	resp, err := c.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var integration []Integration
	if err = json.NewDecoder(resp.Body).Decode(&integration); err != nil {
		return nil, err
	}
	return integration, nil
}

// AddGuildIntegration attaches an integration from the current user to the given guild.
// Requires the 'MANAGE_GUILD' permission. Fires a Guild Integrations Update Gateway event.
func (c *Client) AddGuildIntegration(guildID, integrationID, integrationType string) error {
	s := struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}{
		ID:   integrationID,
		Type: integrationType,
	}
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}

	e := endpoint.AddGuildIntegration(guildID)
	resp, err := c.doReq(http.MethodPost, e, b)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// IntegrationSettings describes a guild integration's settings.
type IntegrationSettings struct {
	// The behavior when an integration subscription lapses.
	ExpireBehavior int
	// Period (in seconds) where the integration will ignore lapsed subscriptions.
	ExpireGracePeriod int
	// Whether emoticons should be synced for this integration (twitch only currently).
	EnableEmoticons bool
}

// ModifyGuildIntegration modifies the behavior and settings of a guild integration.
// Requires the 'MANAGE_GUILD' permission. Fires a Guild Integrations Update Gateway event.
func (c *Client) ModifyGuildIntegration(guildID, integrationID string, s *IntegrationSettings) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}

	e := endpoint.ModifyGuildIntegration(guildID, integrationID)
	resp, err := c.doReq(http.MethodPost, e, b)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// RemoveGuildIntegration removes the attached integration for the given guild.
// Requires the 'MANAGE_GUILD' permission. Fires a Guild Integrations Update Gateway event.
func (c *Client) RemoveGuildIntegration(guildID, integrationID string) error {
	e := endpoint.RemoveGuildIntegration(guildID, integrationID)
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

// SyncGuildIntegration syncs a guild integration. Requires the 'MANAGE_GUILD'
// permission.
func (c *Client) SyncGuildIntegration(guildID, integrationID string) error {
	e := endpoint.SyncGuildIntegration(guildID, integrationID)
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

type GuildEmbed struct {
	Enabled   bool   `json:"enabled,omitempty"`
	ChannelID string `json:"channel_id,omitempty"`
}

// GetGuildEmbed returns the given guild's embed. Requires the 'MANAGE_GUILD'
// permission.
func (c *Client) GetGuildEmbed(guildID string) (*GuildEmbed, error) {
	e := endpoint.GetGuildEmbed(guildID)
	resp, err := c.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var embed GuildEmbed
	if err = json.NewDecoder(resp.Body).Decode(&embed); err != nil {
		return nil, err
	}
	return &embed, nil
}

// ModifyGuildEmbed modifies a guild embed for the given guild. Requires the
// 'MANAGE_GUILD' permission.
func (c *Client) ModifyGuildEmbed(guildID string, embed *GuildEmbed) (*GuildEmbed, error) {
	b, err := json.Marshal(embed)
	if err != nil {
		return nil, err
	}

	e := endpoint.ModifyGuildEmbed(guildID)
	resp, err := c.doReq(http.MethodPatch, e, b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(embed); err != nil {
		return nil, err
	}
	return embed, nil
}

// GetGuildVanityURL returns a partial invite for guilds with that feature
// enabled. Requires the 'MANAGE_GUILD' permission.
func (c *Client) GetGuildVanityURL(guildID string) (*Invite, error) {
	e := endpoint.GetGuildVanityURL(guildID)
	resp, err := c.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var invite Invite
	if err = json.NewDecoder(resp.Body).Decode(invite); err != nil {
		return nil, err
	}
	return &invite, nil
}

type requestGuildMembers struct {
	GuildID string `json:"guild_id"`
	Query   string `json:"query"`
	Limit   int    `json:"limit"`
}

// RequestGuildMembers is used to request offline members for a guild. When initially
// connecting, the gateway will only send offline members if a guild has less than
// the large_threshold members (value in the Gateway Identify). If a client wishes
// to receive additional members, they need to explicitly request them via this
// operation. The server will send Guild Members Chunk events in response with up
// to 1000 members per chunk until all members that match the request have been sent.
// query is a string that username starts with, or an empty string to return all members.
// limit is the maximum number of members to send or 0 to request all members matched.
// You need to be connected to the Gateway to call this method, else it will
// return ErrGatewayNotConnected.
func (c *Client) RequestGuildMembers(guildID, query string, limit int) error {
	if !c.isConnected() {
		return ErrGatewayNotConnected
	}

	return c.sendPayload(8, &requestGuildMembers{
		GuildID: guildID,
		Query:   query,
		Limit:   limit,
	})
}
