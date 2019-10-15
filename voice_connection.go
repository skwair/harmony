package harmony

import (
	"encoding/binary"
	"encoding/json"
	"errors"
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

// VoiceConnection represents a Discord voice connection.
type VoiceConnection struct {
	// General lock for long operations that should
	// not happen concurrently like Disconnect.
	mu sync.Mutex

	// Send is used to send Opus encoded audio packets.
	Send chan []byte
	// Recv is used to receive audio packets
	// containing Opus encoded audio data.
	Recv chan *AudioPacket

	// Keep a reference to the harmony Client
	// so this voice connection is able to
	// send some payloads to the Gateway
	// (needed when disconnecting).
	client *Client

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

	mute, deaf bool

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

// Logger is here to make the logger available to third party packages that
// need to report errors related to this voice connection.
func (vc *VoiceConnection) Logger() log.Logger {
	return vc.logger
}

// VoiceConnectionOption is a function that configures a VoiceConnection.
// It is used in ConnectToVoice.
type VoiceConnectionOption func(*VoiceConnection)

// WithMute allows you to specify whether this voice connection should be muted when connecting.
func WithMute(t bool) VoiceConnectionOption {
	return func(c *VoiceConnection) {
		c.mute = t
	}
}

// WithDeaf allows you to specify whether this voice connection should be deafened when connecting.
func WithDeaf(t bool) VoiceConnectionOption {
	return func(c *VoiceConnection) {
		c.deaf = t
	}
}

// voiceStateUpdate is sent to notify a voice server that
// the client wants to connect to a voice channel.
type voiceStateUpdate struct {
	GuildID   string  `json:"guild_id"`
	ChannelID *string `json:"channel_id"`
	SelfMute  bool    `json:"self_mute"`
	SelfDeaf  bool    `json:"self_deaf"`
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

// ConnectToVoice will create a new VoiceConnection to the given voice channel.
// This method is safe to call from multiple goroutines, but connections will happen
// sequentially.
func (c *Client) ConnectToVoice(guildID, channelID string, opts ...VoiceConnectionOption) (*VoiceConnection, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected() {
		return nil, ErrGatewayNotConnected
	}

	// Initialize a new VoiceConnection.
	vc := VoiceConnection{
		Send:      make(chan []byte, 2),
		Recv:      make(chan *AudioPacket),
		payloads:  make(chan *payload.Payload),
		error:     make(chan error),
		stop:      make(chan struct{}),
		client:    c,
		guildID:   guildID,
		channelID: channelID,
		mute:      false,
		deaf:      false,
		logger:    c.logger,
	}

	for _, opt := range opts {
		opt(&vc)
	}

	if err := vc.connect(); err != nil {
		return nil, err
	}
	return &vc, nil
}

var silenceFrame = []byte{0xf8, 0xff, 0xfe}

func (vc *VoiceConnection) connect() error {
	// This is used to notify the already started event handler that
	// some specific payloads should be sent through to c.payloads.
	atomic.StoreInt32(&vc.client.connectingToVoice, 1)
	defer atomic.StoreInt32(&vc.client.connectingToVoice, 0)

	// Notify a voice server that we want to connect to a voice channel.
	vsu := &voiceStateUpdate{
		GuildID:   vc.guildID,
		ChannelID: &vc.channelID,
		SelfMute:  vc.mute,
		SelfDeaf:  vc.deaf,
	}
	if err := vc.client.sendPayload(gatewayOpcodeVoiceStateUpdate, vsu); err != nil {
		return err
	}

	// The voice server should answer with two payloads,
	// describing the voice state and the voice server
	// to connect to.
	_, server, err := getStateAndServer(vc.client.voicePayloads)
	if err != nil {
		return err
	}
	vc.token = server.Token

	// Open the voice websocket connection.
	vc.endpoint = fmt.Sprintf("wss://%s?v=3", strings.TrimSuffix(server.Endpoint, ":80"))
	vc.logger.Debugf("connecting to voice server: %s", vc.endpoint)
	vc.conn, _, err = websocket.DefaultDialer.Dial(vc.endpoint, nil)
	if err != nil {
		return err
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

	i := &voiceIdentify{
		ServerID:  vc.guildID,
		UserID:    vc.client.userID,
		SessionID: vc.client.sessionID,
		Token:     vc.token,
	}
	vc.logger.Debug("identifying to the voice server")
	if err = vc.sendPayload(voiceOpcodeIdentify, i); err != nil {
		return err
	}

	// The Gateway should send us a Hello packet defining the heartbeat
	// interval when we connect to the websocket.
	p := <-vc.payloads
	if p.Op != voiceOpcodeHello {
		return fmt.Errorf("expected Opcode 8 Hello; got Opcode %d", p.Op)
	}

	var h struct {
		V                 int `json:"v"`
		HeartbeatInterval int `json:"heartbeat_interval"`
	}
	if err = json.Unmarshal(p.D, &h); err != nil {
		return err
	}
	// There is currently a bug in the Hello payload heartbeat interval.
	// See https://discordapp.com/developers/docs/topics/voice-connections#heartbeating
	every := float64(h.HeartbeatInterval) * .75
	vc.wg.Add(1)
	go vc.heartbeat(time.Duration(every) * time.Millisecond)

	// Now we should receive a Ready packet.
	p = <-vc.payloads
	if p.Op != voiceOpcodeReady {
		return fmt.Errorf("expected Opcode 2 Ready; got Opcode %d", p.Op)
	}

	var vr voiceReady
	if err = json.Unmarshal(p.D, &vr); err != nil {
		return err
	}
	vc.ssrc = vr.SSRC
	// We should now be able to open the voice UDP connection.
	host := fmt.Sprintf("%s:%d", vr.IP, vr.Port)
	vc.logger.Debugf("resolving voice connection UDP endpoint: %s", host)
	addr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		return err
	}

	vc.logger.Debugf("dialing voice connection endpoint: %s", host)
	vc.udpConn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	// From now on, close the UDP connection if any error occurs.
	defer func() {
		if err != nil {
			vc.udpConn.Close()
		}
	}()

	// IP discovery.
	ip, port, err := ipDiscovery(vc.udpConn, vc.ssrc)
	if err != nil {
		return err
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
		return err
	}

	// Now we should receive a Session Description packet.
	p = <-vc.payloads
	if p.Op != voiceOpcodeSessionDescription {
		return fmt.Errorf("expected Opcode 4 Session Description; got Opcode %d", p.Op)
	}

	var sd sessionDescription
	if err = json.Unmarshal(p.D, &sd); err != nil {
		return err
	}

	copy(vc.secret[:], sd.SecretKey[0:32])

	vc.wg.Add(3) // opusReceiver starts an additional goroutine.
	vc.opusReadinessWG.Add(2)
	go vc.opusReceiver()
	go vc.opusSender()

	// Making sure Opus receiver and sender are started.
	vc.opusReadinessWG.Wait()

	if err = vc.sendSilenceFrame(); err != nil {
		return err
	}

	atomic.StoreInt32(&vc.connected, 1)
	vc.logger.Debug("connected to voice server")
	return nil
}

// Must send some audio packets so the voice server starts to send us audio packets.
// This appears to be a bug from Discord.
func (vc *VoiceConnection) sendSilenceFrame() error {
	if err := vc.Speaking(true); err != nil {
		return err
	}

	vc.Send <- silenceFrame

	if err := vc.Speaking(false); err != nil {
		return err
	}

	return nil
}

// Disconnect closes the voice connection.
func (vc *VoiceConnection) Disconnect() {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	connected := atomic.LoadInt32(&vc.connected) == 1
	if !connected {
		return
	}

	// Notify the voice server that we want to disconnect.
	vsu := &voiceStateUpdate{
		GuildID: vc.guildID,
	}
	if err := vc.client.sendPayload(gatewayOpcodeVoiceStateUpdate, vsu); err != nil {
		vc.logger.Errorf("voice connection: %v", err)
	}

	close(vc.stop)
	vc.wg.Wait()
	// NOTE: maybe we should explicitly close
	// other channels here.
	close(vc.Recv)
}

func (vc *VoiceConnection) wait() {
	defer vc.wg.Done()

	vc.logger.Debug("starting voice connection manager")
	defer vc.logger.Debug("stopped voice connection manager")

	select {
	case err := <-vc.error:
		vc.onError(err)

	case <-vc.stop:
		vc.onDisconnect()
	}

	vc.conn.Close()
	vc.udpConn.Close()

	atomic.StoreInt32(&vc.connected, 0)
}

func (vc *VoiceConnection) onError(err error) {
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

func (vc *VoiceConnection) onDisconnect() {
	if err := vc.conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(time.Second*10),
	); err != nil {
		vc.logger.Errorf("could not properly close voice websocket: %v", err)
	}
	atomic.StoreUint64(&vc.udpHeartbeatSequence, 0)
}

// ipDiscovery uses Discord's IP discovery service to get the external ip and port the
// given UDP connection is using.
func ipDiscovery(conn *net.UDPConn, ssrc uint32) (ip string, port uint16, err error) {
	b := make([]byte, 70)
	binary.BigEndian.PutUint32(b, ssrc)
	if _, err = conn.Write(b); err != nil {
		return "", 0, err
	}

	b = make([]byte, 70)
	l, err := conn.Read(b)
	if err != nil {
		return "", 0, err
	}
	if l < 70 {
		return "", 0, errors.New("ipDiscovery: did not receive enough bytes")
	}

	s := strings.Builder{}
	for i := 4; b[i] != 0; i++ {
		if err = s.WriteByte(b[i]); err != nil {
			return "", 0, err
		}
	}
	return s.String(), binary.LittleEndian.Uint16(b[68:70]), nil
}

// getStateAndServer will receive exactly two payloads from ch and extract the voice state
// and the voice server information from them. The order of the payloads is not relevant
// although only those two payloads must be sent through ch and only once each.
// NOTE: check if those events are always sequentially sent in the same order, if so,
// refactor this function.
func getStateAndServer(ch chan *payload.Payload) (*VoiceState, *VoiceServerUpdate, error) {
	var (
		server        VoiceServerUpdate
		state         VoiceState
		first, second bool
	)

	for i := 0; i < 2; i++ {
		p := <-ch
		if p.T == eventVoiceStateUpdate {
			if first {
				return nil, nil, errors.New("already received Voice Server Update payload")
			}
			first = true

			if err := json.Unmarshal(p.D, &state); err != nil {
				return nil, nil, err
			}
		} else if p.T == eventVoiceServerUpdate {
			if second {
				return nil, nil, errors.New("already received Voice State payload")
			}
			second = true

			if err := json.Unmarshal(p.D, &server); err != nil {
				return nil, nil, err
			}
		} else {
			return nil, nil, fmt.Errorf(
				"expected Opcode 0 VOICE_STATE_UPDATE or VOICE_SERVER_UPDATE; got Opcode %d %s",
				p.Op, p.T)
		}
	}
	return &state, &server, nil
}

// isConnecting returns whether this voice connection is currently connecting
// to a voice channel.
func (vc *VoiceConnection) isConnecting() bool {
	return atomic.LoadInt32(&vc.connectingToVoice) == 1
}
