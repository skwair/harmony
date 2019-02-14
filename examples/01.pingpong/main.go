package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/skwair/harmony"
)

// Creating a struct that will hold your bot's dependencies.
// For this simple example, there is just the harmony client,
// but for a more complex bot, you could have a logger, a
// database, etc.
//
// This is not mandatory but it's a good way to keep your
// code clean and readable.
type bot struct {
	client *harmony.Client
}

func main() {
	// Fetch the bot token from env.
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		// Not using log.Fatal() or anything that calls os.Exit()
		// because defers are not run, thus we won't disconnect
		// properly from the Gateway.
		fmt.Fprint(os.Stderr, "Environment variable BOT_TOKEN must be set.")
		return
	}

	// Create a harmony client with a bot token.
	// NewClient automatically prepends your bot token with "Bot ",
	// which is a requirement by Discord for bot users.
	c, err := harmony.NewClient(token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}

	b := &bot{client: c}

	// Register a callback for MESSAGE_CREATE events.
	// Note that we won't receive events before the client
	// is actually connected to the Gateway.
	c.OnMessageCreate(b.onNewMessage)

	// Connect to the Gateway. From now on, our registered
	// handler for MESSAGE_CREATE will be called when there
	// are new messages.
	// This connection is designed to be long lived and to survive
	// network failures, attempting to reconnect whenever a problem occurs.
	if err = c.Connect(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}
	defer c.Disconnect()

	log.Println("Bot is running, press ctrl+C to exit.")

	// Wait for ctrl-C, then exit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}

// By declaring your handlers as methods of the bot struct, they
// have access to your bot's dependencies (here, the harmony client).
func (b *bot) onNewMessage(m *harmony.Message) {
	// If the new message's content is "ping",
	// Reply with "pong", logging any error
	// that occurs.
	msg := withReplier(m, b.client)

	if msg.Content == "ping" {
		if _, err := msg.Reply(context.Background(), harmony.WithContent("pong")); err != nil {
			log.Println(err)
		}
		// if _, err := b.client.Channel(msg.ChannelID).SendMessage(context.Background(), "pong"); err != nil {
		// 	log.Println(err)
		// }
	}
}

type message struct {
	*harmony.Message
	replyer
}

type replyer struct {
	client    *harmony.Client
	channelID string
}

func withReplier(msg *harmony.Message, client *harmony.Client) *message {
	return &message{msg, replyer{client: client, channelID: msg.ChannelID}}
}

func (r *replyer) Reply(ctx context.Context, opts ...harmony.MessageOption) (*harmony.Message, error) {
	return r.client.Channel(r.channelID).Send(ctx, opts...)
}
