package main

import (
	"context"
	"fmt"

	"github.com/skwair/harmony/discord"
	"github.com/skwair/harmony/resource/channel"
	"github.com/skwair/harmony/voice"
)

func (b *bot) handleCommand(ctx context.Context, m *discord.Message) {
	switch m.Content {
	case "!play":
		// First, find the voice channel the user is in.
		channelID := findUser(b.client.State.Guild(m.GuildID).VoiceStates, m.Author.ID)
		if channelID == "" {
			b.sendError(ctx, m.ChannelID, "user must be in a voice channel")
			return
		}

		// Then, check if we already have an active player in this guild.
		p, ok := b.players[m.GuildID]
		if !ok {
			// If we don't, create a new one.
			var err error
			p, err = b.createNewPlayer(ctx, m.GuildID, channelID)
			if err != nil {
				b.sendError(ctx, m.ChannelID, err.Error())
				return
			}

			b.players[m.GuildID] = p
		} else {
			// If we already have an active player in this guild, then switch it to the user's channel.
			// This is a no-op if it already was in the correct channel.
			if err := b.client.SwitchVoiceChannel(ctx, m.GuildID, channelID); err != nil {
				fmt.Println("could not join voice channel: ", err)
				return
			}
		}

		// Start playing a track with the player.
		if err := p.Play("./happyrock.mp3"); err != nil {
			b.sendError(ctx, m.ChannelID, err.Error())
			return
		}

	case "!stop":
		// If we have a player in this guild, stop it, else reply with an error.
		if player, ok := b.players[m.GuildID]; ok {
			player.Stop()
		} else {
			b.sendError(ctx, m.ChannelID, "bot not playing music in this guild")
		}

	case "!leave":
		player, ok := b.players[m.GuildID]
		if !ok { // If we don't have a player in this guild, reply with an error.
			b.sendError(ctx, m.ChannelID, "bot not playing music in this guild")
			return
		}

		// Else, destroy it and leave the voice channel.
		player.Destroy()
		delete(b.players, m.GuildID)

		if err := b.client.LeaveVoiceChannel(context.TODO(), m.GuildID); err != nil {
			fmt.Println("could not properly leave voice channel:", err)
			return
		}
	}
}

// findUser tries to find the given user among the given voice states.
// Returns the voice channel ID the user is in if found, empty string if not.
func findUser(states []voice.State, userID string) string {
	for _, state := range states {
		if state.UserID == userID && state.ChannelID != nil {
			return *state.ChannelID
		}
	}

	return ""
}

// createNewPlayer creates a new player in the voice channel the given userID is.
func (b *bot) createNewPlayer(ctx context.Context, guildID, channelID string) (*Player, error) {
	vc, err := b.client.JoinVoiceChannel(ctx, guildID, channelID, false, true)
	if err != nil {
		return nil, fmt.Errorf("could not join voice channel: %w", err)
	}

	p, err := NewPlayer(vc)
	if err != nil {
		return nil, fmt.Errorf("could not create new player: %w", err)
	}

	return p, nil
}

// sendError is a helper function that replies to the user in the given
// channel with a nicely formatted error.
func (b *bot) sendError(ctx context.Context, chID string, msg string) {
	e := &discord.MessageEmbed{
		Description: ":x: " + msg,
	}

	if _, err := b.client.Channel(chID).Send(ctx, channel.WithMessageEmbed(e)); err != nil {
		fmt.Println("could not send error:", err)
	}
}
