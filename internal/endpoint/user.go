package endpoint

func GetUser(userID string) *Endpoint {
	return &Endpoint{
		URL: "/users/" + userID,
		Key: "/users",
	}
}

func ModifyCurrentUser() *Endpoint {
	return &Endpoint{
		URL: "/users/@me",
		Key: "/users/@me",
	}
}

func UserGuilds() *Endpoint {
	return &Endpoint{
		URL: "/users/@me/guilds",
		Key: "/users/@me/guilds",
	}
}

func LeaveGuild(guildID string) *Endpoint {
	return &Endpoint{
		URL: "/users/@me/guilds/" + guildID,
		Key: "/users/@me/guilds",
	}
}

func GetDMs() *Endpoint {
	return &Endpoint{
		URL: "/users/@me/channels",
		Key: "/users/@me/channels",
	}
}

func CreateDM() *Endpoint {
	return &Endpoint{
		URL: "/users/@me/channels",
		Key: "/users/@me/channels",
	}
}

func GetUserConnections() *Endpoint {
	return &Endpoint{
		URL: "/users/@me/connections",
		Key: "/users/@me/connections",
	}
}
