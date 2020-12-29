package guild

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/skwair/harmony/discord"
	"github.com/skwair/harmony/internal/endpoint"
	"github.com/skwair/harmony/internal/rest"
)

// Get returns the guild.
func (r *Resource) Get(ctx context.Context) (*discord.Guild, error) {
	e := endpoint.GetGuild(r.guildID)
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var g discord.Guild
	if err = json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return nil, err
	}
	return &g, nil
}

// Modify is like ModifyWithReason but with no particular reason.
func (r *Resource) Modify(ctx context.Context, settings *discord.GuildSettings) (*discord.Guild, error) {
	return r.ModifyWithReason(ctx, settings, "")
}

// ModifyWithReason modifies the guild's settings. Requires the 'MANAGE_GUILD' permission.
// Returns the updated guild on success. Fires a Guild Update Gateway event.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) ModifyWithReason(ctx context.Context, settings *discord.GuildSettings, reason string) (*discord.Guild, error) {
	b, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	e := endpoint.ModifyGuild(r.guildID)
	resp, err := r.client.DoWithHeader(ctx, e, rest.JSONPayload(b), rest.ReasonHeader(reason))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var g discord.Guild
	if err = json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return nil, err
	}
	return &g, nil
}

// Delete deletes the guild permanently. Current user must be owner.
// Fires a Guild Delete Gateway event.
func (r *Resource) Delete(ctx context.Context) error {
	e := endpoint.DeleteGuild(r.guildID)
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

// ChangeNick modifies the nickname of the current user (i.e.: the bot)
// for this guild.
// It returns the nickname on success. Requires the 'CHANGE_NICKNAME'
// permission. Fires a Guild Member Update Gateway event.
func (r *Resource) ChangeNick(ctx context.Context, name string) (string, error) {
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
	resp, err := r.client.Do(ctx, e, rest.JSONPayload(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", discord.NewAPIError(resp)
	}

	return st.Nick, nil
}

// PruneCount returns the number of members that would be removed in a prune
// operation. Requires the 'KICK_MEMBERS' permission.
func (r *Resource) PruneCount(ctx context.Context, days int) (int, error) {
	if days < 1 {
		days = 1
	}

	q := url.Values{}
	q.Set("days", strconv.Itoa(days))
	e := endpoint.GetGuildPruneCount(r.guildID, q.Encode())
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, discord.NewAPIError(resp)
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
func (r *Resource) BeginPrune(ctx context.Context, days int, computePruneCount bool) (pruneCount int, err error) {
	return r.BeginPruneWithReason(ctx, days, computePruneCount, "")
}

// BeginPruneWithReason begins a prune operation. Requires the 'KICK_MEMBERS' permission.
// Returns the number of members that were removed in the prune operation if
// computePruneCount is set to true (not recommended for large guilds).
// Fires multiple Guild Member Remove Gateway events.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) BeginPruneWithReason(ctx context.Context, days int, computePruneCount bool, reason string) (pruneCount int, err error) {
	if days < 1 {
		days = 1
	}

	q := url.Values{}
	q.Set("days", strconv.Itoa(days))
	q.Set("compute_prune_count", strconv.FormatBool(computePruneCount))
	e := endpoint.BeginGuildPrune(r.guildID, q.Encode())
	resp, err := r.client.DoWithHeader(ctx, e, nil, rest.ReasonHeader(reason))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, discord.NewAPIError(resp)
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
// Unlike the similar VoiceRegions method of the Client, this returns VIP
// servers when the guild is VIP-enabled.
func (r *Resource) VoiceRegions(ctx context.Context) ([]discord.VoiceRegion, error) {
	e := endpoint.GetGuildVoiceRegions(r.guildID)
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var regions []discord.VoiceRegion
	if err = json.NewDecoder(resp.Body).Decode(&regions); err != nil {
		return nil, err
	}
	return regions, nil
}

// Invites returns the list of invites (with invite metadata) for the guild.
// Requires the 'MANAGE_GUILD' permission.
func (r *Resource) Invites(ctx context.Context) ([]discord.Invite, error) {
	e := endpoint.GetGuildInvites(r.guildID)
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var invites []discord.Invite
	if err = json.NewDecoder(resp.Body).Decode(&invites); err != nil {
		return nil, err
	}
	return invites, nil
}

// Widget returns the guild's widget. Requires the 'MANAGE_GUILD' permission.
func (r *Resource) Widget(ctx context.Context) (*discord.GuildWidget, error) {
	e := endpoint.GetGuildWidget(r.guildID)
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var embed discord.GuildWidget
	if err = json.NewDecoder(resp.Body).Decode(&embed); err != nil {
		return nil, err
	}
	return &embed, nil
}

// ModifyWidget modifies the widget of the guild. Requires the
// 'MANAGE_GUILD' permission.
func (r *Resource) ModifyWidget(ctx context.Context, embed *discord.GuildWidget) (*discord.GuildWidget, error) {
	b, err := json.Marshal(embed)
	if err != nil {
		return nil, err
	}

	e := endpoint.ModifyGuildWidget(r.guildID)
	resp, err := r.client.Do(ctx, e, rest.JSONPayload(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(embed); err != nil {
		return nil, err
	}
	return embed, nil
}

// VanityURL returns a partial invite for the guild if that feature is
// enabled. Requires the 'MANAGE_GUILD' permission.
func (r *Resource) VanityURL(ctx context.Context) (*discord.Invite, error) {
	e := endpoint.GetGuildVanityURL(r.guildID)
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var invite discord.Invite
	if err = json.NewDecoder(resp.Body).Decode(&invite); err != nil {
		return nil, err
	}
	return &invite, nil
}

// Webhooks returns the list of webhooks in the guild.
// Requires the 'MANAGE_WEBHOOKS' permission.
func (r *Resource) Webhooks(ctx context.Context) ([]discord.Webhook, error) {
	e := endpoint.GetGuildWebhooks(r.guildID)
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var webhooks []discord.Webhook
	if err = json.NewDecoder(resp.Body).Decode(&webhooks); err != nil {
		return nil, err
	}
	return webhooks, nil
}
