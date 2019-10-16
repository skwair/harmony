package harmony

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/skwair/harmony/internal/endpoint"
)

// VoiceState represents the voice state of a user.
type VoiceState struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Deaf      bool   `json:"deaf"`
	Mute      bool   `json:"mute"`
	SelfDeaf  bool   `json:"self_deaf"`
	SelfMute  bool   `json:"self_mute"`
	Suppress  bool   `json:"suppress"` // Whether this user is muted by the current user.
}

// VoiceRegion represents a voice region a guild can use or is using for its voice channels.
type VoiceRegion struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	// Whether this is a vip-only server.
	VIP bool `json:"vip,omitempty"`
	// Whether this is a single server that is closest to the current user's client.
	Optimal bool `json:"optimal,omitempty"`
	// Whether this is a deprecated voice region (avoid switching to these.
	Deprecated bool `json:"deprecated,omitempty"`
	// Whether this is a custom voice region (used for events/etc).
	Custom bool `json:"custom,omitempty"`
}

// VoiceRegions returns a list of available voice regions that can be used when creating
// or updating servers.
func (c *Client) VoiceRegions(ctx context.Context, guildID string) ([]VoiceRegion, error) {
	e := endpoint.GetVoiceRegions()
	resp, err := c.doReq(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var regions []VoiceRegion
	if err = json.NewDecoder(resp.Body).Decode(&regions); err != nil {
		return nil, err
	}
	return regions, nil
}
