package harmony

import (
	"context"
	"encoding/json"
	"math/rand"
	"sync/atomic"
	"time"
)

// handleEvent handles all events received from Discord's Gateway once connected to it.
func (c *Client) handleEvent(p *payload) error {
	switch p.Op {
	case gatewayOpcodeDispatch:
		atomic.StoreInt64(&c.sequence, p.S)
		// Those two events should be sent through the payloads channel if the
		// client is currently connecting to a voice channel so the ConnectToVoice
		// method can receive them.
		if (p.T == eventVoiceStateUpdate || p.T == eventVoiceServerUpdate) &&
			c.isConnectingToVoice() {
			c.voicePayloads <- p
		}
		if err := c.dispatch(p.T, p.D); err != nil {
			return err
		}

	// Heartbeat requested from the Gateway (used for ping checking).
	case gatewayOpcodeHeartbeat:
		if err := c.sendHeartbeatPayload(); err != nil {
			return err
		}

	// Gateway is asking us to reconnect.
	case gatewayOpcodeReconnect:
		c.Disconnect()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := c.Connect(ctx); err != nil {
			return err
		}

	// Gateway is telling us our session ID is invalid.
	case gatewayOpcodeInvalidSession:
		var resumable bool
		if err := json.Unmarshal(p.D, &resumable); err != nil {
			return err
		}

		if resumable {
			if err := c.resume(); err != nil {
				return err
			}
		} else {
			// If we could not resume a session in time, we will receive an
			// Invalid Session payload and are expected to wait a bit before
			// sending a fresh Identify payload.
			// https://discordapp.com/developers/docs/topics/gateway#resuming.
			time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)

			c.resetGatewaySession()
			if err := c.identify(); err != nil {
				return err
			}
		}

	case gatewayOpcodeHello:
		// Handled by Connect()

	case gatewayOpcodeHeartbeatACK:
		if c.withStateTracking {
			c.State.setRTT(time.Since(time.Unix(0, c.lastHeartbeatSend)))
		}
		atomic.StoreInt64(&c.lastHeartbeatACK, time.Now().UnixNano())
	}
	return nil
}

func (vc *VoiceConnection) handleEvent(p *payload) error {
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
