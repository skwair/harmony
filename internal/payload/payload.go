package payload

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

// Send sends the given Payload, ensuring no concurrent call to conn.WriteJSON can occur.
func Send(connWMu *sync.Mutex, conn *websocket.Conn, p *Payload) error {
	connWMu.Lock()
	err := conn.WriteJSON(p)
	connWMu.Unlock()
	return err
}

// Recv receives a single Payload from the provided connection, ensuring
// no concurrent call to conn.ReadMessage can occur.
// It also takes care of optionally decompressing the message and decoding
// it into a payload.
func Recv(connRMu *sync.Mutex, conn *websocket.Conn) (*Payload, error) {
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

	var p Payload
	if err = json.NewDecoder(rc).Decode(&p); err != nil {
		return nil, err
	}
	rc.Close()
	return &p, nil
}
