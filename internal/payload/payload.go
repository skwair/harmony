package payload

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// Payload is the content of a Discord Gateway or Voice event.
type Payload struct {
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
func (p *Payload) String() string {
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

// Send sends the given Payload on the given connection.
func Send(ctx context.Context, conn *websocket.Conn, p *Payload) error {
	return wsjson.Write(ctx, conn, p)
}

// Recv receives a single message from the provided connection, ensuring
// no concurrent call to conn.ReadMessage can occur.
// It also takes care of optionally decompressing the message and decoding
// it into a payload.
func Recv(ctx context.Context, connRMu *sync.Mutex, conn *websocket.Conn) (*Payload, error) {
	connRMu.Lock()
	typ, b, err := conn.Read(ctx)
	connRMu.Unlock()
	if err != nil {
		return nil, err
	}

	var rc io.ReadCloser
	br := bytes.NewReader(b)
	rc = ioutil.NopCloser(br)
	// If the payload is compressed, we first need to decompress it.
	if typ == websocket.MessageBinary {
		rc, err = zlib.NewReader(rc)
		if err != nil {
			return nil, err
		}
	}

	var p Payload
	if err = json.NewDecoder(rc).Decode(&p); err != nil {
		return nil, err
	}

	if err = rc.Close(); err != nil {
		return nil, err
	}

	return &p, nil
}
