package heartbeat

import (
	"fmt"
	"net"
	"time"

	"go.uber.org/atomic"
)

// Hearbeater is a function that sends a heartbeat.
type Hearbeater func() error

// Run periodically calls the given heartbeater to send a heartbeat payload.
// It can be stopped by closing the stop channel and will report any error
// that occurs using the given errReporter.
func Run(
	every time.Duration,
	h Hearbeater,
	lastHeartbeatACK *atomic.Int64,
	stop chan struct{},
	errReporter func(err error),
) {
	ticker := time.NewTicker(every)
	defer ticker.Stop()

	first := true
	for {
		// If we haven't received a heartbeat ACK since the
		// last heartbeat we sent, we should consider the
		// connection as stale and return an error.
		lastACK := time.Unix(0, lastHeartbeatACK.Load())
		if !first && time.Since(lastACK) > every {
			errReporter(fmt.Errorf("no heartbeat received since %v (%v ago)", lastACK, time.Since(lastACK)))
			return
		}

		// Send the heartbeat payload.
		if err := h(); err != nil {
			errReporter(err)
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
// It can be stopped by closing the stop channel and will report any error that
// occurs using the given errReporter.
func RunUDP(
	every time.Duration,
	h Hearbeater,
	lastUDPHeartbeatACK *atomic.Int64,
	stop chan struct{},
	errReporter func(err error),
) {
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
		lastACK := time.Unix(0, lastUDPHeartbeatACK.Load())
		if !first && time.Since(lastACK) > every {
			errReporter(fmt.Errorf("no UDP heartbeat received since %v (%v ago)", lastACK, time.Since(lastACK)))
			return
		}

		// Send the heartbeat packet.
		if err := h(); err != nil {
			// Silently break out of this loop because
			// the connection was closed by the client.
			if isConnectionClosed(err) {
				return
			}

			errReporter(err)
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
