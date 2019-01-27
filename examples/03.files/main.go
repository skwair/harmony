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
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		fmt.Fprint(os.Stderr, "Environment variable BOT_TOKEN must be set.")
		return
	}

	c, err := harmony.NewClient(harmony.WithBotToken(botToken))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}

	b := &bot{client: c}

	c.OnMessageCreate(b.onNewMessage)

	if err = c.Connect(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}
	defer c.Disconnect()

	log.Println("Bot is running, press ctrl+C to exit.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}

func (b *bot) onNewMessage(m *harmony.Message) {
	if m.Content == "!file" {
		f, err := os.Open("discord-gopher.png")
		if err != nil {
			log.Println(err)
			return
		}

		file := harmony.File{
			// Setting a valid extension type here,
			// such as "png", will allow Discord
			// applications to display the files
			// inline, instead of asking users
			// to download them.
			Name:   "discord-gopher.png",
			Reader: f,
		}

		if _, err = b.client.Channel(m.ChannelID).Send(context.Background(), harmony.WithFiles(file)); err != nil {
			log.Println(err)
		}
	}
}
