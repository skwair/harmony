package harmony

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

// payload is the content of a Discord Gateway or Voice event.
type payload struct {
	// Opcode for the payload.
	Op int `json:"op"`
	// Event data.
	D json.RawMessage `json:"d"`
	// Sequence number, used for resuming sessions
	// and heartbeats. Only for Opcode 0.
	S int64 `json:"s,omitempty"`
	// The event name for this payload.
	// Only for Opcode 0.
	T string `json:"t,omitempty"`
}

// String implements fmt.Stringer.
func (p *payload) String() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("{code: %d", p.Op))

	if p.S != 0 {
		s.WriteString(fmt.Sprintf(", sequence: %d", p.S))
	}

	if p.T != "" {
		s.WriteString(fmt.Sprintf(", type: %s", p.T))
	}

	if len(p.D) > 0 && !bytes.Equal(p.D, []byte("null")) {
		s.WriteString(fmt.Sprintf(", data: %s", string(p.D)))
	}

	s.WriteRune('}')

	return s.String()
}

// sendPayload sends a single Payload to the Gateway with
// the given op and data.
func (c *Client) sendPayload(op int, d interface{}) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	p := &payload{Op: op, D: b}
	c.logger.Debugf("sent payload: %s", p)
	return sendPayload(&c.connWMu, c.conn, p)
}

// recvPayload receives a single Payload from the Gateway.
func (c *Client) recvPayload() (*payload, error) {
	p, err := recvPayload(&c.connRMu, c.conn)
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
	return sendPayload(&vc.connWMu, vc.conn, &payload{Op: op, D: b})
}

// recvPayload receives a single Payload from the Voice server.
func (vc *VoiceConnection) recvPayload() (*payload, error) {
	return recvPayload(&vc.connRMu, vc.conn)
}

// sendPayload ensures no concurrent call to conn.WriteJSON can occur.
func sendPayload(connWMu *sync.Mutex, conn *websocket.Conn, p *payload) error {
	connWMu.Lock()
	err := conn.WriteJSON(p)
	connWMu.Unlock()
	return err
}

// recvPayload receives a single message from the provided connection, ensuring
// no concurrent call to conn.ReadMessage can occur.
// It also takes care of optionally decompressing the message and decoding
// it into a payload.
func recvPayload(connRMu *sync.Mutex, conn *websocket.Conn) (*payload, error) {
	connRMu.Lock()
	typ, b, err := conn.ReadMessage()
	connRMu.Unlock()
	if err != nil {
		return nil, err
	}

	var rc io.ReadCloser
	br := bytes.NewReader(b)
	rc = ioutil.NopCloser(br)
	// If the payload is compressed, we first need to decompress it.
	if typ == websocket.BinaryMessage {
		rc, err = zlib.NewReader(rc)
		if err != nil {
			return nil, err
		}
	}

	var p payload
	if err = json.NewDecoder(rc).Decode(&p); err != nil {
		return nil, err
	}
	rc.Close()
	return &p, nil
}
