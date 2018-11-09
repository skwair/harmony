package discord

import (
	"encoding/json"
	"net/http"

	"github.com/skwair/discord/internal/endpoint"
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

// CreateRole describes the creation of a new guild role.
type CreateRole struct {
	Name        string `json:"name,omitempty"`
	Permissions int    `json:"permissions,omitempty"`
	Color       int    `json:"color,omitempty"`
	Hoist       bool   `json:"hoist,omitempty"`
	Mentionable bool   `json:"mentionable,omitempty"`
}

// GetRoles returns a list of roles for the given guild. Requires the
// 'MANAGE_ROLES' permission.
func (c *Client) GetRoles(guildID string) ([]Role, error) {
	e := endpoint.GetRoles(guildID)
	resp, err := c.doReq(http.MethodGet, e, nil)
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

// CreateRole creates a new role for the given guild. Requires the 'MANAGE_ROLES'
// permission. Fires a Guild Role Create Gateway event.
func (c *Client) CreateRole(guildID string, r *CreateRole) (*Role, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	e := endpoint.CreateRole(guildID)
	resp, err := c.doReq(http.MethodPost, e, b)
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

// ModifyRolePositions modifies the positions of a set of roles for the given guild.
// Requires 'MANAGE_ROLES' permission. Fires multiple Guild Role Update Gateway events.
func (c *Client) ModifyRolePositions(guildID string, positions []RolePosition) ([]Role, error) {
	b, err := json.Marshal(positions)
	if err != nil {
		return nil, err
	}

	e := endpoint.ModifyRolePositions(guildID)
	resp, err := c.doReq(http.MethodPatch, e, b)
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

// ModifyRole describes how to modify a guild role. All fields are optional.
type ModifyRole struct {
	Name        string `json:"name,omitempty"`
	Permissions int    `json:"permissions,omitempty"`
	Color       int    `json:"color,omitempty"`
	Hoist       bool   `json:"hoist,omitempty"`
	Mentionable bool   `json:"mentionable,omitempty"`
}

// ModifyRole modifies a guild role. Requires the 'MANAGE_ROLES' permission.
// Fires a Guild Role Update Gateway event.
func (c *Client) ModifyRole(guildID, roleID string, r *ModifyRole) (*Role, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	e := endpoint.ModifyRole(guildID, roleID)
	resp, err := c.doReq(http.MethodPatch, e, b)
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
func (c *Client) DeleteRole(guildID, roleID string) error {
	e := endpoint.DeleteRole(guildID, roleID)
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

// AddGuildMemberRole adds a role to a guild member. Requires the 'MANAGE_ROLES'
// permission. Fires a Guild Member Update Gateway event.
func (c *Client) AddGuildMemberRole(guildID, userID, roleID string) error {
	e := endpoint.AddGuildMemberRole(guildID, userID, roleID)
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

// RemoveGuildMemberRole removes a role from a guild member. Requires the
// 'MANAGE_ROLES' permission. Fires a Guild Member Update Gateway event.
func (c *Client) RemoveGuildMemberRole(guildID, userID, roleID string) error {
	e := endpoint.RemoveGuildMemberRole(guildID, userID, roleID)
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
