package harmony

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/skwair/harmony/internal/endpoint"
	"github.com/skwair/harmony/role"
)

// Role represents a set of permissions attached to a group of users.
// Roles have unique names, colors, and can be "pinned" to the side bar,
// causing their members to be listed separately. Roles are unique per guild,
// and can have separate permission profiles for the global context (guild)
// and channel context.
type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Color       int    `json:"color"`    // Integer representation of hexadecimal color code.
	Hoist       bool   `json:"hoist"`    // Whether this role is pinned in the user listing.
	Position    int    `json:"position"` // Integer	position of this role.
	Permissions int    `json:"permissions"`
	Managed     bool   `json:"managed"` // Whether this role is managed by an integration.
	Mentionable bool   `json:"mentionable"`
}

// Roles returns a list of roles for the guild. Requires the 'MANAGE_ROLES'
// permission.
func (r *GuildResource) Roles(ctx context.Context) ([]Role, error) {
	e := endpoint.GetGuildRoles(r.guildID)
	resp, err := r.client.doReq(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var roles []Role
	if err = json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return nil, err
	}
	return roles, nil
}

// NewRole creates a new role for the guild. Requires the 'MANAGE_ROLES'
// permission. Fires a Guild Role Create Gateway event.
func (r *GuildResource) NewRole(ctx context.Context, settings *role.Settings) (*Role, error) {
	b, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	e := endpoint.CreateGuildRole(r.guildID)
	resp, err := r.client.doReq(ctx, e, b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var role Role
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
func (r *GuildResource) ModifyRolePositions(ctx context.Context, pos []RolePosition) ([]Role, error) {
	b, err := json.Marshal(pos)
	if err != nil {
		return nil, err
	}

	e := endpoint.ModifyGuildRolePositions(r.guildID)
	resp, err := r.client.doReq(ctx, e, b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var roles []Role
	if err = json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return nil, err
	}
	return roles, nil
}

// ModifyRole modifies a guild role. Requires the 'MANAGE_ROLES' permission.
// Fires a Guild Role Update Gateway event.
func (r *GuildResource) ModifyRole(ctx context.Context, id string, settings *role.Settings) (*Role, error) {
	b, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	e := endpoint.ModifyGuildRole(r.guildID, id)
	resp, err := r.client.doReq(ctx, e, b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var role Role
	if err = json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return nil, err
	}
	return &role, nil
}

// DeleteRole deletes a guild role. Requires the 'MANAGE_ROLES' permission.
// Fires a Guild Role Delete Gateway event.
func (r *GuildResource) DeleteRole(ctx context.Context, id string) error {
	e := endpoint.DeleteGuildRole(r.guildID, id)
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

// AddMemberRole adds a role to a guild member. Requires the 'MANAGE_ROLES'
// permission. Fires a Guild Member Update Gateway event.
func (r *GuildResource) AddMemberRole(ctx context.Context, userID, roleID string) error {
	e := endpoint.AddGuildMemberRole(r.guildID, userID, roleID)
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

// RemoveMemberRole removes a role from a guild member. Requires the
// 'MANAGE_ROLES' permission. Fires a Guild Member Update Gateway event.
func (r *GuildResource) RemoveMemberRole(ctx context.Context, userID, roleID string) error {
	e := endpoint.RemoveGuildMemberRole(r.guildID, userID, roleID)
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
