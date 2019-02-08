package endpoint

import "net/http"

func GetUser(userID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		URL:    "/users/" + userID,
		Key:    "/users",
	}
}

func ModifyCurrentUser() *Endpoint {
	return &Endpoint{
		Method: http.MethodPatch,
		URL:    "/users/@me",
		Key:    "/users/@me",
	}
}

func GetCurrentUserGuilds() *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		URL:    "/users/@me/guilds",
		Key:    "/users/@me/guilds",
	}
}

func LeaveGuild(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodDelete,
		URL:    "/users/@me/guilds/" + guildID,
		Key:    "/users/@me/guilds",
	}
}

func GetUserDMs() *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		URL:    "/users/@me/channels",
		Key:    "/users/@me/channels",
	}
}

func CreateDM() *Endpoint {
	return &Endpoint{
		Method: http.MethodPost,
		URL:    "/users/@me/channels",
		Key:    "/users/@me/channels",
	}
}

func GetUserConnections() *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		URL:    "/users/@me/connections",
		Key:    "/users/@me/connections",
	}
}
