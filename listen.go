package harmony

import (
	"github.com/skwair/harmony/internal/payload"
)

// listen listens for payloads sent by the Discord Gateway
// and handles them as they are received.
func (c *Client) listen() {
	c.logger.Debug("starting gateway event listener")
	defer c.logger.Debug("stopped gateway event listener")

	payload.Listen(
		&c.wg,
		c.stop,
		c.reportErr,
		c.recvPayloads,
		c.handleEvent,
	)
}

func (c *Client) recvPayloads(ch chan<- *payload.Payload) {
	payload.RecvAll(&c.wg, ch, c.reportErr, c.stop, c.recvPayload)
}
