package harmony

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/skwair/harmony/integration"

	"github.com/skwair/harmony/internal/endpoint"
)

type Integration struct {
	ID                string              `json:"id"`
	Name              string              `json:"name"`
	Type              string              `json:"type"`
	Enabled           bool                `json:"enabled"`
	Syncing           bool                `json:"syncing"`
	RoleID            string              `json:"role_id"`
	ExpireBehavior    int                 `json:"expire_behavior"`
	ExpireGravePeriod int                 `json:"expire_grave_period"`
	User              *User               `json:"user"`
	Account           *IntegrationAccount `json:"account"`
	SyncedAt          time.Time           `json:"synced_at"`
}

type IntegrationAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Integrations returns the list of integrations for the guild.
// Requires the 'MANAGE_GUILD' permission.
func (r *GuildResource) Integrations(ctx context.Context) ([]Integration, error) {
	e := endpoint.GetGuildIntegrations(r.guildID)
	resp, err := r.client.doReq(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var integrations []Integration
	if err = json.NewDecoder(resp.Body).Decode(&integrations); err != nil {
		return nil, err
	}
	return integrations, nil
}

// AddIntegration attaches an integration from the current user to the guild.
// Requires the 'MANAGE_GUILD' permission. Fires a Guild Integrations Update
// Gateway event.
func (r *GuildResource) AddIntegration(ctx context.Context, id, typ string) error {
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
	resp, err := r.client.doReq(ctx, e, b)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// ModifyIntegration modifies the behavior and settings of a guild integration.
// Requires the 'MANAGE_GUILD' permission. Fires a Guild Integrations Update Gateway event.
func (r *GuildResource) ModifyIntegration(ctx context.Context, id string, settings *integration.Settings) error {
	b, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	e := endpoint.ModifyGuildIntegration(r.guildID, id)
	resp, err := r.client.doReq(ctx, e, b)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// RemoveIntegration removes the attached integration for the guild.
// Requires the 'MANAGE_GUILD' permission. Fires a Guild Integrations Update Gateway event.
func (r *GuildResource) RemoveIntegration(ctx context.Context, id string) error {
	e := endpoint.DeleteGuildIntegration(r.guildID, id)
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

// SyncIntegration syncs a guild integration. Requires the 'MANAGE_GUILD'
// permission.
func (r *GuildResource) SyncIntegration(ctx context.Context, id string) error {
	e := endpoint.SyncGuildIntegration(r.guildID, id)
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
