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
	lastHeartbeatAck *atomic.Int64,
	lastHeartbeatSent *atomic.Int64,
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
		lastAck := time.Unix(0, lastHeartbeatAck.Load())
		lastSent := time.Unix(0, lastHeartbeatSent.Load())
		if !first && lastSent.After(lastAck) {
			errReporter(fmt.Errorf("no heartbeat received since %v (%v ago)", lastAck, time.Since(lastAck)))
			return
		}

		// Send the heartbeat payload.
		if err := h(); err != nil {
			errReporter(err)
			return
		}

		lastHeartbeatSent.Store(time.Now().UnixNano())

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
	lastUDPHeartbeatAck *atomic.Int64,
	lastUDPHeartbeatSent *atomic.Int64,
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
		lastAck := time.Unix(0, lastUDPHeartbeatAck.Load())
		lastSent := time.Unix(0, lastUDPHeartbeatSent.Load())
		if !first && lastSent.After(lastAck) {
			errReporter(fmt.Errorf("no UDP heartbeat received since %v (%v ago)", lastAck, time.Since(lastAck)))
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

		lastUDPHeartbeatSent.Store(time.Now().UnixNano())

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
