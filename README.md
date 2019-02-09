[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/skwair/harmony)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat-square)](LICENSE)
[![Discord](https://img.shields.io/badge/Discord-online-7289DA.svg?style=flat-square)](https://discord.gg/3sVFWQC)
[![Build Status](https://travis-ci.org/skwair/harmony.svg?branch=master)](https://travis-ci.org/skwair/harmony)


# Harmony

<img align="right" height="200" src=".github/discord-gopher.png">

Harmony is a peaceful [Go](https://golang.org) module for interacting with [Discord](http://discordapp.com)'s API.

Although this package is usable, it still is under active development so please don't use it for anything other than experiments, yet.

**Contents**

- [Installation](#installation)
- [Usage](#usage)
- [Testing](#testing)
- [How does it compare to DiscordGo?](#how-does-it-compare-to-discordgo-)
- [License](#license)

# Installation

Make sure you have a working Go installation, if not see [this page](https://golang.org/dl) first.

Then, install this package with the `go get` command:

```sh
go get -u github.com/skwair/harmony
```

> Note that `go get -u` will always pull the latest version from the master branch before Go 1.11. With newer versions and Go modules enabled, the latest minor or patch release will be downloaded. `go get github.com/skwair/harmony@major.minor.patch` can be used to download a specific version. See [Go modules](https://github.com/golang/go/wiki/Modules#how-to-upgrade-and-downgrade-dependencies) for more information.

# Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/skwair/harmony"
)

func main() {
    c, err := harmony.NewClient("your.bot.token")
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

For information about how to create bots and more examples on how to use this package, check out the [examples](https://github.com/skwair/harmony/blob/master/examples) directory and the [tests](https://github.com/skwair/harmony/blob/master/harmony_test.go).

# Testing

For now, only some end to end tests are provided with this module. To run them, you will need a valid bot token and a valid Discord server ID. The bot attached to the token must be in the server with administrator permissions.

1. Create a Discord test server

From a Discord client and with you main account, simply create a new server. Then, right click on the new server and get it's ID.

> Note that for the UI to have the `Copy ID` option when right clicking on the server, you will need to enable developer mode. You can find this option in `User settings > Appearance > Advanced > Developer Mode`.

2. Create a bot and add it to the test Discord server

Create a bot (or use an existing one) and add it to the freshly created server.

> See the [example directory](https://github.com/skwair/harmony/blob/master/examples) for information on how to create a bot and add it to a server.

3. Set required environment variables and run the tests

Set `HARMONY_TEST_BOT_TOKEN` to the token of your bot and `HARMONY_TEST_GUILD_ID` to the ID of the server you created and simply run:

⚠️ **For the tests to be reproducible, they will start by deleting ALL channels in the provided server. Please make sure to provide a server created ONLY for those tests.** ⚠️

```bash
go test -v -race ./...
```

> Step 1 and 2 must be done only once for initial setup. Once you have your bot token and the ID of your test server, you can run the tests as many times as you want.

# How does it compare to [DiscordGo](https://github.com/bwmarrin/discordgo)?

DiscordGo offers some additional features right now, such as a way to create and manage your [Discord applications](https://discordapp.com/developers/applications/me). The majority of features though, such a receiving events, sending messages, receiving and sending voice data, etc. are also implemented in this library. Another thing that this library does not support is self bot, as they are a violation of Discord's TOS.

The main difference resides in the Gateway (websocket) real time API implementation. This library takes a different approach (using more synchronisation mechanisms) to avoid having to rely on hacks like [this](https://github.com/bwmarrin/discordgo/blob/8325a6bf6dd6c91ed4040a1617b07287b8fb0eba/wsapi.go#L868) or [this](https://github.com/bwmarrin/discordgo/blob/8325a6bf6dd6c91ed4040a1617b07287b8fb0eba/wsapi.go#L822), hopefully providing a more robust implementation as well as a better user experience.

Another difference is in the "event handler" mechanism. Instead of having a single [method](https://github.com/bwmarrin/discordgo/blob/8325a6bf6dd6c91ed4040a1617b07287b8fb0eba/event.go#L120) that takes an `interface{}` as a parameter and guesses for which event you registered a handler based on its concrete type, this library provides one method per event type, making it clear what signature your handler must have and ensuring it at compile time, not at runtime.

Finally, this library has a full support of the [context](https://golang.org/pkg/context/) package, allowing the use of timeouts, deadlines and cancellation when interacting with Discord's API.

# License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/skwair/harmony/blob/master/LICENSE) file for details.

Original logo by [Renee French](https://instagram.com/reneefrench), dressed with the cool t-shirt by [@HlneChd](https://twitter.com/hlnechd).
