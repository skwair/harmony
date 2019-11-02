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
	"sync/atomic"
	"time"

	"nhooyr.io/websocket"

	"github.com/skwair/harmony/internal/payload"
	"github.com/skwair/harmony/log"
)

var silenceFrame = []byte{0xf8, 0xff, 0xfe}

// Connection represents a Discord voice connection.
type Connection struct {
	// General lock for long operations that should
	// not happen concurrently like Close.
	mu sync.Mutex

	// Send is used to send Opus encoded audio packets.
	Send chan []byte
	// Recv is used to receive audio packets
	// containing Opus encoded audio data.
	Recv chan *AudioPacket

	// User and session this connection is established with.
	userID, sessionID string
	// guild and channel IDs this voice connection is attached to.
	guildID, channelID string

	// Token used to identify to the voice server.
	token string
	// Websocket endpoint to connect to.
	endpoint string

	connRMu sync.Mutex
	conn    *websocket.Conn
	// Accessed atomically, acts as a boolean and is
	// set to 1 when the client is connected to voice.
	connected int32

	// UDP connection voice data is sent across.
	udpConn *net.UDPConn

	// Accessed atomically, acts as a boolean and
	// is set to 1 when the client is speaking.
	speaking int32

	// Secret used to encrypt voice data.
	secret [32]byte
	// SSRC of this user.
	ssrc uint32

	// Accessed atomically, acts as a boolean
	// and is set to 1 when the connection is
	// being established.
	connectingToVoice int32
	// When connectingToVoice is set to 1, some
	// payloads received by the event handler will
	// be sent through this channel.
	payloads chan *payload.Payload

	// wg keeps track of all goroutines that are
	// started when establishing a voice connection.
	wg sync.WaitGroup
	// The first fatal error encountered when connected
	// to a voice server will be reported to this channel.
	error chan error
	// Closing this channel will stop the voice connection.
	stop chan struct{}

	// Shared context used for sending and receiving websocket
	// payloads. Will be canceled when the client disconnects
	// or an error occurs.
	ctx    context.Context
	cancel context.CancelFunc

	// Accessed atomically, UNIX timestamps in nanoseconds.
	lastHeartbeatACK, lastUDPHeartbeatACK int64
	// Accessed atomically, sequence number of the last
	// UDP heartbeat we sent.
	udpHeartbeatSequence uint64

	// opusReadinessWG is a wait group used to make sure
	// the Opus sender and receiver are correctly started
	// before assuming we are connected to the voice channel.
	opusReadinessWG sync.WaitGroup

	logger log.Logger
}

// ConnectionOption is a function that configures a Connection.
// It is used in EstablishNewConnection.
type ConnectionOption func(*Connection)

// WithLogger can be used to set the logger used by this connection.
// Defaults to a standard logger reporting only errors.
// See the log package for more information about logging with Harmony.
func WithLogger(l log.Logger) ConnectionOption {
	return func(c *Connection) {
		c.logger = l
	}
}

// voiceIdentify is the payload sent to identify to a voice server.
type voiceIdentify struct {
	ServerID  string `json:"server_id"`
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Token     string `json:"token"`
}

// voiceReady payload is received when the client successfully identified
// with the voice server.
type voiceReady struct {
	SSRC  uint32   `json:"ssrc"`
	IP    string   `json:"ip"`
	Port  int      `json:"port"`
	Modes []string `json:"modes"`
}

// selectProtocol is sent by the client through the voice
// websocket to start the voice UDP connection.
type selectProtocol struct {
	Protocol string              `json:"protocol"`
	Data     *selectProtocolData `json:"data"`
}

type selectProtocolData struct {
	Address string `json:"address"`
	Port    uint16 `json:"port"`
	Mode    string `json:"mode"`
}

// sessionDescription is received when the client selected the UDP
// voice protocol. It contains the key to encrypt voice data.
type sessionDescription struct {
	Mode           string `json:"mode"`
	SecretKey      []byte `json:"secret_key"`
	VideoCodec     string `json:"video_codec"`
	AudioCodec     string `json:"audio_codec"`
	MediaSessionID string `json:"media_session_id"`
}

