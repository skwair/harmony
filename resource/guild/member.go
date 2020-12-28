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

// Member returns a single guild member given its user ID.
func (r *Resource) Member(ctx context.Context, userID string) (*discord.GuildMember, error) {
	e := endpoint.GetGuildMember(r.guildID, userID)
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var m discord.GuildMember
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

// Members returns a list of at most limit guild members, starting at after.
// limit must be between 1 and 1000 and will be set to those values if higher/lower.
// after is the ID of the guild member you want to get the list from, leave it
// empty to start from the beginning.
func (r *Resource) Members(ctx context.Context, limit int, after string) ([]discord.GuildMember, error) {
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
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var members []discord.GuildMember
	if err = json.NewDecoder(resp.Body).Decode(&members); err != nil {
		return nil, err
	}
	return members, nil
}

// AddMember adds a user to the guild, provided you have a valid oauth2 access
// token for the user with the guilds.join scope. Fires a Guild Member Add Gateway event.
// Requires the bot to have the CREATE_INSTANT_INVITE permission.
func (r *Resource) AddMember(ctx context.Context, userID, token string, settings *discord.GuildMemberSettings) (*discord.GuildMember, error) {
	st := struct {
		AccessToken string `json:"access_token,omitempty"`
		*discord.GuildMemberSettings
	}{
		AccessToken:         token,
		GuildMemberSettings: settings,
	}
	b, err := json.Marshal(st)
	if err != nil {
		return nil, err
	}

	e := endpoint.AddGuildMember(r.guildID, userID)
	resp, err := r.client.Do(ctx, e, rest.JSONPayload(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var member discord.GuildMember
	if err = json.NewDecoder(resp.Body).Decode(&member); err != nil {
		return nil, err
	}
	return &member, nil
}

// ModifyMember is like ModifyMemberWithReason but with no particular reason.
func (r *Resource) ModifyMember(ctx context.Context, userID string, settings *discord.GuildMemberSettings) error {
	return r.ModifyMemberWithReason(ctx, userID, settings, "")
}

// ModifyMember modifies attributes of a guild member. Fires a Guild Member
// Update Gateway event.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) ModifyMemberWithReason(ctx context.Context, userID string, settings *discord.GuildMemberSettings, reason string) error {
	b, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	e := endpoint.ModifyGuildMember(r.guildID, userID)
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

// Kick is like KickWithReason but with no particular reason.
func (r *Resource) Kick(ctx context.Context, userID string) error {
	return r.KickWithReason(ctx, userID, "")
}

// Kick removes the given user from the guild. Requires 'KICK_MEMBERS'
// permission. Fires a Guild Member Remove Gateway event.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) KickWithReason(ctx context.Context, userID, reason string) error {
	e := endpoint.RemoveGuildMember(r.guildID, userID)
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

// Bans returns a list of bans for the users banned from this guild.
// Requires the 'BAN_MEMBERS' permission.
func (r *Resource) Bans(ctx context.Context) ([]discord.Ban, error) {
	e := endpoint.GetGuildBans(r.guildID)
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var bans []discord.Ban
	if err = json.NewDecoder(resp.Body).Decode(&bans); err != nil {
		return nil, err
	}
	return bans, nil
}

// Ban is a shorthand to ban a user with no reason and without
// deleting his messages. Requires the 'BAN_MEMBERS' permission.
// For more control, use the BanWithReason method.
func (r *Resource) Ban(ctx context.Context, userID string) error {
	return r.BanWithReason(ctx, userID, 0, "")
}

// BanWithReason creates a guild ban, and optionally delete previous messages
// sent by the banned user. Requires the 'BAN_MEMBERS' permission.
// Parameter delMsgDays is the number of days to delete messages for (0-7).
// Fires a Guild Ban Add Gateway event.
func (r *Resource) BanWithReason(ctx context.Context, userID string, delMsgDays int, reason string) error {
	q := url.Values{}
	if reason != "" {
		q.Set("reason", reason)
	}
	if delMsgDays > 0 {
		q.Set("delete_message_days", strconv.Itoa(delMsgDays))
	}

	e := endpoint.CreateGuildBan(r.guildID, userID, q.Encode())
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

// Unban is like UnbanWithReason but with no particular reason.
func (r *Resource) Unban(ctx context.Context, userID string) error {
	return r.UnbanWithReason(ctx, userID, "")
}

// Unban removes the ban for a user. Requires the 'BAN_MEMBERS' permissions.
// Fires a Guild Ban Remove Gateway event.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) UnbanWithReason(ctx context.Context, userID, reason string) error {
	e := endpoint.RemoveGuildBan(r.guildID, userID)
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
