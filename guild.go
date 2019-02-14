package harmony

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/skwair/harmony/channel"
	"github.com/skwair/harmony/guild"
	"github.com/skwair/harmony/internal/endpoint"
)

// Guild in Discord represents an isolated collection of users and channels,
// and are often referred to as "servers" in the UI.
type Guild struct {
	ID                          string                         `json:"id"`
	Name                        string                         `json:"name,omitempty"`
	Icon                        *string                        `json:"icon,omitempty"`
	Splash                      *string                        `json:"splash,omitempty"`
	Owner                       bool                           `json:"owner,omitempty"`
	OwnerID                     string                         `json:"owner_id,omitempty"`
	Permissions                 int                            `json:"permissions,omitempty"`
	Region                      string                         `json:"region,omitempty"`
	AFKChannelID                *string                        `json:"afk_channel_id,omitempty"`
	AFKTimeout                  int                            `json:"afk_timeout,omitempty"`
	EmbedEnabled                bool                           `json:"embed_enabled,omitempty"`
	EmbedChannelID              string                         `json:"embed_channel_id,omitempty"`
	VerificationLevel           guild.VerificationLevel        `json:"verification_level,omitempty"`
	DefaultMessageNotifications guild.DefaultNotificationLevel `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       guild.ExplicitContentFilter    `json:"explicit_content_filter,omitempty"`
	Roles                       []Role                         `json:"roles,omitempty"`
	Emojis                      []Emoji                        `json:"emojis,omitempty"`
	Features                    []string                       `json:"features,omitempty"`
	MFALevel                    int                            `json:"mfa_level,omitempty"`
	ApplicationID               *string                        `json:"application_id,omitempty"`
	WidgetEnabled               bool                           `json:"widget_enabled,omitempty"`
	WidgetChannelID             string                         `json:"widget_channel_id,omitempty"`
	SystemChannelID             *string                        `json:"system_channel_id,omitempty"`

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

// CreateGuild creates a new guild with the given name.
// Returns the created guild on success. Fires a Guild Create Gateway event.
func (c *Client) CreateGuild(ctx context.Context, name string) (*Guild, error) {
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
	resp, err := c.doReq(ctx, e, jsonPayload(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, apiError(resp)
	}

	var g Guild
	if err = json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return nil, err
	}
	return &g, nil
}

// GuildResource is a resource that allows to perform various actions on a Discord guild.
// Create one with Client.Guild.
type GuildResource struct {
	guildID string
	client  *Client
}

// Guild returns a new guild resource to manage the guild with the given ID.
func (c *Client) Guild(id string) *GuildResource {
	return &GuildResource{guildID: id, client: c}
}

// Get returns the guild.
func (r *GuildResource) Get(ctx context.Context) (*Guild, error) {
	e := endpoint.GetGuild(r.guildID)
	resp, err := r.client.doReq(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var g Guild
	if err = json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return nil, err
	}
	return &g, nil
}

// Modify is like ModifyWithReason but with no particular reason.
func (r *GuildResource) Modify(ctx context.Context, settings *guild.Settings) (*Guild, error) {
	return r.ModifyWithReason(ctx, settings, "")
}

// ModifyWithReason modifies the guild's settings. Requires the 'MANAGE_GUILD' permission.
// Returns the updated guild on success. Fires a Guild Update Gateway event.
// The given reason will be set in the audit log entry for this action.
func (r *GuildResource) ModifyWithReason(ctx context.Context, settings *guild.Settings, reason string) (*Guild, error) {
	b, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	e := endpoint.ModifyGuild(r.guildID)
	resp, err := r.client.doReqWithHeader(ctx, e, jsonPayload(b), reasonHeader(reason))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var g Guild
	if err = json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return nil, err
	}
	return &g, nil
}

// Delete deletes the guild permanently. Current user must be owner.
// Fires a Guild Delete Gateway event.
func (r *GuildResource) Delete(ctx context.Context) error {
	e := endpoint.DeleteGuild(r.guildID)
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

// Channels returns the list of channels in the guild.
func (r *GuildResource) Channels(ctx context.Context) ([]Channel, error) {
	e := endpoint.GetGuildChannels(r.guildID)
	resp, err := r.client.doReq(ctx, e, nil)
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

// NewChannel is like NewChannelWithReason but with no particular reason.
func (r *GuildResource) NewChannel(ctx context.Context, settings *channel.Settings) (*Channel, error) {
	return r.NewChannelWithReason(ctx, settings, "")
}

// NewChannelWithReason creates a new channel in the guild. Requires the MANAGE_CHANNELS permission.
// Fires a Channel Create Gateway event.
// The given reason will be set in the audit log entry for this action.
func (r *GuildResource) NewChannelWithReason(ctx context.Context, settings *channel.Settings, reason string) (*Channel, error) {
	b, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	e := endpoint.CreateGuildChannel(r.guildID)
	resp, err := r.client.doReqWithHeader(ctx, e, jsonPayload(b), reasonHeader(reason))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, apiError(resp)
	}

	var ch Channel
	if err = json.NewDecoder(resp.Body).Decode(&ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

// ChannelPosition is a pair of channel ID with its position.
type ChannelPosition struct {
	ID       string `json:"id"`
	Position int    `json:"position"`
}

// ModifyChannelPosition modifies the positions of a set of channel for the guild.
// Requires 'MANAGE_CHANNELS' permission. Fires multiple Channel Update Gateway events.
//
// Only channels to be modified are required, with the minimum being a swap between at
// least two channels.
func (r *GuildResource) ModifyChannelPosition(ctx context.Context, pos []ChannelPosition) error {
	b, err := json.Marshal(pos)
	if err != nil {
		return err
	}

	e := endpoint.ModifyChannelPositions(r.guildID)
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

// ChangeNick modifies the nickname of the current user in the guild.
// It returns the nickname on success. Requires the 'CHANGE_NICKNAME'
// permission. Fires a Guild Member Update Gateway event.
func (r *GuildResource) ChangeNick(ctx context.Context, name string) (string, error) {
	st := struct {
		Nick string `json:"nick"`
	}{
		Nick: name,
	}
	b, err := json.Marshal(st)
	if err != nil {
		return "", err
	}

	e := endpoint.ModifyCurrentUserNick(r.guildID)
	resp, err := r.client.doReq(ctx, e, jsonPayload(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", apiError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}
	return st.Nick, nil
}

// PruneCount returns the number of members that would be removed in a prune
// operation. Requires the 'KICK_MEMBERS' permission.
func (r *GuildResource) PruneCount(ctx context.Context, days int) (int, error) {
	if days < 1 {
		days = 1
	}

	q := url.Values{}
	q.Set("days", strconv.Itoa(days))
	e := endpoint.GetGuildPruneCount(r.guildID, q.Encode())
	resp, err := r.client.doReq(ctx, e, nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, apiError(resp)
	}

	st := struct {
		Pruned int `json:"pruned"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&st); err != nil {
		return 0, err
	}
	return st.Pruned, nil
}

