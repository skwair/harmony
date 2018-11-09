/*
Package discord provides an interface to the Discord API
(https://discordapp.com/developers/docs/intro).

The first thing you do is to create a Client. For a normal user,
you can get one like this :

	c := discord.NewClient(discord.WithToken("userToken"))

If you want to create a client for a bot, use WithBotToken instead
of WithToken, without prefixing the token with "Bot :" :

	c := discord.NewClient(discord.WithBotToken("botToken"))

You can pass more configuration parameters to NewClient. Review the
documentation of NewClient for more information.

With this client, you can start interacting with the Discord API, but
some methods (such as event handlers) won't be available until you
connect to the Gateway :

	if err = c.Connect(); err != nil {
		// Handle error
	}
	defer c.Disconnect()

Once connected to the Gateway, you have full access to the Discord API.
*/
package discord
