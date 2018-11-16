[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/skwair/discord)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat-square)](LICENSE)
[![Discord](https://img.shields.io/badge/Discord-online-7289DA.svg?style=flat-square)](https://discord.gg/3sVFWQC)


# Discord

<img align="right" height="200" src=".github/discord-gopher.png">

An unofficial client for the [Discord](http://discordapp.com) API written in [Go](https://golang.org).

Although this package is usable, it still is under active development so please don't use it for anything other than experiments, yet.

**Contents :**

- [Installation](#installation)
- [Usage](#usage)
- [How does it compare to DiscordGo ?](#how-does-it-compare-to-discordgo-)
- [License](#license)

# Installation

Make sure you have a working Go installation, if not see [this page](https://golang.org/dl) first.

Then, install this package with the `go get` command :

```sh
go get -u github.com/skwair/discord
```

> Note that `go get -u` will always pull the latest version from the master branch before Go 1.11. With newer versions and Go modules enabled, the latest minor or patch release will be downloaded. `go get github.com/skwair/discord@major.minor.patch` can be used to download a specific version. See [Go modules](https://github.com/golang/go/wiki/Modules#how-to-upgrade-and-downgrade-dependencies) for more information.

# Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/skwair/discord"
)

func main() {
    c, err := discord.NewClient(discord.WithBotToken("your.bot.token"))
    if err != nil {
        log.Fatal(err)
    }

    // Get information about the current user.
    u, err := c.CurrentUser().Get()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(u)
}
```

For information about how to create bots and more examples on how to use this package, check out the [examples](https://github.com/skwair/discord/blob/master/examples) directory.

# How does it compare to [DiscordGo](https://github.com/bwmarrin/discordgo) ?

DiscordGo offers some additional features right now, such as a way to create and manage your [Discord applications](https://discordapp.com/developers/applications/me). The majority of features though, such a receiving events, sending messages, receiving and sending voice data, etc. are also implemented in this library.

The main difference resides in the Gateway (websocket) real time API implementation. This library takes a different approach (using more synchronisation mechanisms) to avoid having to rely on hacks like [this](https://github.com/bwmarrin/discordgo/blob/8325a6bf6dd6c91ed4040a1617b07287b8fb0eba/wsapi.go#L868) or [this](https://github.com/bwmarrin/discordgo/blob/8325a6bf6dd6c91ed4040a1617b07287b8fb0eba/wsapi.go#L822), hopefully providing a more robust implementation as well as a better user experience.

Another difference is in the "event handler" mechanism. Instead of having a single [method](https://github.com/bwmarrin/discordgo/blob/73f6772a2b7cc95e29c462e4f15bf07cbe0d3854/event.go#L111) that takes an `interface{}` as a parameter and guesses for which event you registered a handler based on its concrete type, this library provides one method per event type, making it clear what signature your handler must have and ensuring it at compile time, not at runtime.

# License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/skwair/discord/blob/master/LICENSE) file for details.

Original logo by [Renee French](https://instagram.com/reneefrench), dressed with the cool t-shirt by [@HlneChd](https://twitter.com/hlnechd).