// BeginPrune is like BeginPruneWithReason but with no particular reason.
func (r *GuildResource) BeginPrune(ctx context.Context, days int, computePruneCount bool) (pruneCount int, err error) {
	return r.BeginPruneWithReason(ctx, days, computePruneCount, "")
}

// BeginPruneWithReason begins a prune operation. Requires the 'KICK_MEMBERS' permission.
// Returns the number of members that were removed in the prune operation if
// computePruneCount is set to true (not recommended for large guilds).
// Fires multiple Guild Member Remove Gateway events.
// The given reason will be set in the audit log entry for this action.
func (r *GuildResource) BeginPruneWithReason(ctx context.Context, days int, computePruneCount bool, reason string) (pruneCount int, err error) {
	if days < 1 {
		days = 1
	}

	q := url.Values{}
	q.Set("days", strconv.Itoa(days))
	q.Set("compute_prune_count", strconv.FormatBool(computePruneCount))
	e := endpoint.BeginGuildPrune(r.guildID, q.Encode())
	resp, err := r.client.doReqWithHeader(ctx, e, nil, reasonHeader(reason))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, apiError(resp)
	}

	var st struct {
		Pruned *int `json:"pruned"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&st); err != nil {
		return 0, err
	}

	if st.Pruned != nil {
		pruneCount = *st.Pruned
	}

	return pruneCount, nil
}

// VoiceRegions returns a list of available voice regions for the guild.
// Unlike the similar GetVoiceRegions method of the Client, this returns VIP
// servers when the guild is VIP-enabled.
func (r *GuildResource) VoiceRegions(ctx context.Context) ([]VoiceRegion, error) {
	e := endpoint.GetGuildVoiceRegions(r.guildID)
	resp, err := r.client.doReq(ctx, e, nil)
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

// Invites returns the list of invites (with invite metadata) for the guild.
// Requires the 'MANAGE_GUILD' permission.
func (r *GuildResource) Invites(ctx context.Context) ([]Invite, error) {
	e := endpoint.GetGuildInvites(r.guildID)
	resp, err := r.client.doReq(ctx, e, nil)
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

// Embed returns the guild's embed. Requires the 'MANAGE_GUILD' permission.
func (r *GuildResource) Embed(ctx context.Context) (*guild.Embed, error) {
	e := endpoint.GetGuildEmbed(r.guildID)
	resp, err := r.client.doReq(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var embed guild.Embed
	if err = json.NewDecoder(resp.Body).Decode(&embed); err != nil {
		return nil, err
	}
	return &embed, nil
}

// ModifyEmbed modifies the guild embed of the guild. Requires the
// 'MANAGE_GUILD' permission.
func (r *GuildResource) ModifyEmbed(ctx context.Context, embed *guild.Embed) (*guild.Embed, error) {
	b, err := json.Marshal(embed)
	if err != nil {
		return nil, err
	}

	e := endpoint.ModifyGuildEmbed(r.guildID)
	resp, err := r.client.doReq(ctx, e, jsonPayload(b))
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

// VanityURL returns a partial invite for the guild if that feature is
// enabled. Requires the 'MANAGE_GUILD' permission.
func (r *GuildResource) VanityURL(ctx context.Context) (*Invite, error) {
	e := endpoint.GetGuildVanityURL(r.guildID)
	resp, err := r.client.doReq(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var invite Invite
	if err = json.NewDecoder(resp.Body).Decode(&invite); err != nil {
		return nil, err
	}
	return &invite, nil
}

// Webhooks returns the list of webhooks in the guild.
// Requires the 'MANAGE_WEBHOOKS' permission.
func (r *GuildResource) Webhooks(ctx context.Context) ([]Webhook, error) {
	e := endpoint.GetGuildWebhooks(r.guildID)
	return r.client.getWebhooks(ctx, e)
}

type requestGuildMembers struct {
	GuildID string `json:"guild_id"`
	Query   string `json:"query"`
	Limit   int    `json:"limit"`
}

// RequestGuildMembers is used to request offline members for the guild. When initially
// connecting, the gateway will only send offline members if a guild has less than
// the large_threshold members (value in the Gateway Identify). If a client wishes
// to receive additional members, they need to explicitly request them via this
// operation. The server will send Guild Members Chunk events in response with up
// to 1000 members per chunk until all members that match the request have been sent.
// query is a string that username starts with, or an empty string to return all members.
// limit is the maximum number of members to send or 0 to request all members matched.
// You need to be connected to the Gateway to call this method, else it will
// return ErrGatewayNotConnected.
func (r *GuildResource) RequestGuildMembers(query string, limit int) error {
	if !r.client.isConnected() {
		return ErrGatewayNotConnected
	}

	return r.client.sendPayload(gatewayOpcodeRequestGuildMembers, &requestGuildMembers{
		GuildID: r.guildID,
		Query:   query,
		Limit:   limit,
	})
}
