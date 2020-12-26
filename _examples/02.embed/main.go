package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/skwair/harmony"
	"github.com/skwair/harmony/discord"
	"github.com/skwair/harmony/resource/channel"
)

// NOTE: the structure of this bot is detailed in the
// first example: 01.pingpong.

type bot struct {
	client *harmony.Client
}

func main() {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "Environment variable BOT_TOKEN must be set.")
		return
	}

	client, err := harmony.NewClient(token)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	b := &bot{client: client}

	client.OnMessageCreate(b.onNewMessage)

	if err = client.Connect(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer client.Disconnect()

	log.Println("Bot is running, press ctrl+C to exit.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}

func (b *bot) onNewMessage(m *discord.Message) {
	// If the new message's content is "!embed",
	// reply with a complex message containing embedded content.
	// This message corresponds to what you can see here:
	// https://leovoel.github.io/embed-visualizer
	if m.Content == "!embed" {
		e := &discord.MessageEmbed{
			Title:       "title ~~(did you know you can have markdown here too?)~~",
			Type:        "rich",
			Description: "this supports [named links](https://discord.com) on top of the previously shown subset of markdown. ```\nyes, even code blocks```",
			URL:         "https://discord.com",
			Timestamp:   time.Now(),
			Color:       0x2491ec,
			Footer: &discord.MessageEmbedFooter{
				Text:    "footer text",
				IconURL: "https://cdn.discordapp.com/embed/avatars/0.png",
			},
			Image: &discord.MessageEmbedImage{
				URL: "https://cdn.discordapp.com/embed/avatars/0.png",
			},
			Author: &discord.MessageEmbedAuthor{
				Name:    "author name",
				URL:     "https://discord.com",
				IconURL: "https://cdn.discordapp.com/embed/avatars/0.png",
			},
			Fields: []discord.MessageEmbedField{
				{
					Name:   "🤔",
					Value:  "some of these properties have certain limits...",
					Inline: false,
				},
				{
					Name:   "😱",
					Value:  "try exceeding some of them!",
					Inline: false,
				},
				{
					Name:   "🙄",
					Value:  "an informative error should show up, and this view will remain as-is until all issues are fixed",
					Inline: false,
				},
				{
					Name:   "field 1",
					Value:  "these last two",
					Inline: true,
				},
				{
					Name:   "field 2",
					Value:  "are inline fields",
					Inline: true,
				},
			},
		}

		if _, err := b.client.Channel(m.ChannelID).Send(context.Background(), channel.WithMessageEmbed(e)); err != nil {
			log.Println(err)
		}
	}
}
