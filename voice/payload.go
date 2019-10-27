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

// sendPayload sends a single Payload to the Voice server with
// the given op and data.
func (vc *Connection) sendPayload(ctx context.Context, op int, d interface{}) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	return payload.Send(ctx, vc.conn, &payload.Payload{Op: op, D: b})
}

// recvPayload receives a single Payload from the Voice server.
func (vc *Connection) recvPayload() (*payload.Payload, error) {
	return payload.Recv(vc.ctx, &vc.connRMu, vc.conn)
}
