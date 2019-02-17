package harmony

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/skwair/harmony/internal/endpoint"
)

// Ban represents a Guild ban.
type Ban struct {
	Reason string
	User   *User
}

// Bans returns a list of bans for the users banned from this guild.
// Requires the 'BAN_MEMBERS' permission.
func (r *GuildResource) Bans(ctx context.Context) ([]Ban, error) {
	e := endpoint.GetGuildBans(r.guildID)
	resp, err := r.client.doReq(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var bans []Ban
	if err = json.NewDecoder(resp.Body).Decode(&bans); err != nil {
		return nil, err
	}
	return bans, nil
}

// Ban is a shorthand to ban a user with no reason and without
// deleting his messages. Requires the 'BAN_MEMBERS' permission.
// For more control, use the BanWithReason method.
func (r *GuildResource) Ban(ctx context.Context, userID string) error {
	return r.BanWithReason(ctx, userID, 0, "")
}

// BanWithReason creates a guild ban, and optionally delete previous messages
// sent by the banned user. Requires the 'BAN_MEMBERS' permission.
// Parameter delMsgDays is the number of days to delete messages for (0-7).
// Fires a Guild Ban Add Gateway event.
func (r *GuildResource) BanWithReason(ctx context.Context, userID string, delMsgDays int, reason string) error {
	q := url.Values{}
	if reason != "" {
		q.Set("reason", reason)
	}
	if delMsgDays > 0 {
		q.Set("delete-message-days", strconv.Itoa(delMsgDays))
	}

	e := endpoint.CreateGuildBan(r.guildID, userID, q.Encode())
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

// Unban is like UnbanWithReason but with no particular reason.
func (r *GuildResource) Unban(ctx context.Context, userID string) error {
	return r.UnbanWithReason(ctx, userID, "")
}

// Unban removes the ban for a user. Requires the 'BAN_MEMBERS' permissions.
// Fires a Guild Ban Remove Gateway event.
// The given reason will be set in the audit log entry for this action.
func (r *GuildResource) UnbanWithReason(ctx context.Context, userID, reason string) error {
	e := endpoint.RemoveGuildBan(r.guildID, userID)
	resp, err := r.client.doReqWithHeader(ctx, e, nil, reasonHeader(reason))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}
