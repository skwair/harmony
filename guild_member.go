package harmony

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/skwair/harmony/guild"
	"github.com/skwair/harmony/internal/endpoint"
)

// GuildMember represents a User in a Guild.
type GuildMember struct {
	User     *User     `json:"user,omitempty"`
	Nick     string    `json:"nick,omitempty"`
	Roles    []string  `json:"roles,omitempty"` // Role IDs.
	JoinedAt time.Time `json:"joined_at,omitempty"`
	Deaf     bool      `json:"deaf,omitempty"`
	Mute     bool      `json:"mute,omitempty"`
}

// PermissionsIn returns the permissions of the Guild member in the given Guild and channel.
func (m *GuildMember) PermissionsIn(g *Guild, ch *Channel) (permissions int) {
	base := computeBasePermissions(g, m)
	return computeOverwrites(ch, m, base)
}

// HasRole returns whether this member has the given role.
// Note that this method does not try to fetch this member latest roles, it instead looks
// in the roles it already had when this member object was created.
func (m *GuildMember) HasRole(id string) bool {
	for _, roleID := range m.Roles {
		if roleID == id {
			return true
		}
	}
	return false
}

// Member returns a single guild member given its user ID.
func (r *GuildResource) Member(ctx context.Context, userID string) (*GuildMember, error) {
	e := endpoint.GetGuildMember(r.guildID, userID)
	resp, err := r.client.doReq(ctx, http.MethodGet, e, nil)
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

// Members returns a list of at most limit guild members, starting at after.
// limit must be between 1 and 1000 and will be set to those values if higher/lower.
// after is the ID of the guild member you want to get the list from, leave it
// empty to start from the beginning.
func (r *GuildResource) Members(ctx context.Context, limit int, after string) ([]GuildMember, error) {
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

	e := endpoint.ListGuildMembers(r.guildID, q.Encode())
	resp, err := r.client.doReq(ctx, http.MethodGet, e, nil)
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

// AddMember adds a user to the guild, provided you have a valid oauth2 access
// token for the user with the guilds.join scope. Fires a Guild Member Add Gateway event.
// Requires the bot to have the CREATE_INSTANT_INVITE permission.
func (r *GuildResource) AddMember(ctx context.Context, userID, token string, settings *guild.MemberSettings) (*GuildMember, error) {
	st := struct {
		AccessToken string `json:"access_token,omitempty"`
		*guild.MemberSettings
	}{
		AccessToken:    token,
		MemberSettings: settings,
	}
	b, err := json.Marshal(st)
	if err != nil {
		return nil, err
	}

	e := endpoint.AddGuildMember(r.guildID, userID)
	resp, err := r.client.doReq(ctx, http.MethodPatch, e, b)
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

// RemoveMember removes the given user from the guild. Requires 'KICK_MEMBERS'
// permission. Fires a Guild Member Remove Gateway event.
func (r *GuildResource) RemoveMember(ctx context.Context, userID string) error {
	e := endpoint.RemoveGuildMember(r.guildID, userID)
	resp, err := r.client.doReq(ctx, http.MethodDelete, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// ModifyMember modifies attributes of a guild member. Fires a Guild Member
// Update Gateway event.
func (r *GuildResource) ModifyMember(ctx context.Context, userID string, settings *guild.MemberSettings) error {
	b, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	e := endpoint.ModifyGuildMember(r.guildID, userID)
	resp, err := r.client.doReq(ctx, http.MethodPatch, e, b)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}
