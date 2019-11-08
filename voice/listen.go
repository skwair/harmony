package voice

import "github.com/skwair/harmony/internal/payload"

// listen listens for payloads sent by the voice server.
func (vc *Connection) listen() {
	vc.logger.Debug("starting voice connection event listener")
	defer vc.logger.Debug("stopped voice connection event listener")

	payload.Listen(
		&vc.wg,
		vc.stop,
		vc.reportErr,
		vc.recvPayloads,
		vc.handleEvent,
	)
}

func (vc *Connection) recvPayloads(ch chan<- *payload.Payload) {
	payload.RecvAll(&vc.wg, ch, vc.reportErr, vc.stop, vc.recvPayload)
}
