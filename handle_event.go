package harmony

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/skwair/harmony/internal/payload"
)

// handleEvent handles all events received from Discord's Gateway once connected to it.
func (c *Client) handleEvent(p *payload.Payload) error {
	switch p.Op {
	case gatewayOpcodeDispatch:
		c.sequence.Store(p.S)

		// Those two events should be sent through the voice payloads channel if the
		// client is currently connecting to a voice channel so the JoinVoiceChannel
		// method can receive them.
		if (p.T == eventVoiceStateUpdate || p.T == eventVoiceServerUpdate) &&
			c.isConnectingToVoice() {
			c.voicePayloads <- p
		}

		if err := c.dispatch(p.T, p.D); err != nil {
			return fmt.Errorf("dispatch: %w", err)
		}

	// Heartbeat requested from the Gateway (used for ping checking).
	case gatewayOpcodeHeartbeat:
		if err := c.sendHeartbeatPayload(); err != nil {
			return fmt.Errorf("send heartbeat payload: %w", err)
		}

	// Gateway is asking us to reconnect.
	case gatewayOpcodeReconnect:
		return errMustReconnect

	// Gateway is telling us our session ID is invalid.
	case gatewayOpcodeInvalidSession:
		var resumable bool
		if err := json.Unmarshal(p.D, &resumable); err != nil {
			return fmt.Errorf("unmarshal resume: %w", err)
		}

		if resumable {
			if err := c.resume(c.ctx); err != nil {
				return fmt.Errorf("resume: %w", err)
			}
		} else {
			// If we could not resume a session in time, we will receive an
			// Invalid Session payload and are expected to wait a bit before
			// sending a fresh Identify payload.
			// https://discord.com/developers/docs/topics/gateway#resuming.
			time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)

			c.resetGatewaySession()
			if err := c.identify(c.ctx); err != nil {
				return fmt.Errorf("identify: %w", err)
			}
		}

	case gatewayOpcodeHello:
		// Handled by Connect()

	case gatewayOpcodeHeartbeatAck:
		if c.withStateTracking {
			c.State.setRTT(time.Since(time.Unix(0, c.lastHeartbeatSent.Load())))
		}
		c.lastHeartbeatAck.Store(time.Now().UnixNano())
	}
	return nil
}
