package harmony

import (
	"encoding/binary"
	"sync/atomic"
	"time"

	"github.com/skwair/harmony/internal/heartbeat"
)

// heartbeat periodically sends a heartbeat payload to the Gateway.
func (c *Client) heartbeat(every time.Duration) {
	c.logger.Debug("starting gateway heartbeater")
	defer c.logger.Debug("stopped gateway heartbeater")

	heartbeat.Run(&c.wg, c.stop, c.error, every, c.sendHeartbeatPayload, &c.lastHeartbeatACK)
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

	heartbeat.Run(&vc.wg, vc.stop, vc.error, every, vc.sendHeartbeatPayload, &vc.lastHeartbeatACK)
}

// sendHeartbeatPayload sends a single heartbeat payload
// to the voice server containing a nonce.
func (vc *VoiceConnection) sendHeartbeatPayload() error {
	return vc.sendPayload(voiceOpcodeHeartbeat, time.Now().Unix())
}

// sendUDPHeartbeat sends a single UDP heartbeat packet and increments the sequence number.
func (vc *VoiceConnection) sendUDPHeartbeat() error {
	packet := make([]byte, 8)

	// Load and increment the UDP sequence atomically,
	// but send the value before the increment.
	binary.LittleEndian.PutUint64(packet, atomic.AddUint64(&vc.udpHeartbeatSequence, 1)-1)
	if _, err := vc.udpConn.Write(packet); err != nil {
		return err
	}

	return nil
}
