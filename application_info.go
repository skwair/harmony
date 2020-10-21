package harmony

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/skwair/harmony/discord"
	"github.com/skwair/harmony/internal/endpoint"
)

// GetApplicationInfo returns the current user's OAuth2 application info.
func (c *Client) GetApplicationInfo(ctx context.Context) (*discord.Application, error) {
	e := endpoint.GetApplication()
	resp, err := c.restClient.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var app discord.Application
	if err = json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, err
	}
	return &app, nil
}
