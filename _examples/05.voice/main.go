package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/skwair/harmony"
	"github.com/skwair/harmony/discord"
	"github.com/skwair/harmony/log"
	"github.com/skwair/harmony/voice"
)

type bot struct {
	client *harmony.Client

	botID   string
	players map[string]*Player // Player per guild ID.
}

func main() {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "Environment variable BOT_TOKEN must be set.")
		return
	}

	client, err := harmony.NewClient(token, harmony.WithLogger(log.NewStd(os.Stdout, log.LevelDebug)))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	b := &bot{
		client:  client,
		players: make(map[string]*Player),
	}

	// Subscribe to the ready event to know the bot user ID.
	client.OnReady(func(r *harmony.Ready) { b.botID = r.User.ID })
	// Subscribe to new messages so the bot can receive commands.
	client.OnMessageCreate(b.onNewMessage)
	// Also subscribe to voice state updates to know if the bot is kicked
	// or disconnected from a voice channel in order to free allocated resource.
	client.OnVoiceStateUpdate(b.onVoiceUpdate)

	if err = client.Connect(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer client.Disconnect()

	fmt.Println("Bot is running, press ctrl+C to exit.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}

func (b *bot) onNewMessage(m *discord.Message) {
	// Create a base context with a 30 seconds timeout to handle commands.
	// See https://github.com/skwair/harmony/wiki/Event-handlers-and-context
	// for more information.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	b.handleCommand(ctx, m)
}

func (b *bot) onVoiceUpdate(vsu *voice.StateUpdate) {
	// If the bot gets disconnected or kicked by someone, we must clear
	// the player that was associated with it.
	if vsu.UserID == b.botID && vsu.ChannelID == nil {
		if player, ok := b.players[vsu.GuildID]; ok {
			player.Destroy()
			delete(b.players, vsu.GuildID)
		}
	}
}
