package voice

import (
	"context"
	"encoding/json"

	"github.com/skwair/harmony/internal/payload"
)

// StateUpdate is the payload describing the update of the voice state of a user.
type StateUpdate struct {
	State
}

// State represents the voice state of a user.
type State struct {
	GuildID    string  `json:"guild_id"`
	ChannelID  *string `json:"channel_id"`
	UserID     string  `json:"user_id"`
	SessionID  string  `json:"session_id"`
	Deaf       bool    `json:"deaf"`
	Mute       bool    `json:"mute"`
	SelfDeaf   bool    `json:"self_deaf"`
	SelfMute   bool    `json:"self_mute"`
	SelfStream bool    `json:"self_stream"`
	Suppress   bool    `json:"suppress"` // Whether this user is muted by the current user.
}

// Clone returns a clone of this StateUpdate.
func (v *State) Clone() *State {
	if v == nil {
		return nil
	}

	return &State{
		GuildID:   v.GuildID,
		ChannelID: v.ChannelID,
		UserID:    v.UserID,
		SessionID: v.SessionID,
		Deaf:      v.Deaf,
		Mute:      v.Mute,
		SelfDeaf:  v.SelfDeaf,
		SelfMute:  v.SelfMute,
		Suppress:  v.Suppress,
	}
}

// ServerUpdate is the payload describing the update of a voice server.
type ServerUpdate struct {
	Token    string `json:"token"`
	GuildID  string `json:"guild_id"`
	Endpoint string `json:"endpoint"`
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

// resume is sent by the client to notify a voice server it is trying
// to resume a connection which was unexpectedly ended.
type resume struct {
	ServerID  string `json:"server_id"`
	SessionID string `json:"session_id"`
	Token     string `json:"token"`
}

// sendPayload sends a single Payload to the Voice server with
// the given op and data.
func (vc *Connection) sendPayload(ctx context.Context, op int, d interface{}) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	p := &payload.Payload{Op: op, D: b}
	vc.logger.Debugf("sent voice payload (guild=%q): %s", vc.State().GuildID, p)
	return payload.Send(ctx, vc.conn, p)
}

// recvPayload receives a single Payload from the Voice server.
func (vc *Connection) recvPayload() (*payload.Payload, error) {
	p, err := payload.Recv(vc.ctx, &vc.connRMu, vc.conn)
	if err != nil {
		return nil, err
	}

	vc.logger.Debugf("received voice payload (guild=%q): %s", vc.State().GuildID, p)

	return p, nil
}

// listenAndHandlePayloads listens for payloads sent by the voice server and
// handles them as they are received.
func (vc *Connection) listenAndHandlePayloads() {
	defer vc.wg.Done()

	vc.logger.Debug("starting voice connection event listener")
	defer vc.logger.Debug("stopped voice connection event listener")

	payload.ListenAndHandle(vc.recvPayload, vc.handleEvent, vc.reportErr)
}
