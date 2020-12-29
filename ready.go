package harmony

import (
	"encoding/json"
	"fmt"

	"github.com/skwair/harmony/discord"
)

// Ready is the Event fired by the Gateway after the client sent
// a valid Identify payload.
type Ready struct {
	V                    int                        `json:"v"` // Gateway version.
	User                 *discord.User              `json:"user"`
	Guilds               []discord.UnavailableGuild `json:"guilds"`
	SessionID            string                     `json:"session_id"`
	Application          discord.PartialApplication `json:"application"`
	GeoOrderedRTCRegions []string                   `json:"geo_ordered_rtc_regions"`
	Shard                [2]int                     `json:"shard"`
}

// recvReady expects to receive a Ready payload from the Gateway and will set the
// session ID of the client if it receive it, else an error is returned.
func (c *Client) recvReady() error {
	p, err := c.recvPayload()
	if err != nil {
		return fmt.Errorf("could not receive ready payload from gateway: %w", err)
	}
	if p.Op != gatewayOpcodeDispatch || p.T != eventReady {
		return fmt.Errorf("expected Opcode 0 Ready; got Opcode %d %s", p.Op, p.T)
	}

	var rdy Ready
	if err = json.Unmarshal(p.D, &rdy); err != nil {
		return err
	}
	c.sessionID = rdy.SessionID
	c.userID = rdy.User.ID

	if c.withStateTracking {
		c.logger.Debug("initializing state tracker")
		c.State.setInitialState(&rdy)
	}

	// Let this event be dispatched so the user
	// can get the initial state of the connection.
	return c.handleEvent(p)
}
