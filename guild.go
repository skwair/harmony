package harmony

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/skwair/harmony/discord"
	"github.com/skwair/harmony/internal/endpoint"
	"github.com/skwair/harmony/internal/rest"
)

// CreateGuild creates a new guild with the given name.
// Returns the created guild on success. Fires a Guild Create Gateway event.
func (c *Client) CreateGuild(ctx context.Context, name string) (*discord.Guild, error) {
	s := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	e := endpoint.CreateGuild()
	resp, err := c.restClient.Do(ctx, e, rest.JSONPayload(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, discord.NewAPIError(resp)
	}

	var g discord.Guild
	if err = json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return nil, err
	}
	return &g, nil
}

// RequestGuildMembers is used to request offline members for the guild. When initially
// connecting, the gateway will only send offline members if a guild has less than
// the large_threshold members (value in the Gateway Identify). If a client wishes
// to receive additional members, they need to explicitly request them via this
// operation. The server will send Guild Members Chunk events in response with up
// to 1000 members per chunk until all members that match the request have been sent.
// query is a string that username starts with, or an empty string to return all members.
// limit is the maximum number of members to send or 0 to request all members matched.
// You need to be connected to the Gateway to call this method, else it will
// return ErrGatewayNotConnected.
func (c *Client) RequestGuildMembers(guildID, query string, limit int) error {
	if !c.isConnected() {
		return discord.ErrGatewayNotConnected
	}

	req := struct {
		GuildID string `json:"guild_id"`
		Query   string `json:"query"`
		Limit   int    `json:"limit"`
	}{
		GuildID: guildID,
		Query:   query,
		Limit:   limit,
	}
	return c.sendPayload(c.ctx, gatewayOpcodeRequestGuildMembers, &req)
}
