package harmony

import (
	"encoding/json"
	"math/rand"
	"sync/atomic"
	"time"
)

func (c *Client) handleEvent(p *payload) error {
	switch p.Op {
	// Dispatch.
	case 0:
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
	case 1:
		if err := c.sendHeartbeatPayload(); err != nil {
			return err
		}

	// Reconnect.
	case 7:
		c.Disconnect()
		if err := c.Connect(); err != nil {
			return err
		}

	// Invalid Session.
	case 9:
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
			c.sessionID = ""
			atomic.StoreInt64(&c.sequence, 0)
			if err := c.identify(); err != nil {
				return err
			}
		}

	// Hello.
	case 10:
		// Handled by Connect()

	// Heartbeat ACK.
	case 11:
		if c.withStateTracking {
			c.State.setRTT(time.Since(time.Unix(0, c.lastHeartbeatSend)))
		}
		atomic.StoreInt64(&c.lastHeartbeatACK, time.Now().UnixNano())
	}
	return nil
}

func (vc *VoiceConnection) handleEvent(p *payload) error {
	switch p.Op {
	// Ready.
	case 2:
		// Those two events should be sent through the payloads channel if this
		// voice connection is currently connecting to a voice channel so the
		// ConnectToVoice method can receive them.
		if atomic.LoadInt32(&vc.connectingToVoice) == 1 {
			vc.payloads <- p
		}

	// Session description.
	case 4:
		// Those two events should be sent through the payloads channel if this
		// voice connection is currently connecting to a voice channel so the
		// ConnectToVoice method can receive them.
		if atomic.LoadInt32(&vc.connectingToVoice) == 1 {
			vc.payloads <- p
		}

	// Heartbeat ACK.
	case 6:
		// TODO: Check nonce ?
		atomic.StoreInt64(&vc.lastHeartbeatACK, time.Now().UnixNano())

	// Hello.
	case 8:
		// Those two events should be sent through the payloads channel if this
		// voice connection is currently connecting to a voice channel so the
		// ConnectToVoice method can receive them.
		if atomic.LoadInt32(&vc.connectingToVoice) == 1 {
			vc.payloads <- p
		}

	// Resumed.
	case 9:

	// Client disconnect.
	case 13:
	}
	return nil
}
