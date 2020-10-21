package voice

import (
	"encoding/binary"
	"time"

	"github.com/skwair/harmony/internal/heartbeat"
)

// heartbeat periodically sends a heartbeat payload to the voice server.
func (vc *Connection) heartbeat(every time.Duration) {
	defer vc.wg.Done()

	vc.logger.Debug("starting voice connection heartbeater")
	defer vc.logger.Debug("stopped voice connection heartbeater")

	heartbeat.Run(
		every,
		vc.sendHeartbeatPayload,
		vc.lastHeartbeatAck,
		vc.lastHeartbeatSent,
		vc.stop,
		vc.reportErr,
	)
}

// sendHeartbeatPayload sends a single heartbeat payload
// to the voice server containing a nonce.
func (vc *Connection) sendHeartbeatPayload() error {
	return vc.sendPayload(vc.ctx, voiceOpcodeHeartbeat, time.Now().Unix())
}

// udpHeartbeat periodically sends a UDP heartbeat packet to the voice server.
func (vc *Connection) udpHeartbeat(every time.Duration) {
	defer vc.wg.Done()

	vc.logger.Debug("starting UDP heartbeater")
	defer vc.logger.Debug("stopped UDP heartbeater")

	heartbeat.RunUDP(
		every,
		vc.sendUDPHeartbeat,
		vc.lastUDPHeartbeatAck,
		vc.lastUDPHeartbeatSent,
		vc.stop,
		vc.reportErr,
	)
}

// sendUDPHeartbeat sends a single UDP heartbeat packet and increments the sequence number.
func (vc *Connection) sendUDPHeartbeat() error {
	packet := make([]byte, 8)

	// Load and increment the UDP sequence atomically,
	// but send the value before the increment.
	binary.LittleEndian.PutUint64(packet, vc.udpHeartbeatSequence.Add(1)-1)
	if _, err := vc.udpConn.Write(packet); err != nil {
		return err
	}

	return nil
}
