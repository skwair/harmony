package voice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"nhooyr.io/websocket"

	"github.com/skwair/harmony/internal/payload"
)

var errInvalidSession = errors.New("invalid voice session")

func (vc *Connection) reconnectWithBackoff() {
	vc.reconnecting.Store(true)
	defer vc.reconnecting.Store(false)

	vc.logger.Debug("trying to reconnect to the voice server")

	for i := 0; true; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		if err := vc.reconnect(ctx); err != nil {
			cancel()

			if err == errInvalidSession {
				vc.logger.Error("invalid session, can not recover")
				return
			}

			duration := time.Second * 2
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
	vc.payloads = make(chan *payload.Payload)
	vc.error = make(chan error)
	vc.stop = make(chan struct{})

	vc.ctx, vc.cancel = context.WithCancel(context.Background())

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

	// This is used to notify the event handler that some
	// specific payloads should be sent through to vc.payloads
	// while we are reconnecting to the voice server.
	vc.connectingToVoice.Store(true)
	defer vc.connectingToVoice.Store(false)

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
		if websocket.CloseStatus(err) == 4006 {
			return errInvalidSession
		}
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
		ServerID:  vc.guildID,
		SessionID: vc.sessionID,
		Token:     vc.token,
	}
	if err = vc.sendPayload(ctx, voiceOpcodeResume, r); err != nil {
		return err
	}

	// We should receive an Opcode 9 Resumed payload to acknowledge the resume.
	p, err = vc.recvPayload()
	if err != nil {
		if websocket.CloseStatus(err) == 4006 {
			return errInvalidSession
		}
	}
	if p.Op != voiceOpcodeResumed {
		return fmt.Errorf("expected Opcode 9 Resumed; got Opcode %d", p.Op)
	}

	vc.wg.Add(2) // listen starts an additional goroutine.
	go vc.listen()

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
