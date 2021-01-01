package channel

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/skwair/harmony/discord"
	"github.com/skwair/harmony/internal/endpoint"
	"github.com/skwair/harmony/internal/rest"
)

// Get returns the channel.
func (r *Resource) Get(ctx context.Context) (*discord.Channel, error) {
	e := endpoint.GetChannel(r.channelID)
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var ch discord.Channel
	if err = json.NewDecoder(resp.Body).Decode(&ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

// Modify is like ModifyWithReason but with no particular reason.
func (r *Resource) Modify(ctx context.Context, settings *discord.ChannelSettings) (*discord.Channel, error) {
	return r.ModifyWithReason(ctx, settings, "")
}

// ModifyWithReason updates the channel's settings. Requires the 'MANAGE_CHANNELS'
// permission for the discord. Fires a Channel Update Gateway event. If modifying
// category, individual Channel Update events will fire for each child channel
// that also changes.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) ModifyWithReason(ctx context.Context, settings *discord.ChannelSettings, reason string) (*discord.Channel, error) {
	b, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	e := endpoint.ModifyChannel(r.channelID)
	resp, err := r.client.DoWithHeader(ctx, e, rest.JSONPayload(b), rest.ReasonHeader(reason))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var ch discord.Channel
	if err = json.NewDecoder(resp.Body).Decode(&ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

// Delete is like DeleteWithReason but with no particular reason.
func (r *Resource) Delete(ctx context.Context) (*discord.Channel, error) {
	return r.DeleteWithReason(ctx, "")
}

// DeleteWithReason deletes the channel, or closes the private message. Requires the 'MANAGE_CHANNELS'
// permission for the discord. Deleting a category does not delete its child channels; they will
// have their parent_id removed and a Channel Update Gateway event will fire for each of them.
// Returns the deleted channel on success. Fires a Channel Delete Gateway event.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) DeleteWithReason(ctx context.Context, reason string) (*discord.Channel, error) {
	e := endpoint.DeleteChannel(r.channelID)
	resp, err := r.client.DoWithHeader(ctx, e, nil, rest.ReasonHeader(reason))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var ch discord.Channel
	if err = json.NewDecoder(resp.Body).Decode(&ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

// UpdatePermissions is like UpdatePermissionsWithReason but with no particular reason.
func (r *Resource) UpdatePermissions(ctx context.Context, perms discord.PermissionOverwrite) error {
	return r.UpdatePermissionsWithReason(ctx, perms, "")
}

// UpdatePermissionsWithReason updates the channel permission overwrites for a user or
// role in the channel.
// If the channel permission overwrites do not not exist, they are created.
// Only usable for guild channels. Requires the 'MANAGE_ROLES' permission.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) UpdatePermissionsWithReason(ctx context.Context, perms discord.PermissionOverwrite, reason string) error {
	b, err := json.Marshal(perms)
	if err != nil {
		return err
	}

	e := endpoint.EditChannelPermissions(r.channelID, perms.ID)
	resp, err := r.client.DoWithHeader(ctx, e, rest.JSONPayload(b), rest.ReasonHeader(reason))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return discord.NewAPIError(resp)
	}
	return nil
}

// DeletePermission is like DeletePermissionWithReason but with no particular reason.
func (r *Resource) DeletePermission(ctx context.Context, channelID, targetID string) error {
	return r.DeletePermissionWithReason(ctx, channelID, targetID, "")
}

// DeletePermissionWithReason deletes the channel permission overwrite for a user or
// role in a channel. Only usable for guild channels. Requires the 'MANAGE_ROLES'
// permission.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) DeletePermissionWithReason(ctx context.Context, channelID, targetID, reason string) error {
	e := endpoint.DeleteChannelPermission(channelID, targetID)
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

// Invites returns a list of invites (with invite metadata) for the channel.
// Only usable for guild channels. Requires the 'MANAGE_CHANNELS' permission.
func (r *Resource) Invites(ctx context.Context) ([]discord.Invite, error) {
	e := endpoint.GetChannelInvites(r.channelID)
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

// NewInvite is like NewInviteWithReason but with no particular reason.
func (r *Resource) NewInvite(ctx context.Context, settings *discord.InviteSettings) (*discord.Invite, error) {
	return r.NewInviteWithReason(ctx, settings, "")
}

// NewInviteWithReason creates a new invite for the channel. Only usable
// for guild channels. Requires the CREATE_INSTANT_INVITE permission.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) NewInviteWithReason(ctx context.Context, settings *discord.InviteSettings, reason string) (*discord.Invite, error) {
	b, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	e := endpoint.CreateChannelInvite(r.channelID)
	resp, err := r.client.DoWithHeader(ctx, e, rest.JSONPayload(b), rest.ReasonHeader(reason))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var i discord.Invite
	if err = json.NewDecoder(resp.Body).Decode(&i); err != nil {
		return nil, err
	}
	return &i, nil
}

// TriggerTyping triggers a typing indicator for the channel.
// Generally bots should not implement this route. However, if a bot is
// responding to a command and expects the computation to take a few
// seconds, this endpoint may be called to let the user know that the
// bot is processing their message. Fires a Typing Start Gateway event.
func (r *Resource) TriggerTyping(ctx context.Context) error {
	e := endpoint.TriggerTypingIndicator(r.channelID)
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

// NewWebhook is like NewWebhookWithReason but with no particular reason.
func (r *Resource) NewWebhook(ctx context.Context, name, avatar string) (*discord.Webhook, error) {
	return r.NewWebhookWithReason(ctx, name, avatar, "")
}

// NewWebhookWithReason creates a new webhook for the channel. Requires the 'MANAGE_WEBHOOKS'
// permission.
// name must contain between 2 and 32 characters. avatar is an avatar data string,
// see https://discord.com/developers/docs/resources/user#avatar-data for more info.
// It can be left empty to have the default avatar.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) NewWebhookWithReason(ctx context.Context, name, avatar, reason string) (*discord.Webhook, error) {
	st := struct {
		Name   string `json:"name,omitempty"`
		Avatar string `json:"avatar,omitempty"`
	}{
		Name:   name,
		Avatar: avatar,
	}
	b, err := json.Marshal(st)
	if err != nil {
		return nil, err
	}

	e := endpoint.CreateWebhook(r.channelID)
	resp, err := r.client.DoWithHeader(ctx, e, rest.JSONPayload(b), rest.ReasonHeader(reason))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var w discord.Webhook
	if err = json.NewDecoder(resp.Body).Decode(&w); err != nil {
		return nil, err
	}
	return &w, nil
}

// Webhooks returns webhooks for the channel.
func (r *Resource) Webhooks(ctx context.Context) ([]discord.Webhook, error) {
	e := endpoint.GetChannelWebhooks(r.channelID)
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
