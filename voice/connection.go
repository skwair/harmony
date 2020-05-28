package voice

import (
	"context"
	"net"
	"sync"
	"time"

	"go.uber.org/atomic"
	"nhooyr.io/websocket"

	"github.com/skwair/harmony/internal/payload"
	"github.com/skwair/harmony/log"
)

const (
	gatewayVersion = 4
)

// Five silence frames should be sent when there is a break in the sent data.
// See https://discord.com/developers/docs/topics/voice-connections#voice-data-interpolation for more information.
var SilenceFrame = []byte{0xf8, 0xff, 0xfe}

// Connection represents a Discord voice connection.
type Connection struct {
	// Send is used to send Opus encoded audio packets.
	Send chan []byte
	// Recv is used to receive audio packets
	// containing Opus encoded audio data.
	Recv chan *AudioPacket

	// General lock for long operations that should
	// not happen concurrently like Close or SetSpeakingMode.
	mu sync.Mutex

	// Current state of this voice connection.
	// This is set when initially establishing
	// the connection and should be kept up to
	// date with the SetState method as voice
	// server update events are received on the
	// main Gateway connection.
	stateMu sync.RWMutex
	state   *State

	// Token used to identify to the voice server.
	token string
	// Websocket endpoint to connect to.
	endpoint string
	// UDP endpoint to send voice data to.
	dataEndpoint *net.UDPAddr

	connRMu sync.Mutex
	conn    *websocket.Conn

	// UDP connection voice data is sent across.
	udpConn *net.UDPConn

	// Holds the value of the last
	// speaking (opcode 5) payload sent.
	speakingModeMu sync.Mutex
	speakingMode   SpeakingMode

	// Secret used to encrypt voice data.
	secret [32]byte
	// SSRC of this user.
	ssrc uint32

	// Whether this voice connection is up.
	connected *atomic.Bool
	// Whether this voice connection is currently
	// connecting to a voice server.
	connecting *atomic.Bool
	// Whether this voice connection is currently
	// trying to reconnect to a voice server.
	reconnecting *atomic.Bool

	// When connecting is true, some payloads
	// received by the event handler will be
	// sent through this channel.
	payloads chan *payload.Payload

	// wg keeps track of all goroutines that are
	// started when establishing a voice connection.
	wg sync.WaitGroup
	// The first fatal error encountered while connected
	// to a voice server will be sent through this channel.
	error chan error
	// Used to ensure we only report the first error.
	reportErrorOnce sync.Once

	// Closing this channel will gracefully shutdown the
	// voice connection.
	stop chan struct{}

	// Shared context used for sending and receiving websocket
	// payloads. Will be canceled when the client disconnects
	// or an error occurs.
	ctx    context.Context
	cancel context.CancelFunc

	// Accessed atomically, UNIX timestamps in nanoseconds.
	lastHeartbeatACK, lastUDPHeartbeatACK *atomic.Int64
	// Accessed atomically, sequence number of the last
	// UDP heartbeat we sent.
	udpHeartbeatSequence *atomic.Uint64

	// opusReadinessWG is a wait group used to make sure
	// the Opus sender and receiver are correctly started
	// before assuming we are connected to the voice channel.
	opusReadinessWG sync.WaitGroup

	logger log.Logger
}

// Logger is here to make the logger available to third party packages that
// need to report errors related to this voice connection.
func (vc *Connection) Logger() log.Logger {
	return vc.logger
}

// State returns the current state of this voice connection.
func (vc *Connection) State() *State {
	vc.stateMu.RLock()
	defer vc.stateMu.RUnlock()

	return vc.state
}

// SetState updates the state for this voice connections.
func (vc *Connection) SetState(s *State) {
	vc.stateMu.Lock()
	defer vc.stateMu.Unlock()

	vc.state = s
}

// UpdateServer updates the voice server this connection is using.
// This closes the connection to the old server and establishes a new
// connection to the updated server.
func (vc *Connection) UpdateServer(server *ServerUpdate) error {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	close(vc.stop)
	vc.wg.Wait()

	// Explicitly set the speaking mode to off in the
	// event where some audio is being sent while updating
	// the voice server. This will allow the connect method
	// to correctly send the initial silence frame.
	vc.speakingModeMu.Lock()
	vc.speakingMode = SpeakingModeOff
	vc.speakingModeMu.Unlock()

	vc.reset()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	return vc.connect(ctx, server)
}

// reset resets the voice connection so a new connect or reconnect attempt can be issued.
func (vc *Connection) reset() {
	vc.payloads = make(chan *payload.Payload)
	vc.error = make(chan error)
	vc.reportErrorOnce = sync.Once{}
	vc.stop = make(chan struct{})
	vc.lastHeartbeatACK = atomic.NewInt64(0)
	vc.udpHeartbeatSequence = atomic.NewUint64(0)
	vc.lastUDPHeartbeatACK = atomic.NewInt64(0)

	vc.ctx, vc.cancel = context.WithCancel(context.Background())
}

// reportErr reports the first fatal error encountered while a voice
// connection is up. Calls after the first one are no-ops.
func (vc *Connection) reportErr(err error) {
	vc.reportErrorOnce.Do(func() {
		select {
		case vc.error <- err:

		// Discard the error if we are already closing the connection.
		case <-vc.stop:
		}
		close(vc.error)
	})
}

// isConnecting returns whether this voice connection is currently connecting
// to a voice channel.
func (vc *Connection) isConnecting() bool {
	return vc.connecting.Load()
}

// isConnected reports whether the voice connection is fully established.
func (vc *Connection) isConnected() bool {
	return vc.connected.Load()
}

func (vc *Connection) isReconnecting() bool {
	return vc.reconnecting.Load()
}
