package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/skwair/discord"
)

// Creating a struct that will hold your bot's dependencies.
// For this simple example, there is just the Discord client,
// but for a more complex bot, you could have a logger, a
// database, etc.
//
// This is not mandatory but it's a good way to keep your
// code clean and readable.
type bot struct {
	client *discord.Client
}

func main() {
	// Fetch the bot token from env.
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		// Not using log.Fatal() or anything that calls os.Exit()
		// because defers are not run, thus we won't disconnect
		// properly from the Gateway.
		fmt.Fprint(os.Stderr, "Environment variable BOT_TOKEN must be set.")
		return
	}

	// Create a discord client with a bot token.
	// Using WithBotToken will automatically prepend your bot token
	// with "Bot ", which is a requirement by Discord for bot users.
	c, err := discord.NewClient(discord.WithBotToken(botToken))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}

	b := &bot{client: c}

	// Register a callback for MESSAGE_CREATE events.
	// Note that we won't receive events before the client
	// is actually connected to the Gateway.
	c.HandleMessageCreate(b.onNewMessage)

	// Connect to the Gateway. From now on, our registered
	// handler for MESSAGE_CREATE will be called when there
	// are new messages.
	// This connection is designed to be long lived and to survive
	// network failures, attempting to reconnect whenever a problem occurs.
	if err = c.Connect(); err != nil {
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
// have access to your bot's dependencies (here, the Discord client).
func (b *bot) onNewMessage(m *discord.Message) {
	// If the new message's content is "ping",
	// Reply with "pong", logging any error
	// that occurs.
	if m.Content == "ping" {
		if _, err := b.client.Channel(m.ChannelID).SendMessage("pong"); err != nil {
			log.Println(err)
		}
	}
}
