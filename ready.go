package discord

import (
	"encoding/json"
	"fmt"
)

// Ready is the Event fired by the Gateway after the client sent
// a valid Identify payload.
type Ready struct {
	V               int            `json:"v"` // Gateway version.
	User            *User          `json:"user"`
	PrivateChannels []Channel      `json:"private_channels"`
	Guilds          []PartialGuild `json:"guilds"`
	SessionID       string         `json:"session_id"`
	Trace           []string       `json:"_trace"`
}

// ready expects to receive a Ready payload from the Gateway and will set the
// session ID of the client if it receive it, else an error is returned.
func (c *Client) ready() error {
	p, err := c.recvPayload()
	if err != nil {
		return fmt.Errorf("could not receive ready payload from gateway: %v", err)
	}
	if p.Op != 0 || p.T != eventReady {
		return fmt.Errorf("expected Opcode 0 Ready; got Opcode %d %s", p.Op, p.T)
	}

	var rdy Ready
	if err := json.Unmarshal(p.D, &rdy); err != nil {
		return err
	}
	c.sessionID = rdy.SessionID
	c.userID = rdy.User.ID

	if c.withStateTracking {
		c.State.setInitialState(&rdy)
	}

	// Let this event be dispatched so the user
	// can get the initial state of the connection.
	return c.handleEvent(p)
}
