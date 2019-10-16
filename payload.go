package harmony

import (
	"encoding/json"

	"github.com/skwair/harmony/internal/payload"
)

// sendPayload sends a single Payload to the Gateway with
// the given op and data.
func (c *Client) sendPayload(op int, d interface{}) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	p := &payload.Payload{Op: op, D: b}
	c.logger.Debugf("sent payload: %s", p)
	return payload.Send(&c.connWMu, c.conn, p)
}

// recvPayload receives a single Payload from the Gateway.
func (c *Client) recvPayload() (*payload.Payload, error) {
	p, err := payload.Recv(&c.connRMu, c.conn)
	if err != nil {
		return nil, err
	}

	c.logger.Debugf("received payload: %s", p)

	return p, nil
}

// sendPayload sends a single Payload to the Voice server with
// the given op and data.
func (vc *VoiceConnection) sendPayload(op int, d interface{}) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	return payload.Send(&vc.connWMu, vc.conn, &payload.Payload{Op: op, D: b})
}

// recvPayload receives a single Payload from the Voice server.
func (vc *VoiceConnection) recvPayload() (*payload.Payload, error) {
	return payload.Recv(&vc.connRMu, vc.conn)
}
