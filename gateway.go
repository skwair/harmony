package harmony

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/skwair/harmony/internal/endpoint"
)

// Gateway returns a valid WSS URL, which the client can use for connecting.
func (c *Client) Gateway(ctx context.Context) (string, error) {
	e := endpoint.Gateway()
	resp, err := c.doReq(ctx, http.MethodGet, e, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var gateway struct {
		URL string
	}
	if err = json.NewDecoder(resp.Body).Decode(&gateway); err != nil {
		return "", err
	}
	c.gatewayURL = gateway.URL
	return gateway.URL, nil
}

// GatewayBot returns a valid WSS URL and the recommended number of shards to connect with.
func (c *Client) GatewayBot(ctx context.Context) (string, int, error) {
	e := endpoint.GatewayBot()
	resp, err := c.doReq(ctx, http.MethodGet, e, nil)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	var gateway struct {
		URL    string
		Shards int
	}
	if err = json.NewDecoder(resp.Body).Decode(&gateway); err != nil {
		return "", 0, err
	}
	c.gatewayURL = gateway.URL
	return gateway.URL, gateway.Shards, nil
}
