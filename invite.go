package harmony

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/skwair/harmony/internal/endpoint"
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

// InviteResource is a resource that allows to perform various actions on a Discord invite.
// Create one with Client.Invite.
type InviteResource struct {
	code   string
	client *Client
}

// Invite returns a new invite resource to manage the invite with the given code.
func (c *Client) Invite(code string) *InviteResource {
	return &InviteResource{code: code, client: c}
}

// Get returns the invite. If withCounts is set to true,
// the returned invite will contain the approximate member counts.
func (r *InviteResource) Get(ctx context.Context, withCounts bool) (*Invite, error) {
	q := url.Values{}
	q.Set("with_counts", strconv.FormatBool(withCounts))

	e := endpoint.GetInvite(r.code, q.Encode())
	resp, err := r.client.doReq(ctx, e, nil)
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

// Delete deletes the invite. Requires the MANAGE_CHANNELS permission.
// Returns the deleted invite on success.
func (r *InviteResource) Delete(ctx context.Context) (*Invite, error) {
	e := endpoint.DeleteInvite(r.code)
	resp, err := r.client.doReq(ctx, e, nil)
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
