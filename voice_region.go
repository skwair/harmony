package harmony

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/skwair/harmony/discord"
	"github.com/skwair/harmony/internal/endpoint"
)

// GetVoiceRegions returns a list of available voice regions that can be used when creating
// or updating servers.
func (c *Client) GetVoiceRegions(ctx context.Context) ([]discord.VoiceRegion, error) {
	e := endpoint.GetVoiceRegions()
	resp, err := c.restClient.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var regions []discord.VoiceRegion
	if err = json.NewDecoder(resp.Body).Decode(&regions); err != nil {
		return nil, err
	}
	return regions, nil
}
