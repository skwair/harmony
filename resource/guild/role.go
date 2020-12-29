package guild

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/skwair/harmony/discord"
	"github.com/skwair/harmony/internal/endpoint"
	"github.com/skwair/harmony/internal/rest"
)

// Roles returns a list of roles for the guild. Requires the 'MANAGE_ROLES'
// permission.
func (r *Resource) Roles(ctx context.Context) ([]discord.Role, error) {
	e := endpoint.GetGuildRoles(r.guildID)
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var roles []discord.Role
	if err = json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return nil, err
	}
	return roles, nil
}

// NewRole is like NewRoleWithReason but with no particular reason.
func (r *Resource) NewRole(ctx context.Context, settings *discord.RoleSettings) (*discord.Role, error) {
	return r.NewRoleWithReason(ctx, settings, "")
}

// NewRole creates a new role for the guild. Requires the 'MANAGE_ROLES'
// permission. Fires a Guild Role Create Gateway event.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) NewRoleWithReason(ctx context.Context, settings *discord.RoleSettings, reason string) (*discord.Role, error) {
	b, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	e := endpoint.CreateGuildRole(r.guildID)
	resp, err := r.client.DoWithHeader(ctx, e, rest.JSONPayload(b), rest.ReasonHeader(reason))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var role discord.Role
	if err = json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return nil, err
	}
	return &role, nil
}

// RolePosition is a pair of role ID with its position.
// A higher position means it will appear before in the list.
type RolePosition struct {
	ID       string `json:"id"`
	Position int    `json:"position"`
}

// ModifyRolePositions modifies the positions of a set of roles for the guild.
// Requires 'MANAGE_ROLES' permission. Fires multiple Guild Role Update Gateway events.
func (r *Resource) ModifyRolePositions(ctx context.Context, pos []RolePosition) ([]discord.Role, error) {
	b, err := json.Marshal(pos)
	if err != nil {
		return nil, err
	}

	e := endpoint.ModifyGuildRolePositions(r.guildID)
	resp, err := r.client.Do(ctx, e, rest.JSONPayload(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var roles []discord.Role
	if err = json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return nil, err
	}
	return roles, nil
}

// ModifyRole is like ModifyRoleWithReason but with no particular reason.
func (r *Resource) ModifyRole(ctx context.Context, id string, settings *discord.RoleSettings) (*discord.Role, error) {
	return r.ModifyRoleWithReason(ctx, id, settings, "")
}

// ModifyRole modifies a guild role. Requires the 'MANAGE_ROLES' permission.
// Fires a Guild Role Update Gateway event.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) ModifyRoleWithReason(ctx context.Context, id string, settings *discord.RoleSettings, reason string) (*discord.Role, error) {
	b, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	e := endpoint.ModifyGuildRole(r.guildID, id)
	resp, err := r.client.DoWithHeader(ctx, e, rest.JSONPayload(b), rest.ReasonHeader(reason))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var role discord.Role
	if err = json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return nil, err
	}
	return &role, nil
}

// DeleteRole is like DeleteRoleWithReason but with no particular reason.
func (r *Resource) DeleteRole(ctx context.Context, id string) error {
	return r.DeleteRoleWithReason(ctx, id, "")
}

// DeleteRole deletes a guild role. Requires the 'MANAGE_ROLES' permission.
// Fires a Guild Role Delete Gateway event.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) DeleteRoleWithReason(ctx context.Context, id, reason string) error {
	e := endpoint.DeleteGuildRole(r.guildID, id)
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

// AddMemberRole is like AddMemberRoleWithReason but with no particular reason.
func (r *Resource) AddMemberRole(ctx context.Context, userID, roleID string) error {
	return r.AddMemberRoleWithReason(ctx, userID, roleID, "")
}

// AddMemberRole adds a role to a guild member. Requires the 'MANAGE_ROLES'
// permission. Fires a Guild Member Update Gateway event.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) AddMemberRoleWithReason(ctx context.Context, userID, roleID, reason string) error {
	e := endpoint.AddGuildMemberRole(r.guildID, userID, roleID)
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

// RemoveMemberRole is like RemoveMemberRoleWithReason but with no particular reason.
func (r *Resource) RemoveMemberRole(ctx context.Context, userID, roleID string) error {
	return r.RemoveMemberRoleWithReason(ctx, userID, roleID, "")
}

// RemoveMemberRoleWithReason removes a role from a guild member. Requires the
// 'MANAGE_ROLES' permission. Fires a Guild Member Update Gateway event.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) RemoveMemberRoleWithReason(ctx context.Context, userID, roleID, reason string) error {
	e := endpoint.RemoveGuildMemberRole(r.guildID, userID, roleID)
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
