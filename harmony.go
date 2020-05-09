/*
Package harmony provides an interface to the Discord API
(https://discordapp.com/developers/docs/intro).

Getting started

The first thing you do is to create a Client. NewClient returns a new Client
configured with sain defaults which should work just fine in most cases.
However, should you need a more specific configuration, you can always
tweak it with optional `ClientOption`s. See the documentation of NewClient
and the ClientOption type for more information on how to do so.

	client := harmony.NewClient("your.bot.token")

Once you have a Client, you can start interacting with the Discord API,
but some methods (such as event handlers) won't be available until you
connect to Discord's Gateway. You can do so by simply calling the Connect
method of the Client:

	if err = client.Connect(); err != nil {
		// Handle error
	}
	defer client.Disconnect() // Gracefully disconnect

It is only when successfully connected to the Gateway that your bot will
appear as online and your Client will be able to receive events and send
messages.

Using the HTTP API

Harmony's HTTP API is organized by resource. A resource maps to a core
concept in the Discord world, such as a User or a Channel. Here is the
list of resources you can interact with:

	- Guild
	- Channel
	- CurrentUser
	- Webhook
	- Invite

Every interaction you can have with a resource can be accessed via
methods attached to it. For example, if you wish to send a message
to a channel:

	msg, err := client.Channel("channel-id").SendMessage("content of the message")
	if err != nil {
		// Handle error
	}
	// msg is the message sent

Registering event handlers

To receive messages, use the OnMessageCreate method and give it your
handler. It will be called each time a message is sent to a channel your
bot is in with the message as a parameter.

	client.OnMessageCreate(func(msg *harmony.Message) {
		fmt.Println(msg.Content)
	})

To register handlers for other types of events, see Client.On* methods.

Note that your handlers are called in their own goroutine, meaning
whatever you do inside of them won't block future events.

Using the state

When connecting to Discord, a session state is created with initial data
sent by Discord's Gateway. As events are received by the client, this state
is constantly updated so it always have the newest data available.

This session state acts as a cache to avoid making requests over the HTTP API
each time. If you need to get information about the current user:

	user := client.State.CurrentUser()

Because this state might become memory hungry for bots that are in a very
large number of servers, it can be disabled with the WithStateTracking option
while creating the harmony client.
*/
package harmony

const version = "0.16.0"
