package heartbeat

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	errStaleConnection    = errors.New("stale connection")
	errStaleUDPConnection = errors.New("stale UDP connection")
)

// Hearbeater is a function that sends a heartbeat.
type Hearbeater func() error

// Run periodically calls the given heartbeater to send a heartbeat payload.
// It should be called in a separate goroutine. It will decrement the given
// wait group when done, can be stopped by closing the stop channel and will
// report any error that occurs through the provided errCh.
func Run(
	wg *sync.WaitGroup,
	stop chan struct{},
	errCh chan<- error,
	every time.Duration,
	h Hearbeater,
	lastHeartbeatACK *int64,
) {
	defer wg.Done()

	ticker := time.NewTicker(every)
	defer ticker.Stop()

	first := true
	for {
		// If we haven't received a heartbeat ACK since the
		// last heartbeat we sent, we should consider the
		// connection as stale and return an error.
		t := atomic.LoadInt64(lastHeartbeatACK)
		if !first && time.Now().UTC().Sub(time.Unix(0, t).UTC()) > every {
			errCh <- errStaleConnection
			return
		}

		// Send the heartbeat payload.
		if err := h(); err != nil {
			errCh <- err
			return
		}

		if first {
			first = false
		}

		select {
		case <-stop:
			return
		case <-ticker.C:
		}
	}
}

// RunUDP periodically calls the given heartbeater to send a heartbeat packet.
// It should be called in a separate goroutine. It will decrement the given
// wait group when done, can be stopped by closing the stop channel and will
// report any error that occurs through the provided errCh.
func RunUDP(
	wg *sync.WaitGroup,
	stop chan struct{},
	errCh chan<- error,
	every time.Duration,
	h Hearbeater,
	lastUDPHeartbeatACK *int64,
) {
	defer wg.Done()

	ticker := time.NewTicker(every)
	defer ticker.Stop()

	first := true
	for {
		// If we haven't received a heartbeat ACK since the
		// last heartbeat we sent, we should consider the
		// connection as stale and return an error.
		// NOTE: since we're dealing with UDP, this might
		// not be the best idea. Maybe consider adding a threshold
		// before assuming the connection is down?
		t := atomic.LoadInt64(lastUDPHeartbeatACK)
		if !first && time.Now().UTC().Sub(time.Unix(0, t).UTC()) > every {
			errCh <- errStaleUDPConnection
			return
		}

		// Send the heartbeat packet.
		if err := h(); err != nil {
			// Silently break out of this loop because
			// the connection was closed by the client.
			if isConnectionClosed(err) {
				return
			}

			errCh <- err
			return
		}

		if first {
			first = false
		}

		select {
		case <-stop:
			return
		case <-ticker.C:
		}
	}
}

func isConnectionClosed(err error) bool {
	if e, ok := err.(*net.OpError); ok {
		// Ugly but : https://github.com/golang/go/issues/4373
		if e.Err.Error() == "use of closed network connection" {
			return true
		}
	}
	return false
}
