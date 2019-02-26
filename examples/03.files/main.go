package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/skwair/harmony"
)

// NOTE: the structure of this bot is detailed in the
// first example : 01.pingpong.

type bot struct {
	client *harmony.Client
}

func main() {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		fmt.Fprint(os.Stderr, "Environment variable BOT_TOKEN must be set.")
		return
	}

	client, err := harmony.NewClient(token)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}

	b := &bot{client: client}

	client.OnMessageCreate(b.onNewMessage)

	if err = client.Connect(context.Background()); err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	defer client.Disconnect()

	log.Println("Bot is running, press ctrl+C to exit.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}

func (b *bot) onNewMessage(m *harmony.Message) {
	if m.Content == "!file" {
		// Here, we demonstrate the FileFromDisk function to send a file present
		// on the local filesystem. If the resource you want to send is a URL,
		// use FileFromURL instead.
		// If you already have your own reader, then FileFromReadCloser is the
		// function you want to use.
		file, err := harmony.FileFromDisk("discord-gopher.png", "zob")
		if err != nil {
			log.Println(err)
			return
		}

		if _, err = b.client.Channel(m.ChannelID).Send(context.Background(), harmony.WithFiles(file)); err != nil {
			log.Println(err)
		}
	}
}
