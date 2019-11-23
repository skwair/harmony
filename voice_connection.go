package harmony

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/skwair/harmony/internal/payload"
	"github.com/skwair/harmony/voice"
)

// JoinVoiceChannel will create a new voice connection to the given voice channel.
// This method is safe to call from multiple goroutines, but connections will happen
// sequentially.
// To properly leave the voice channel, call LeaveVoiceChannel.
func (c *Client) JoinVoiceChannel(ctx context.Context, guildID, channelID string, mute, deaf bool) (*voice.Connection, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected() {
		return nil, ErrGatewayNotConnected
	}

	// This is used to notify the already started event handler that
	// some specific payloads should be sent through to c.payloads.
	c.connectingToVoice.Store(true)
	defer c.connectingToVoice.Store(false)

	// Notify a voice server that we want to connect to a voice channel.
	vsu := &voice.StateUpdate{
		State: voice.State{
			GuildID:   guildID,
			ChannelID: &channelID,
			SelfMute:  mute,
			SelfDeaf:  deaf,
		},
	}
	if err := c.sendPayload(ctx, gatewayOpcodeVoiceStateUpdate, vsu); err != nil {
		return nil, err
	}

	// The voice server should answer with two payloads,
	// describing the voice state and the voice server
	// to connect to.
	state, server, err := getStateAndServer(c.voicePayloads)
	if err != nil {
		return nil, err
	}

	// Establish the voice connection.
	conn, err := voice.Connect(ctx, state, server, voice.WithLogger(c.logger))
	if err != nil {
		return nil, err
	}

	c.voiceConnections[guildID] = conn

	return conn, nil
}

// LeaveVoiceChannel notifies the Gateway we want the voice channel we are
// connected to in the given guild.
func (c *Client) LeaveVoiceChannel(ctx context.Context, guildID string) error {
	conn, ok := c.voiceConnections[guildID]
	if ok {
		conn.Close()
		delete(c.voiceConnections, guildID)
	}

	vsu := &voice.StateUpdate{
		State: voice.State{
			GuildID: guildID,
		},
	}
	if err := c.sendPayload(ctx, gatewayOpcodeVoiceStateUpdate, vsu); err != nil {
		return fmt.Errorf("could not send voice state update payload: %w", err)
	}
	return nil
}

// getStateAndServer will receive exactly two payloads from ch and extract the voice state
// and the voice server information from them. The order of the payloads is not relevant
// although only those two payloads must be sent through ch and only once each.
// NOTE: check if those events are always sequentially sent in the same order, if so,
// refactor this function.
func getStateAndServer(ch chan *payload.Payload) (*voice.StateUpdate, *voice.ServerUpdate, error) {
	var (
		server        voice.ServerUpdate
		state         voice.StateUpdate
		first, second bool
	)

	for i := 0; i < 2; i++ {
		p := <-ch
		if p.T == eventVoiceStateUpdate {
			if first {
				return nil, nil, errors.New("already received voice state update payload")
			}
			first = true

			if err := json.Unmarshal(p.D, &state); err != nil {
				return nil, nil, err
			}
		} else if p.T == eventVoiceServerUpdate {
			if second {
				return nil, nil, errors.New("already received voice server update payload")
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
