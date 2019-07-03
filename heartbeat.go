package harmony

import (
	"encoding/binary"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
	errStaleConnection    = errors.New("stale connection")
	errStaleUDPConnection = errors.New("stale UDP connection")
)

// heartbeat periodically sends a heartbeat payload to the Gateway.
func (c *Client) heartbeat(every time.Duration) {
	c.logger.Debug("starting gateway heartbeater")
	defer c.logger.Debug("stopped gateway heartbeater")

	heartbeat(&c.wg, c.stop, c.error, every, c.sendHeartbeatPayload, &c.lastHeartbeatACK)
}

// sendHeartbeatPayload sends a single heartbeat payload
// to the Gateway containing the sequence number.
func (c *Client) sendHeartbeatPayload() error {
	var sequence *int64 // nil or seq if seq > 0
	if seq := atomic.LoadInt64(&c.sequence); seq != 0 {
		sequence = &seq
	}
	atomic.StoreInt64(&c.lastHeartbeatSend, time.Now().UnixNano())
	return c.sendPayload(gatewayOpcodeHeartbeat, sequence)
}

// heartbeat periodically sends a heartbeat payload to the voice server.
func (vc *VoiceConnection) heartbeat(every time.Duration) {
	vc.logger.Debug("starting voice connection heartbeater")
	defer vc.logger.Debug("stopped voice connection heartbeater")

	heartbeat(&vc.wg, vc.stop, vc.error, every, vc.sendHeartbeatPayload, &vc.lastHeartbeatACK)
}

// sendHeartbeatPayload sends a single heartbeat payload
// to the voice server containing a nonce.
func (vc *VoiceConnection) sendHeartbeatPayload() error {
	return vc.sendPayload(voiceOpcodeHeartbeat, time.Now().Unix())
}

// heartbeat periodically calls the heartbeater function to send a heartbeat payload.
// It should be called in a separate goroutine. It will decrement the given
// wait group when done, can be stopped by closing the stop channel and will
// report any error that occurs with the errCh channel.
func heartbeat(
	wg *sync.WaitGroup,
	stop chan struct{},
	errCh chan<- error,
	every time.Duration,
	heartbeater func() error,
	lastHeartbeatACK *int64,
) {
	defer wg.Done()

	ticker := time.NewTicker(every)
	defer ticker.Stop()

	first := true
	for {
		// If we haven't received a heartbeat ACK since the
		// last heartbeat we sent, we should consider the
		// connection as stale and return an error.
		t := atomic.LoadInt64(lastHeartbeatACK)
		if !first && time.Now().UTC().Sub(time.Unix(0, t).UTC()) > every {
			errCh <- errStaleConnection
			return
		}

		if err := heartbeater(); err != nil {
			errCh <- err
			return
		}

		if first {
			first = false
		}

		select {
		case <-stop:
			return
		case <-ticker.C:
		}
	}
}

// udpHeartbeat sends a UDP packet with a sequence number at a defined
// internal (every), incremented by one each time it heartbeats.
// It should be called in a separate goroutine. It will decrement the given
// wait group while running, can be stopped by closing the stop channel and
// will report any error that occurs with the errCh channel.
func (vc *VoiceConnection) udpHeartbeat(every time.Duration) {
	defer vc.wg.Done()

	vc.logger.Debug("starting UDP heartbeater")
	defer vc.logger.Debug("stopped UDP heartbeater")

	ticker := time.NewTicker(every)
	defer ticker.Stop()

	packet := make([]byte, 8)

	first := true
	for {
		// If we haven't received a heartbeat ACK since the
		// last heartbeat we sent, we should consider the
		// connection as stale and return an error.
		t := atomic.LoadInt64(&vc.lastUDPHeartbeatACK)
		if !first && time.Now().UTC().Sub(time.Unix(0, t).UTC()) > every {
			vc.error <- errStaleUDPConnection
			return
		}

		// Load and increment the UDP sequence atomically,
		// but send the value before the increment.
		binary.LittleEndian.PutUint64(packet, atomic.AddUint64(&vc.udpHeartbeatSequence, 1)-1)
		if _, err := vc.udpConn.Write(packet); err != nil {
			// Silently break out of this loop because
			// the connection was closed by the client.
			if isConnectionClosed(err) {
				return
			}

			vc.error <- err
			return
		}

		if first {
			first = false
		}

		select {
		case <-vc.stop:
			return
		case <-ticker.C:
		}
	}
}
