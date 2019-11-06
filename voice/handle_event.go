package voice

import (
	"time"

	"github.com/skwair/harmony/internal/payload"
)

func (vc *Connection) handleEvent(p *payload.Payload) error {
	switch p.Op {
	case voiceOpcodeReady, voiceOpcodeSessionDescription, voiceOpcodeHello:
		// Those events should be sent through the payloads channel if this
		// voice connection is currently being established so EstablishNewConnection
		// can receive them.
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
		// Not sure what to do with this event as it contains no additional info
		// and the main Gateway connection will receive a Voice State Update.
	}

	return nil
}
