package voice

import (
	"sync/atomic"
	"time"

	"github.com/skwair/harmony/internal/payload"
)

func (vc *Connection) handleEvent(p *payload.Payload) error {
	switch p.Op {
	case voiceOpcodeReady, voiceOpcodeSessionDescription, voiceOpcodeHello:
		// Those events should be sent through the payloads channel if this
		// voice connection is currently connecting to a voice channel so the
		// connect method can receive them.
		if vc.isConnecting() {
			vc.payloads <- p
		}

	// Heartbeat ACK.
	case voiceOpcodeHeartbeatACK:
		// TODO: Check nonce ?
		atomic.StoreInt64(&vc.lastHeartbeatACK, time.Now().UnixNano())

	// Resume acknowledged by the voice server.
	case voiceOpcodeResumed:

	// A client has disconnected from the voice channel.
	case voiceOpcodeClientDisconnect:
	}

	return nil
}
