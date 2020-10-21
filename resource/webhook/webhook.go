package webhook

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/skwair/harmony/discord"
	"github.com/skwair/harmony/internal/endpoint"
	"github.com/skwair/harmony/internal/rest"
)

// Get returns the webhook.
func (r *Resource) Get(ctx context.Context) (*discord.Webhook, error) {
	e := endpoint.GetWebhook(r.webhookID)
	resp, err := r.client.Do(ctx, e, nil)
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

// Modify is like ModifyWithReason but with no particular reason.
func (r *Resource) Modify(ctx context.Context, settings *discord.WebhookSettings) (*discord.Webhook, error) {
	return r.ModifyWithReason(ctx, settings, "")
}

// ModifyWithReason modifies the webhook. Requires the 'MANAGE_WEBHOOKS' permission.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) ModifyWithReason(ctx context.Context, settings *discord.WebhookSettings, reason string) (*discord.Webhook, error) {
	b, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	e := endpoint.ModifyWebhook(r.webhookID)
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

// Delete is like DeleteWithReason but with no particular reason.
func (r *Resource) Delete(ctx context.Context) error {
	return r.DeleteWithReason(ctx, "")
}

// DeleteWithReason deletes the webhook.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) DeleteWithReason(ctx context.Context, reason string) error {
	e := endpoint.DeleteWebhook(r.webhookID)
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
