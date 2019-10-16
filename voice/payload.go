package voice

import (
	"encoding/json"

	"github.com/skwair/harmony/internal/payload"
)

// sendPayload sends a single Payload to the Voice server with
// the given op and data.
func (vc *Connection) sendPayload(op int, d interface{}) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	return payload.Send(&vc.connWMu, vc.conn, &payload.Payload{Op: op, D: b})
}

// recvPayload receives a single Payload from the Voice server.
func (vc *Connection) recvPayload() (*payload.Payload, error) {
	return payload.Recv(&vc.connRMu, vc.conn)
}
