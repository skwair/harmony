package invite

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

// Get returns the invite. If withCounts is set to true,
// the returned invite will contain the approximate member counts.
func (r *Resource) Get(ctx context.Context, withCounts bool) (*discord.Invite, error) {
	q := url.Values{}
	q.Set("with_counts", strconv.FormatBool(withCounts))

	e := endpoint.GetInvite(r.code, q.Encode())
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var invite discord.Invite
	if err = json.NewDecoder(resp.Body).Decode(&invite); err != nil {
		return nil, err
	}
	return &invite, nil
}

// Delete is like DeleteWithReason but with no particular reason.
func (r *Resource) Delete(ctx context.Context) (*discord.Invite, error) {
	return r.DeleteWithReason(ctx, "")
}

// DeleteWithReason deletes the invite. Requires the MANAGE_CHANNELS permission.
// Returns the deleted invite on success.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) DeleteWithReason(ctx context.Context, reason string) (*discord.Invite, error) {
	e := endpoint.DeleteInvite(r.code)
	resp, err := r.client.DoWithHeader(ctx, e, nil, rest.ReasonHeader(reason))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var invite discord.Invite
	if err = json.NewDecoder(resp.Body).Decode(&invite); err != nil {
		return nil, err
	}
	return &invite, nil
}
