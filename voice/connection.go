package voice

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"

	"github.com/skwair/harmony/internal/payload"
	"github.com/skwair/harmony/log"
)

var silenceFrame = []byte{0xf8, 0xff, 0xfe}

// Connection represents a Discord voice connection.
type Connection struct {
	// General lock for long operations that should
	// not happen concurrently like Disconnect.
	mu sync.Mutex

	// Send is used to send Opus encoded audio packets.
	Send chan []byte
	// Recv is used to receive audio packets
	// containing Opus encoded audio data.
	Recv chan *AudioPacket

	// User and session that initiated this connection.
	userID, sessionID string
	// guild and channel ID this voice
	// connection is attached to.
	guildID, channelID string

	// This is the the token used to connect
	// to the voice server. Used when resuming
	// voice connections.
	token string
	// Websocket endpoint to connect to.
	endpoint string

	connRMu sync.Mutex
	connWMu sync.Mutex
	conn    *websocket.Conn
	// Accessed atomically, acts as a boolean and is
	// set to 1 when the client is connected to voice.
	connected int32

	udpConn *net.UDPConn

	// Accessed atomically, acts as a boolean and
	// is set to 1 when the client is speaking.
	speaking int32

	// Secret used to encrypt voice data.
	secret [32]byte
	// ssrc of this user.
	ssrc uint32

	// Accessed atomically, acts as a boolean
	// and is set to 1 when the client is
	// connecting to voice.
	connectingToVoice int32
	// When connectingToVoice is set to 1, some
	// payloads received by the event handler will
	// be sent through this channel.
	payloads chan *payload.Payload

	// wg keeps track of all goroutines that are
	// started when connecting to a voice channel.
	wg sync.WaitGroup
	// The first fatal error encountered will be reported to this channel.
	error chan error
	// Closing this channel will stop the voice connection.
	stop chan struct{}

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
func EstablishNewConnection(state *State, server *ServerUpdate) (*Connection, error) {
	vc := &Connection{
		Send:      make(chan []byte, 2),
		Recv:      make(chan *AudioPacket),
		payloads:  make(chan *payload.Payload),
		error:     make(chan error),
		stop:      make(chan struct{}),
		logger:    log.NewStd(log.LevelDebug), // FIXME: set a noop logger by default?
		userID:    state.UserID,
		sessionID: state.SessionID,
		guildID:   state.GuildID,
		channelID: state.ChannelID,
		token:     server.Token,
	}

	// Open the voice websocket connection.
	var err error
	vc.endpoint = fmt.Sprintf("wss://%s?v=3", strings.TrimSuffix(server.Endpoint, ":80"))
	vc.logger.Debugf("connecting to voice server: %s", vc.endpoint)
	vc.conn, _, err = websocket.DefaultDialer.Dial(vc.endpoint, nil)
	if err != nil {
		return nil, err
	}

	// From now, if any error occurs during the rest of the
	// voice connection process, we should close the underlying
	// websocket so we can try to reconnect.
	defer func() {
		if err != nil {
			vc.conn.Close()
			atomic.StoreInt32(&vc.connected, 0)
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

	i := &voiceIdentify{
		ServerID:  vc.guildID,
		UserID:    vc.userID,
		SessionID: vc.sessionID,
		Token:     vc.token,
	}
	vc.logger.Debug("identifying to the voice server")
	if err = vc.sendPayload(voiceOpcodeIdentify, i); err != nil {
		return nil, err
	}

	// There is currently a bug in the Hello payload heartbeat interval.
	// See https://discordapp.com/developers/docs/topics/voice-connections#heartbeating
	every := float64(h.HeartbeatInterval) * .75
	// Now we can start heartbeating.
	vc.wg.Add(1)
	go vc.heartbeat(time.Duration(every) * time.Millisecond)

	// Now we should receive a Ready packet.
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
	ipPort := fmt.Sprintf("%s:%d", ip, port)
	vc.logger.Debugf("IP discovery result: %s", ipPort)

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
	if err = vc.sendPayload(voiceOpcodeSelectProtocol, sp); err != nil {
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

	atomic.StoreInt32(&vc.connected, 1)
	vc.logger.Debug("connected to voice server")
	return vc, nil
}

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

	if err := vc.conn.Close(); err != nil {
		vc.logger.Errorf("failed to properly close voice connection: %v", err)
	}
	if vc.udpConn != nil {
		if err := vc.udpConn.Close(); err != nil {
			vc.logger.Errorf("failed to properly close voice UDP connection: %v", err)
		}
	}

	atomic.StoreInt32(&vc.connected, 0)

	// NOTE: maybe try to automatically reconnect if
	// we err != nil here, like done in the Gateway.
}

func (vc *Connection) onError(err error) {
	if err := vc.conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseAbnormalClosure, ""),
		time.Now().Add(time.Second*10),
	); err != nil {
		vc.logger.Errorf("could not properly close voice websocket: %v", err)
	}
	vc.logger.Errorf("voice connection: %v", err)
	close(vc.stop)
}

func (vc *Connection) onDisconnect() {
	if err := vc.conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(time.Second*10),
	); err != nil {
		vc.logger.Errorf("could not properly close voice websocket: %v", err)
	}
	atomic.StoreUint64(&vc.udpHeartbeatSequence, 0)
}

// Must send some audio packets so the voice server starts to send us audio packets.
// This appears to be a bug from Discord.
func (vc *Connection) sendSilenceFrame() error {
	if err := vc.Speaking(true); err != nil {
		return err
	}

	vc.Send <- silenceFrame

	if err := vc.Speaking(false); err != nil {
		return err
	}

	return nil
}

// isConnecting returns whether this voice connection is currently connecting
// to a voice channel.
func (vc *Connection) isConnecting() bool {
	return atomic.LoadInt32(&vc.connectingToVoice) == 1
}

// Speaking sends an Opcode 5 Speaking payload. This does nothing
// if the user is already in the given state.
func (vc *Connection) Speaking(s bool) error {
	// Return early if the user is already in the asked state.
	prev := atomic.LoadInt32(&vc.speaking)
	if (prev == 1) == s {
		return nil
	}

	if s {
		atomic.StoreInt32(&vc.speaking, 1)
	} else {
		atomic.StoreInt32(&vc.speaking, 0)
	}

	p := struct {
		Speaking bool   `json:"speaking"`
		Delay    int    `json:"delay"`
		SSRC     uint32 `json:"ssrc"`
	}{
		Speaking: s,
		Delay:    0,
		SSRC:     vc.ssrc,
	}

	if err := vc.sendPayload(voiceOpcodeSpeaking, p); err != nil {
		// If there is an error, reset our internal value to its previous
		// state because the update was not acknowledged by Discord.
		atomic.StoreInt32(&vc.speaking, prev)
		return err
	}

	return nil
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

// State represents the voice state of a user.
type State struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Deaf      bool   `json:"deaf"`
	Mute      bool   `json:"mute"`
	SelfDeaf  bool   `json:"self_deaf"`
	SelfMute  bool   `json:"self_mute"`
	Suppress  bool   `json:"suppress"` // Whether this user is muted by the current user.
}

type ServerUpdate struct {
	Token    string `json:"token"`
	GuildID  string `json:"guild_id"`
	Endpoint string `json:"endpoint"`
}

// StateUpdate is sent to notify a voice server that
// the client wants to connect to a voice channel.
type StateUpdate struct {
	GuildID   string  `json:"guild_id"`
	ChannelID *string `json:"channel_id"`
	SelfMute  bool    `json:"self_mute"`
	SelfDeaf  bool    `json:"self_deaf"`
}
