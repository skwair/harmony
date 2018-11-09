package discord

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/skwair/discord/internal/endpoint"
)

// Invite represents a code that when used, adds a user to a guild or group DM channel.
type Invite struct {
	Code                     string   `json:"code,omitempty"`
	Guild                    *Guild   `json:"guild,omitempty"` // Nil if this invite is for a group DM channel.
	Channel                  *Channel `json:"channel,omitempty"`
	ApproximatePresenceCount int      `json:"approximate_presence_count,omitempty"`
	ApproximateMemberCount   int      `json:"approximate_member_count,omitempty"`

	InviteMetadata
}

// InviteMetadata contains additional information about an Invite.
type InviteMetadata struct {
	Inviter   *User     `json:"inviter,omitempty"`
	Uses      int       `json:"uses,omitempty"`
	MaxUses   int       `json:"max_uses,omitempty"`
	MaxAge    int       `json:"max_age,omitempty"`
	Temporary bool      `json:"temporary,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Revoked   bool      `json:"revoked,omitempty"`
}

// GetInvite returns the invite corresponding to the given code.
// If withCounts is set to true, the returned invite will contain
// the approximate member counts.
func (c *Client) GetInvite(code string, withCounts bool) (*Invite, error) {
	q := url.Values{}
	q.Set("with_counts", strconv.FormatBool(withCounts))

	e := endpoint.GetInvite(code, q.Encode())
	resp, err := c.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var invite Invite
	if err = json.NewDecoder(resp.Body).Decode(&invite); err != nil {
		return nil, err
	}
	return &invite, nil
}

// DeleteInvite deletes an invite given its code. Requires the MANAGE_CHANNELS permission.
// Returns the deleted invite on success.
func (c *Client) DeleteInvite(code string) (*Invite, error) {
	e := endpoint.DeleteInvite(code)
	resp, err := c.doReq(http.MethodDelete, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var invite Invite
	if err = json.NewDecoder(resp.Body).Decode(&invite); err != nil {
		return nil, err
	}
	return &invite, nil
}
