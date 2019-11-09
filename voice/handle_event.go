package voice

import (
	"time"

	"github.com/skwair/harmony/internal/payload"
)

func (vc *Connection) handleEvent(p *payload.Payload) error {
	switch p.Op {
	case voiceOpcodeReady, voiceOpcodeSessionDescription, voiceOpcodeHello:
		// Those events should be sent through the payloads channel if this
		// voice connection is currently being established so Connect can
		// receive them.
		if vc.isConnecting() {
			vc.payloads <- p
		}

	// Heartbeat ACK.
	case voiceOpcodeHeartbeatACK:
		// TODO: Check nonce ?
		vc.lastHeartbeatACK.Store(time.Now().UnixNano())

	// Resume acknowledged by the voice server.
	case voiceOpcodeResumed:
		if vc.isConnecting() {
			vc.payloads <- p
		}

	// A client has disconnected from the voice channel.
	case voiceOpcodeClientDisconnect:
		// TODO: add a way to register to those events.
		// Example payload: {code: 13, data: {"user_id":"220152355228164927"}}
	}

	return nil
}
