package voice

import (
	"encoding/binary"
	"sync/atomic"
	"time"

	"github.com/skwair/harmony/internal/heartbeat"
)

// heartbeat periodically sends a heartbeat payload to the voice server.
func (vc *Connection) heartbeat(every time.Duration) {
	vc.logger.Debug("starting voice connection heartbeater")
	defer vc.logger.Debug("stopped voice connection heartbeater")

	heartbeat.Run(
		&vc.wg,
		vc.stop,
		vc.error,
		every,
		vc.sendHeartbeatPayload,
		&vc.lastHeartbeatACK,
	)
}

// sendHeartbeatPayload sends a single heartbeat payload
// to the voice server containing a nonce.
func (vc *Connection) sendHeartbeatPayload() error {
	return vc.sendPayload(vc.ctx, voiceOpcodeHeartbeat, time.Now().Unix())
}

// udpHeartbeat periodically sends a UDP heartbeat packet to the voice server.
func (vc *Connection) udpHeartbeat(every time.Duration) {
	vc.logger.Debug("starting UDP heartbeater")
	defer vc.logger.Debug("stopped UDP heartbeater")

	heartbeat.RunUDP(
		&vc.wg,
		vc.stop,
		vc.error,
		time.Second*5,
		vc.sendUDPHeartbeat,
		&vc.lastUDPHeartbeatACK,
	)
}

// sendUDPHeartbeat sends a single UDP heartbeat packet and increments the sequence number.
func (vc *Connection) sendUDPHeartbeat() error {
	packet := make([]byte, 8)

	// Load and increment the UDP sequence atomically,
	// but send the value before the increment.
	binary.LittleEndian.PutUint64(packet, atomic.AddUint64(&vc.udpHeartbeatSequence, 1)-1)
	if _, err := vc.udpConn.Write(packet); err != nil {
		return err
	}

	return nil
}
