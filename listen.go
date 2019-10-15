package harmony

import (
	"github.com/skwair/harmony/internal/payload"
)

// listen listens for payloads sent by the Discord Gateway.
func (c *Client) listen() {
	c.logger.Debug("starting gateway event listener")
	defer c.logger.Debug("stopped gateway event listener")

	payload.Listen(
		&c.wg,
		c.stop,
		c.error,
		c.recvPayloads,
		c.handleEvent,
	)
}

// listen listens for payloads sent by the voice server.
func (vc *VoiceConnection) listen() {
	vc.logger.Debug("starting voice connection event listener")
	defer vc.logger.Debug("stopped voice connection event listener")

	payload.Listen(
		&vc.wg,
		vc.stop,
		vc.error,
		vc.recvPayloads,
		vc.handleEvent,
	)
}

func (c *Client) recvPayloads(ch chan<- *payload.Payload) {
	payload.RecvAll(&c.wg, ch, c.error, c.recvPayload)
}

func (vc *VoiceConnection) recvPayloads(ch chan<- *payload.Payload) {
	payload.RecvAll(&vc.wg, ch, vc.error, vc.recvPayload)
}
