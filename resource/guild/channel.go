package guild

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/skwair/harmony/discord"
	"github.com/skwair/harmony/internal/endpoint"
	"github.com/skwair/harmony/internal/rest"
)

// Channels returns the list of channels in the guild.
func (r *Resource) Channels(ctx context.Context) ([]discord.Channel, error) {
	e := endpoint.GetGuildChannels(r.guildID)
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var channels []discord.Channel
	if err = json.NewDecoder(resp.Body).Decode(&channels); err != nil {
		return nil, err
	}
	return channels, nil
}

// NewChannel is like NewChannelWithReason but with no particular reason.
func (r *Resource) NewChannel(ctx context.Context, settings *discord.ChannelSettings) (*discord.Channel, error) {
	return r.NewChannelWithReason(ctx, settings, "")
}

// NewChannelWithReason creates a new channel in the guild. Requires the MANAGE_CHANNELS permission.
// Fires a Channel Create Gateway event.
// The given reason will be set in the audit log entry for this action.
func (r *Resource) NewChannelWithReason(ctx context.Context, settings *discord.ChannelSettings, reason string) (*discord.Channel, error) {
	b, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	e := endpoint.CreateGuildChannel(r.guildID)
	resp, err := r.client.DoWithHeader(ctx, e, rest.JSONPayload(b), rest.ReasonHeader(reason))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, discord.NewAPIError(resp)
	}

	var ch discord.Channel
	if err = json.NewDecoder(resp.Body).Decode(&ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

// ChannelPosition is a pair of channel ID with its position.
type ChannelPosition struct {
	ID       string `json:"id"`
	Position int    `json:"position"`
}

// ModifyChannelPosition modifies the positions of a set of channel for the guild.
// Requires 'MANAGE_CHANNELS' permission. Fires multiple Channel Update Gateway events.
//
// Only channels to be modified are required, with the minimum being a swap between at
// least two channels.
func (r *Resource) ModifyChannelPosition(ctx context.Context, pos []ChannelPosition) error {
	b, err := json.Marshal(pos)
	if err != nil {
		return err
	}

	e := endpoint.ModifyChannelPositions(r.guildID)
	resp, err := r.client.Do(ctx, e, rest.JSONPayload(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return discord.NewAPIError(resp)
	}
	return nil
}
