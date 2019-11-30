package voice

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"nhooyr.io/websocket"
)

// Determine whether we should try to reconnect based on the error we got.
// See https://discordapp.com/developers/docs/topics/opcodes-and-status-codes#voice-voice-close-event-codes for more information.
func shouldReconnect(err error) bool {
	if err == nil {
		return false
	}

	switch websocket.CloseStatus(err) {
	case 4003, 4004, 4005, 4006, 4011, 4012, 4014, 4016:
		return false
	case 4015:
		return true
	case -1: // Not a websocket error.
		return true
	default: // New (or undocumented?) close status code.
		return true
	}
}

func (vc *Connection) reconnectWithBackoff() {
	vc.reconnecting.Store(true)
	defer vc.reconnecting.Store(false)

	vc.logger.Debug("trying to reconnect to the voice server")

	for i := 0; true; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		if err := vc.reconnect(ctx); err != nil {
			cancel()

			if !shouldReconnect(err) {
				vc.logger.Error("invalid voice session, can not recover: %v", err)
				return
			}

			duration := 2 * time.Second
			vc.logger.Errorf("failed to reconnect to voice server: %v, retrying in %s", err, duration)

			select {
			case <-time.After(duration):
				continue // Make a new connection attempt.
			case <-vc.stop:
				// Client called Disconnect(), stop trying to reconnect.
				vc.logger.Debug("client called Close while trying to reconnect to the voice server, aborting")
				return
			}
		} else {
			// We could reconnect.
			vc.logger.Info("successfully reconnected to the voice server")
			cancel()
			return
		}
	}
}

func (vc *Connection) reconnect(ctx context.Context) error {
	// This is used to notify the event handler that some
	// specific payloads should be sent through to vc.payloads
	// while we are reconnecting to the voice server.
	vc.connecting.Store(true)
	defer vc.connecting.Store(false)

	vc.reset()

	// Start by re-opening the voice websocket connection.
	var err error
	vc.logger.Debugf("connecting to voice server: %s", vc.endpoint)
	vc.conn, _, err = websocket.Dial(ctx, vc.endpoint, nil)
	if err != nil {
		return err
	}
	// From now on, if any error occurs during the rest of the
	// voice reconnection process, we should close the underlying
	// websocket so we can try to reconnect.
	defer func() {
		if err != nil {
			_ = vc.conn.Close(websocket.StatusInternalError, "failed to reestablish voice connection")
			vc.connected.Store(false)
			close(vc.stop)
			vc.cancel()
		}
	}()

	// Then re-establish the voice data UDP connection.
	vc.udpConn, err = net.DialUDP("udp", nil, vc.dataEndpoint)
	if err != nil {
		return err
	}
	// From now on, close the UDP connection if any error occurs.
	defer func() {
		if err != nil {
			_ = vc.udpConn.Close()
		}
	}()

	// Start heartbeating on the UDP connection.
	vc.wg.Add(1)
	go vc.udpHeartbeat(5 * time.Second)

	// Once the websocket connection is re opened, the sever should send us an
	// Opcode 8 Hello payload, indicating the heartbeat interval we should use.
	p, err := vc.recvPayload()
	if err != nil {
		return fmt.Errorf("could not receive Hello payload: %w", err)
	}
	if p.Op != voiceOpcodeHello {
		return fmt.Errorf("expected Opcode 8 Hello; got Opcode %d", p.Op)
	}

	var h struct {
		V                 int     `json:"v"`
		HeartbeatInterval float64 `json:"heartbeat_interval"`
	}
	if err = json.Unmarshal(p.D, &h); err != nil {
		return err
	}

	// Send the resume payload to notify the voice server this is not a new connection.
	r := resume{
		ServerID:  vc.State().GuildID,
		SessionID: vc.State().SessionID,
		Token:     vc.token,
	}
	if err = vc.sendPayload(ctx, voiceOpcodeResume, r); err != nil {
		return err
	}

	// We should receive an Opcode 9 Resumed payload to acknowledge the resume.
	p, err = vc.recvPayload()
	if err != nil {
		return fmt.Errorf("could not receive Resumed payload: %w", err)
	}
	if p.Op != voiceOpcodeResumed {
		return fmt.Errorf("expected Opcode 9 Resumed; got Opcode %d", p.Op)
	}

	vc.wg.Add(1)
	go vc.listenAndHandlePayloads()

	vc.wg.Add(1)
	go vc.wait()

	vc.wg.Add(1)
	go vc.heartbeat(time.Duration(h.HeartbeatInterval) * time.Millisecond)

	vc.wg.Add(3) // opusReceiver starts an additional goroutine.
	vc.opusReadinessWG.Add(2)
	go vc.opusReceiver()
	go vc.opusSender()

	// Making sure Opus receiver and sender are started.
	vc.opusReadinessWG.Wait()

	if err = vc.sendSilenceFrame(); err != nil {
		return err
	}

	vc.connected.Store(true)

	return nil
}
