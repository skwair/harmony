package payload

import (
	"nhooyr.io/websocket"
)

// ReceiverFunc is a function that receives a single Payload.
type ReceiverFunc func() (*Payload, error)

// HandlerFunc is a function that handles a Payload.
type HandlerFunc func(p *Payload) error

// ListenAndHandle loops on payloads received using the given receiver and
// handles them with the given handler as they arrive.
// If an error happens during the process, ListenAndHandle will report it
// using the given errReporter and return.
// It silently stops when the provided receiver's underlying connection is
// closed.
func ListenAndHandle(r ReceiverFunc, h HandlerFunc, errReporter func(err error)) {
	for {
		p, err := r()
		if err != nil {
			// Silently break out of this loop because the connection
			// was closed (either intentionally by calling Disconnect
			// or because we encountered an error).
			// NOTE: this is probably useless now that only the first error
			// is reported. It will be discarded if an error has already been
			// reported.
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusInternalError {
				return
			}

			errReporter(err)
			return
		}

		if err = h(p); err != nil {
			errReporter(err)
			return
		}
	}
}