// EstablishNewConnection establishes a new voice connection with the provided
// information. This connection should be closed by calling its Close method
// when no longer needed.
func EstablishNewConnection(ctx context.Context, state *StateUpdate, server *ServerUpdate, opts ...ConnectionOption) (*Connection, error) {
	if state.ChannelID == nil {
		return nil, errors.New("could not establish voice connection: channel ID in given state is nil")
	}

	vc := &Connection{
		Send:      make(chan []byte, 2),
		Recv:      make(chan *AudioPacket),
		payloads:  make(chan *payload.Payload),
		error:     make(chan error),
		stop:      make(chan struct{}),
		userID:    state.UserID,
		sessionID: state.SessionID,
		guildID:   state.GuildID,
		channelID: *state.ChannelID,
		token:     server.Token,
		logger:    log.NewStd(os.Stderr, log.LevelError),
	}

	vc.ctx, vc.cancel = context.WithCancel(context.Background())

	for _, opt := range opts {
		opt(vc)
	}

	// Start by opening the voice websocket connection.
	var err error
	vc.endpoint = fmt.Sprintf("wss://%s?v=3", strings.TrimSuffix(server.Endpoint, ":80"))
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
			atomic.StoreInt32(&vc.connected, 0)
			close(vc.stop)
			vc.cancel()
		}
	}()

	// This is used to notify the event handler that some
	// specific payloads should be sent through to vc.payloads.
	atomic.StoreInt32(&vc.connectingToVoice, 1)
	defer atomic.StoreInt32(&vc.connectingToVoice, 0)

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
		V                 int `json:"v"`
		HeartbeatInterval int `json:"heartbeat_interval"`
	}
	if err = json.Unmarshal(p.D, &h); err != nil {
		return nil, err
	}
	// NOTE: do not start heartbeating before sending the identify payload
	// to the voice server, else it will close the connection.

	// Identify on the websocket connection. This is the first payload we must sent to the server.
	i := &voiceIdentify{
		ServerID:  vc.guildID,
		UserID:    vc.userID,
		SessionID: vc.sessionID,
		Token:     vc.token,
	}
	vc.logger.Debug("identifying to the voice server")
	if err = vc.sendPayload(ctx, voiceOpcodeIdentify, i); err != nil {
		return nil, err
	}

	// There is currently a bug in the Hello payload heartbeat interval.
	// See https://discordapp.com/developers/docs/topics/voice-connections#heartbeating
	every := float64(h.HeartbeatInterval) * .75
	// Now that we sent the identify payload, we can start heartbeating.
	vc.wg.Add(1)
	go vc.heartbeat(time.Duration(every) * time.Millisecond)

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
	addr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		return nil, err
	}

	vc.logger.Debugf("dialing voice connection endpoint: %s", host)
	vc.udpConn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	// From now on, close the UDP connection if any error occurs.
	defer func() {
		if err != nil {
			vc.udpConn.Close()
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

	if err = vc.sendSilenceFrame(ctx); err != nil {
		return nil, err
	}

	atomic.StoreInt32(&vc.connected, 1)
	vc.logger.Debug("connected to voice server")
	return vc, nil
}

// wait waits for an error to happen while connected to the voice server
// or for a stop signal to be sent.
func (vc *Connection) wait() {
	defer vc.wg.Done()

	vc.logger.Debug("starting voice connection manager")
	defer vc.logger.Debug("stopped voice connection manager")

	select {
	case err := <-vc.error:
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
	atomic.StoreInt32(&vc.connected, 0)

	// NOTE: maybe try to automatically reconnect if
	// we err != nil here, like done in the Gateway.
}

// onError is called when an error occurs while the connection to
// the voice server is up. It closes the underlying websocket connection
// with a 1006 code, logs the error and finally signals to all other
// goroutines (heartbeat, listen, etc.) to stop by closing the stop channel.
func (vc *Connection) onError(err error) {
	if closeErr := vc.conn.Close(websocket.StatusInternalError, "voice error"); closeErr != nil {
		vc.logger.Errorf("could not properly close voice websocket connection: %v", closeErr)
		vc.logger.Errorf("voice connection: %v", err)
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
	atomic.StoreUint64(&vc.udpHeartbeatSequence, 0)
}

// Must send some audio packets so the voice server starts to send us audio packets.
// This appears to be a bug from Discord.
func (vc *Connection) sendSilenceFrame(ctx context.Context) error {
	if err := vc.Speaking(ctx, true); err != nil {
		return err
	}

	vc.Send <- silenceFrame

	if err := vc.Speaking(ctx, false); err != nil {
		return err
	}

	return nil
}

// isConnecting returns whether this voice connection is currently connecting
// to a voice channel.
func (vc *Connection) isConnecting() bool {
	return atomic.LoadInt32(&vc.connectingToVoice) == 1
}

// Disconnect closes the voice connection.
func (vc *Connection) Close() {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	connected := atomic.LoadInt32(&vc.connected) == 1
	if !connected {
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
