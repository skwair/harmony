package payload

import (
	"sync"

	"nhooyr.io/websocket"
)

// Listen uses the given receiver to receive payloads and passes them to the
// given handler as they arrive. It should be called in a separate goroutine.
// It will decrement the given wait group when done, can be stopped
// by closing the stop channel and will report any error that occurs with
// the errCh channel.
func Listen(
	wg *sync.WaitGroup,
	stop chan struct{},
	errCh chan<- error,
	receiver func(ch chan<- *Payload),
	handler func(p *Payload) error,
) {
	defer wg.Done()

	payloads := make(chan *Payload)
	go receiver(payloads)

	for {
		select {
		case p := <-payloads:
			if err := handler(p); err != nil {
				errCh <- err
				return
			}
		case <-stop:
			return
		}
	}
}

// RecvAll uses the receiver to receive payloads and send them
// through payloads as they arrive. It should be called in a separate
// goroutine. It will decrement the given wait group when done, can be
// stopped by closing the stop channel and will report any error that
// occurs with the errCh channel.
func RecvAll(
	wg *sync.WaitGroup,
	payloads chan<- *Payload,
	errCh chan<- error,
	receiver func() (*Payload, error),
) {
	defer wg.Done()

	for {
		p, err := receiver()
		if err != nil {
			// Silently break out of this loop because
			// the connection was closed by the client.
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				return
			}

			// NOTE: maybe treat websocket close errors differently based on their code.
			// See : https://discordapp.com/developers/docs/topics/opcodes-and-status-codes

			errCh <- err
			return
		}

		payloads <- p
	}
}
