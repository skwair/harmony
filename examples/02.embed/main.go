package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/skwair/harmony"
	"github.com/skwair/harmony/embed"
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

	if err = c.Connect(); err != nil {
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
	// If the new message's content is "!embed",
	// reply with a complex message containing embedded content.
	// This message corresponds to what you can see here :
	// https://leovoel.github.io/embed-visualizer
	if m.Content == "!embed" {
		e := embed.New().
			Title("title ~~(did you know you can have markdown here too?)~~").
			Description("this supports [named links](https://discordapp.com) on top of the previously shown subset of markdown. ```\nyes, even code blocks```").
			URL("https://discordapp.com").
			Color(0x2491ec). // Hexadecimal color code.
			Timestamp(time.Now()).
			Footer(embed.NewFooter().
				Text("footer text").
				Icon("https://cdn.discordapp.com/embed/avatars/0.png").
				Build()).
			Image(embed.NewImage("https://cdn.discordapp.com/embed/avatars/0.png")).
			Author(embed.NewAuthor().
				Name("author name").
				URL("https://discordapp.com").
				IconURL("https://cdn.discordapp.com/embed/avatars/0.png").
				Build()).
			Fields(
				embed.NewField().Name("🤔").Value("some of these properties have certain limits...").Inline(false).Build(),
				embed.NewField().Name("😱").Value("try exceeding some of them!").Inline(false).Build(),
				embed.NewField().Name("🙄").Value("an informative error should show up, and this view will remain as-is until all issues are fixed").Inline(false).Build(),
				embed.NewField().Name("field 1").Value("these last two").Inline(true).Build(),
				embed.NewField().Name("field 2").Value("are inline fields").Inline(true).Build(),
			).
			Build()

		if _, err := b.client.Channel(m.ChannelID).SendEmbed(e); err != nil {
			fmt.Fprintf(os.Stderr, "could not send message: %v", err)
		}
	}
}
