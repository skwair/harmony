/*
Package harmony provides an interface to the Discord API
(https://discordapp.com/developers/docs/intro).

Getting started

The first thing you do is to create a Client. For a normal user,
you can get one like this:

	c := harmony.NewClient(harmony.WithToken("userToken"))

If you want to create a client for a bot, use WithBotToken instead
of WithToken, without prefixing the token with "Bot :":

	c := harmony.NewClient(harmony.WithBotToken("botToken"))

You can pass more configuration parameters to NewClient. Review the
documentation of NewClient for more information.

With this client, you can start interacting with the Discord API, but
some methods (such as event handlers) won't be available until you
connect to the Gateway:

	if err = c.Connect(); err != nil {
		// Handle error
	}
	defer c.Disconnect() // Gracefully disconnect

Once connected to the Gateway, you have full access to the Discord API.

Using the HTTP API

Harmony's HTTP API is organized by resource:

	- Guild
	- Channel
	- CurrentUser
	- Webhook
	- Invite

Each resource has its own method on the Client to interact with. For example,
to send a message to a channel:

	msg, err := c.Channel("channel-id").SendMessage("content of the message")
	if err != nil {
		// Handle error
	}
	// msg is the message sent

Registering event handlers

To receive messages, use the OnMessageCreate method and give it your
handler. It will be called each time a message is sent to a channel your
bot is in with the message as a parameter.

	c.OnMessageCreate(func(msg *harmony.Message) {
		fmt.Println(msg.Content)
	})

To register handlers for other types of events, see Client.On* methods.

Using the state

When connecting to Discord, a session state is created with initial data
sent by Discord's Gateway. As events are received by the client, this state
is constantly updated so it always have the newest data available.

This session state acts as a cache to avoid making requests over the HTTP API
each time. If you need to get information about the current user:

	user := c.State.CurrentUser()

Because this state might become memory hungry for bots that are in a very
large number of servers, it can be disabled with the WithStateTracking option
while creating the harmony client.
*/
package harmony

const version = "0.9.0"
