package harmony

import (
	"sync"
)

// listen listens for payloads sent by the Discord Gateway.
func (c *Client) listen() {
	listen(&c.wg, c.stop, c.error, c.recvPayloads, c.handleEvent)
}

// listen listens for payloads sent by the Voice server.
func (vc *VoiceConnection) listen() {
	listen(&vc.wg, vc.stop, vc.error, vc.recvPayloads, vc.handleEvent)
}

// listen uses the receiver to receive payloads and passes them to the
// handler as they arrive. It should be called in a separate goroutine.
// It will decrement the given wait group when done, can be stopped
// by closing the stop channel and will report any error that occurs with
// the errCh channel.
func listen(wg *sync.WaitGroup,
	stop chan struct{},
	errCh chan<- error,
	receiver func(ch chan<- *payload),
	handler func(p *payload) error) {
	defer wg.Done()

	ch := make(chan *payload)
	go receiver(ch)

	for {
		select {
		case p := <-ch:
			if err := handler(p); err != nil {
				errCh <- err
				return
			}
		case <-stop:
			return
		}
	}
}

func (c *Client) recvPayloads(ch chan<- *payload) {
	recvPayloads(&c.wg, ch, c.error, c.recvPayload)
}

func (vc *VoiceConnection) recvPayloads(ch chan<- *payload) {
	recvPayloads(&vc.wg, ch, vc.error, vc.recvPayload)
}

// recvPayloads uses the receiver to receive payloads and send them
// through pCh as they arrive. It should be called in a separate
// goroutine. It will decrement the given wait group when done, can be
// stopped by closing the stop channel and will report any error that
// occurs with the errCh channel.
func recvPayloads(wg *sync.WaitGroup,
	pCh chan<- *payload,
	errCh chan<- error,
	receiver func() (*payload, error)) {
	defer wg.Done()

	for {
		p, err := receiver()
		if err != nil {
			// Silently break out of this loop because
			// the connection was closed by the client.
			if isConnectionClosed(err) {
				return
			}

			// NOTE: maybe treat websocket close errors differently based on their code.
			// See : https://discordapp.com/developers/docs/topics/opcodes-and-status-codes
			// if e, ok := err.(*websocket.CloseError); ok {
			// }

			errCh <- err
			return
		}
		pCh <- p
	}
}
