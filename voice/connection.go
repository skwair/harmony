package voice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
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
// See https://discordapp.com/developers/docs/topics/voice-connections#voice-data-interpolation for more information.
var SilenceFrame = []byte{0xf8, 0xff, 0xfe}

// Connection represents a Discord voice connection.
type Connection struct {
	// General lock for long operations that should
	// not happen concurrently like Close or SetSpeakingMode.
	mu sync.Mutex

	// Send is used to send Opus encoded audio packets.
	Send chan []byte
	// Recv is used to receive audio packets
	// containing Opus encoded audio data.
	Recv chan *AudioPacket

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
	speakingMode SpeakingMode

	// Secret used to encrypt voice data.
	secret [32]byte
	// SSRC of this user.
	ssrc uint32

	// Whether this voice connection is up.
	connected *atomic.Bool
	// Whether this voice connection is currently
	// connecting to a voice server.
	connectingToVoice *atomic.Bool
	// Whether this voice connection is currently
	// trying to reconnect to a voice server.
	reconnecting *atomic.Bool

	// When connectingToVoice is set to 1, some
	// payloads received by the event handler will
	// be sent through this channel.
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

// Connect establishes a new voice connection with the provided
// information. This connection should be closed by calling its Close method
// when no longer needed.
func Connect(ctx context.Context, state *StateUpdate, server *ServerUpdate, opts ...ConnectionOption) (*Connection, error) {
	if state.ChannelID == nil {
		return nil, errors.New("could not establish voice connection: channel ID in given state is nil")
	}

	vc := &Connection{
		Send:                 make(chan []byte),
		Recv:                 make(chan *AudioPacket),
		payloads:             make(chan *payload.Payload),
		error:                make(chan error),
		stop:                 make(chan struct{}),
		state:                &state.State,
		token:                server.Token,
		logger:               log.NewStd(os.Stderr, log.LevelError),
		lastHeartbeatACK:     atomic.NewInt64(0),
		udpHeartbeatSequence: atomic.NewUint64(0),
		lastUDPHeartbeatACK:  atomic.NewInt64(0),
		connected:            atomic.NewBool(false),
		connectingToVoice:    atomic.NewBool(false),
		reconnecting:         atomic.NewBool(false),
	}

	vc.ctx, vc.cancel = context.WithCancel(context.Background())

	for _, opt := range opts {
		opt(vc)
	}

	// Start by opening the voice websocket connection.
	var err error
	vc.endpoint = fmt.Sprintf("wss://%s?v=%d", strings.TrimSuffix(server.Endpoint, ":80"), gatewayVersion)
	vc.logger.Debugf("connecting to voice server: %s", vc.endpoint)
	vc.conn, _, err = websocket.Dial(ctx, vc.endpoint, nil)
	if err != nil {
		return nil, err
	}
	// From now on, if any error occurs during the rest of the
	// voice connection process, we should close the underlying
	// websocket so we can try to reconnect.
	defer func() {
		if err != nil {
			_ = vc.conn.Close(websocket.StatusInternalError, "failed to establish voice connection")
			vc.connected.Store(false)
			close(vc.stop)
			vc.cancel()
		}
	}()

	// This is used to notify the event handler that some
	// specific payloads should be sent through to vc.payloads
	// while we are connecting to the voice server.
	vc.connectingToVoice.Store(true)
	defer vc.connectingToVoice.Store(false)

	vc.wg.Add(2) // listen starts an additional goroutine.
	go vc.listen()

	vc.wg.Add(1)
	go vc.wait()

	// The voice server should first send us a Hello packet defining the heartbeat
	// interval when we connect to the websocket.
	p := <-vc.payloads
	if p.Op != voiceOpcodeHello {
		return nil, fmt.Errorf("expected Opcode 8 Hello; got Opcode %d", p.Op)
	}

	var h struct {
		V                 int     `json:"v"`
		HeartbeatInterval float64 `json:"heartbeat_interval"`
	}
	if err = json.Unmarshal(p.D, &h); err != nil {
		return nil, err
	}
	// NOTE: do not start heartbeating before sending the identify payload
	// to the voice server, else it will close the connection.

	// Identify on the websocket connection. This is the first payload we must sent to the server.
	i := &voiceIdentify{
		ServerID:  vc.State().GuildID,
		UserID:    vc.State().UserID,
		SessionID: vc.State().SessionID,
		Token:     vc.token,
	}
	vc.logger.Debug("identifying to the voice server")
	if err = vc.sendPayload(ctx, voiceOpcodeIdentify, i); err != nil {
		return nil, err
	}

	// Now that we sent the identify payload, we can start heartbeating.
	vc.wg.Add(1)
	go vc.heartbeat(time.Duration(h.HeartbeatInterval) * time.Millisecond)

	// A Ready payload should be sent after we identified.
	p = <-vc.payloads
	if p.Op != voiceOpcodeReady {
		return nil, fmt.Errorf("expected Opcode 2 Ready; got Opcode %d", p.Op)
	}

	var vr voiceReady
	if err = json.Unmarshal(p.D, &vr); err != nil {
		return nil, err
	}
	vc.ssrc = vr.SSRC

	// We should now be able to open the voice UDP connection.
	host := fmt.Sprintf("%s:%d", vr.IP, vr.Port)
	vc.logger.Debug("resolving voice connection UDP endpoint")
	vc.dataEndpoint, err = net.ResolveUDPAddr("udp", host)
	if err != nil {
		return nil, err
	}

	vc.logger.Debugf("dialing voice connection endpoint: %s", host)
	vc.udpConn, err = net.DialUDP("udp", nil, vc.dataEndpoint)
	if err != nil {
		return nil, err
	}
	// From now on, close the UDP connection if any error occurs.
	defer func() {
		if err != nil {
			_ = vc.udpConn.Close()
		}
	}()

	// IP discovery.
	vc.logger.Debug("starting IP discovery")
	ip, port, err := ipDiscovery(vc.udpConn, vc.ssrc)
	if err != nil {
		return nil, err
	}
	vc.logger.Debugf("IP discovery result: %s:%d", ip, port)

	// Start heartbeating on the UDP connection.
	vc.wg.Add(1)
	go vc.udpHeartbeat(5 * time.Second)

	sp := &selectProtocol{
		Protocol: "udp",
		Data: &selectProtocolData{
			Address: ip,
			Port:    port,
			Mode:    "xsalsa20_poly1305",
		},
	}
	if err = vc.sendPayload(ctx, voiceOpcodeSelectProtocol, sp); err != nil {
		return nil, err
	}

	// Now we should receive a Session Description packet.
	p = <-vc.payloads
	if p.Op != voiceOpcodeSessionDescription {
		return nil, fmt.Errorf("expected Opcode 4 Session Description; got Opcode %d", p.Op)
	}

	var sd sessionDescription
	if err = json.Unmarshal(p.D, &sd); err != nil {
		return nil, err
	}

	copy(vc.secret[:], sd.SecretKey[0:32])

	vc.wg.Add(3) // opusReceiver starts an additional goroutine.
	vc.opusReadinessWG.Add(2)
	go vc.opusReceiver()
	go vc.opusSender()

	// Making sure Opus receiver and sender are started.
	vc.opusReadinessWG.Wait()

	if err = vc.sendSilenceFrame(); err != nil {
		return nil, err
	}

	vc.connected.Store(true)

	vc.logger.Debug("connected to voice server")
	return vc, nil
}

// wait waits for an error to happen while connected to the voice server
// or for a stop signal to be sent.
func (vc *Connection) wait() {
	defer vc.wg.Done()

	vc.logger.Debug("starting voice connection manager")
	defer vc.logger.Debug("stopped voice connection manager")

	var err error
	select {
	case err = <-vc.error:
		vc.onError(err)

	case <-vc.stop:
		vc.logger.Debug("disconnecting from the voice server")
		vc.onDisconnect()
	}

	close(vc.payloads)

	if vc.udpConn != nil {
		if err := vc.udpConn.Close(); err != nil {
			vc.logger.Errorf("failed to properly close voice UDP connection: %v", err)
		}
	}

	vc.cancel()
	vc.connected.Store(false)

	// If there was an error, maybe try to reconnect.
	if shouldReconnect(err) && !vc.isReconnecting() {
		vc.reconnectWithBackoff()
	}
}

// Determine whether we should try to reconnect based on the error we got.
// See https://discordapp.com/developers/docs/topics/opcodes-and-status-codes#voice-voice-close-event-codes for more information.
func shouldReconnect(err error) bool {
	if err == nil {
		return false
	}

	switch websocket.CloseStatus(err) {
	case 4003, 4004, 4005, 4006, 4011, 4012, 4014, 4016:
		return false
	case 4015:
		return true
	case -1: // Not a websocket error.
		return true
	default: // New (or undocumented?) close status code.
		return true
	}
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

// onError is called when an error occurs while the connection to
// the voice server is up. It closes the underlying websocket connection
// with a 1006 code, logs the error and finally signals to all other
// goroutines (heartbeat, listen, etc.) to stop by closing the stop channel.
func (vc *Connection) onError(err error) {
	vc.logger.Errorf("voice connection error: %v", err)

	if closeErr := vc.conn.Close(websocket.StatusInternalError, "voice error"); closeErr != nil {
		vc.logger.Errorf("could not properly close voice websocket connection: %v", closeErr)
	}

	// If an error occurred before the connection is established,
	// the stop channel will already be closed, so return early.
	if !vc.isEstablished() {
		return
	}

	close(vc.stop)
}

// onDisconnect is called when a normal disconnection happens (the client
// called the Close() method). It closes the underlying websocket
// connection with a 1000 code and resets the UDP heartbeat sequence.
func (vc *Connection) onDisconnect() {
	if err := vc.conn.Close(websocket.StatusNormalClosure, "disconnecting"); err != nil {
		vc.logger.Errorf("could not properly close voice websocket connection: %v", err)
	}
	vc.udpHeartbeatSequence.Store(0)
}

// Must send some audio packets so the voice server starts to send us audio packets.
// This appears to be a bug from Discord.
func (vc *Connection) sendSilenceFrame() error {
	if err := vc.SetSpeakingMode(SpeakingModeMicrophone); err != nil {
		return err
	}

	vc.Send <- SilenceFrame

	if err := vc.SetSpeakingMode(SpeakingModeOff); err != nil {
		return err
	}

	return nil
}

// Disconnect closes the voice connection.
func (vc *Connection) Close() {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	if !vc.isEstablished() {
		return
	}

	close(vc.stop)
	vc.wg.Wait()
	// NOTE: maybe we should explicitly close
	// other channels here.
	close(vc.Recv)
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

// isConnecting returns whether this voice connection is currently connecting
// to a voice channel.
func (vc *Connection) isConnecting() bool {
	return vc.connectingToVoice.Load()
}

// isEstablished reports whether the voice connection is fully established.
func (vc *Connection) isEstablished() bool {
	return vc.connected.Load()
}

func (vc *Connection) isReconnecting() bool {
	return vc.reconnecting.Load()
}
