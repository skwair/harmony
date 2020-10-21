package guild

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/skwair/harmony/discord"
	"github.com/skwair/harmony/internal/endpoint"
	"github.com/skwair/harmony/internal/rest"
)

// Integrations returns the list of integrations for the guild.
// Requires the 'MANAGE_GUILD' permission.
func (r *Resource) Integrations(ctx context.Context) ([]discord.GuildIntegration, error) {
	e := endpoint.GetGuildIntegrations(r.guildID)
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var integrations []discord.GuildIntegration
	if err = json.NewDecoder(resp.Body).Decode(&integrations); err != nil {
		return nil, err
	}
	return integrations, nil
}

// AddIntegration attaches an integration from the current user to the guild.
// Requires the 'MANAGE_GUILD' permission. Fires a Guild Integrations Update
// Gateway event.
func (r *Resource) AddIntegration(ctx context.Context, id, typ string) error {
	st := struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}{
		ID:   id,
		Type: typ,
	}
	b, err := json.Marshal(st)
	if err != nil {
		return err
	}

	e := endpoint.CreateGuildIntegration(r.guildID)
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

// ModifyIntegration modifies the behavior and settings of a guild integration.
// Requires the 'MANAGE_GUILD' permission. Fires a Guild Integrations Update Gateway event.
func (r *Resource) ModifyIntegration(ctx context.Context, id string, settings *discord.GuildIntegrationSettings) error {
	b, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	e := endpoint.ModifyGuildIntegration(r.guildID, id)
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

// RemoveIntegration removes the attached integration for the guild.
// Requires the 'MANAGE_GUILD' permission. Fires a Guild Integrations Update Gateway event.
func (r *Resource) RemoveIntegration(ctx context.Context, id string) error {
	e := endpoint.DeleteGuildIntegration(r.guildID, id)
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

// SyncIntegration syncs a guild integration. Requires the 'MANAGE_GUILD'
// permission.
func (r *Resource) SyncIntegration(ctx context.Context, id string) error {
	e := endpoint.SyncGuildIntegration(r.guildID, id)
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
