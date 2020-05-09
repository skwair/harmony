# Examples

- 01.pingpong: shows how to create a simple bot that replies with `pong` whenever someone sends a `ping` message.
- 02.embed: demonstrates how to create a bot that replies with some rich embedded content when someone types the `!embed` command.
- 03.files: shows how to send files when someone sends the `!file` command.
- 04.auditlog: shows how to interact with the audit log of a guild.
- 05.voice: a more complex example showcasing how to send voice data with a bot. Available commands: `!play`, `!stop`, `!leave`.

# Creating a Discord bot

## Getting a bot token

In Discord, bots are attached to applications. So you will first need to create one, then attach a bot to this application.

1. Head to [your applications](https://discordapp.com/developers/applications) and create a new one.
2. You will be presented a pop-up where you can set the name for the application as well as the team you want it to be linked to (chose personal if you don't know what it is).
3. Then you will need to go to the `Bot` tab on the right and click the `Add Bot` button to attach a bot to your application.
4. (optional) You can uncheck the `Public Bot` checkbox for now if you do not want anyone to be able to add your bot to their servers.
5. To get the token of your bot, click the `click to reveal` button below to the Token field. Keep it secret, anyone with this token can run *your* bot.

## Adding a bot to a server

1. Go to [your applications](https://discordapp.com/developers/applications) and select the application linked to the bot you want to add to a server.
2. In the `OAuth2` tab on the right, you will find an `OAUTH2 URL GENERATOR` section.
3. Check the `Bot` checkbox under in the `SCOPES` list. A list of permissions will appear below the `SCOPES` list.
4. Select which permissions you want your bot to have when added to a server. Note that those are the permissions your bot will ask, but it doesn't mean your bot users will grant them all. Also, note that if you don't select any permission here, your bot won't be able to do anything.
4. Open the generated link in your browser, and select the server you want to add this bot to. Note that you need the Manage Server permission to be able to add a bot to a server.

# Building and running

To build those examples, simply go to their directory and use `go build`. For example, for `01.pingpong`:

```sh
cd 01.pingpong
go build
```

This will create an executable named after the example directory. To run the examples, you must provide a bot token with the `BOT_TOKEN` environment variable. To prevent the token from leaking in your shell history you can read it into a shell variable that you export so the executable can read it too:

```sh
read BOT_TOKEN
# <paste your bot token here, then press enter>
export BOT_TOKEN
```

Then you can simply run the bot by executing the binary. For the `01.pingpong` example:

```sh
./01.pingpong
2021/07/06 19:01:50 Bot is running, press ctrl+C to exit.
```

# License

Royalty Free Music used in `05.voice` is from Bensound.
