package harmony

import (
	"time"

	"github.com/skwair/harmony/internal/heartbeat"
)

// heartbeat periodically sends a heartbeat payload to the Gateway.
func (c *Client) heartbeat(every time.Duration) {
	c.logger.Debug("starting gateway heartbeater")
	defer c.logger.Debug("stopped gateway heartbeater")

	heartbeat.Run(
		&c.wg,
		c.stop,
		c.reportErr,
		every,
		c.sendHeartbeatPayload,
		c.lastHeartbeatACK,
	)
}

// sendHeartbeatPayload sends a single heartbeat payload
// to the Gateway containing the sequence number.
func (c *Client) sendHeartbeatPayload() error {
	var sequence *int64 // nil or seq if seq > 0
	if seq := c.sequence.Load(); seq != 0 {
		sequence = &seq
	}
	c.lastHeartbeatSend.Store(time.Now().UnixNano())
	return c.sendPayload(c.ctx, gatewayOpcodeHeartbeat, sequence)
}
