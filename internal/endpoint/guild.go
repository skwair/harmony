package endpoint

import "net/http"

func CreateGuild() *Endpoint {
	return &Endpoint{
		Method: http.MethodPost,
		Path:   "/guilds",
		Key:    "/guilds",
	}
}

func GetGuild(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/guilds/" + guildID,
		Key:    "/guilds/" + guildID,
	}
}

func ModifyGuild(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPatch,
		Path:   "/guilds/" + guildID,
		Key:    "/guilds/" + guildID,
	}
}

func DeleteGuild(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodDelete,
		Path:   "/guilds/" + guildID,
		Key:    "/guilds/" + guildID,
	}
}

func GetGuildChannels(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/guilds/" + guildID + "/channels",
		Key:    "/guilds/" + guildID + "/channels",
	}
}

func CreateGuildChannel(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPost,
		Path:   "/guilds/" + guildID + "/channels",
		Key:    "/guilds/" + guildID + "/channels",
	}
}

func ModifyChannelPositions(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPatch,
		Path:   "/guilds/" + guildID + "/channels",
		Key:    "/guilds/" + guildID + "/channels",
	}
}

func GetGuildMember(guildID, userID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/guilds/" + guildID + "/members/" + userID,
		Key:    "/guilds/" + guildID + "/members",
	}
}

func ListGuildMembers(guildID, query string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/guilds/" + guildID + "/members?" + query,
		Key:    "/guilds/" + guildID + "/members",
	}
}

func AddGuildMember(guildID, userID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPut,
		Path:   "/guilds/" + guildID + "/members/" + userID,
		Key:    "/guilds/" + guildID + "/members",
	}
}

func RemoveGuildMember(guildID, userID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodDelete,
		Path:   "/guilds/" + guildID + "/members/" + userID,
		Key:    "/guilds/" + guildID + "/members",
	}
}

func ModifyGuildMember(guildID, userID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPatch,
		Path:   "/guilds/" + guildID + "/members/" + userID,
		Key:    "/guilds/" + guildID + "/members",
	}
}

func ModifyCurrentUserNick(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPatch,
		Path:   "/guilds/" + guildID + "/members/@me/nick",
		Key:    "/guilds/" + guildID + "/members/@me/nick",
	}
}

func GetGuildBans(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/guilds/" + guildID + "/bans",
		Key:    "/guilds/" + guildID + "/bans",
	}
}

func CreateGuildBan(guildID, userID, query string) *Endpoint {
	if query != "" {
		query = "?" + query
	}

	return &Endpoint{
		Method: http.MethodPut,
		Path:   "/guilds/" + guildID + "/bans/" + userID + query,
		Key:    "/guilds/" + guildID + "/bans",
	}
}

func RemoveGuildBan(guildID, userID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodDelete,
		Path:   "/guilds/" + guildID + "/bans/" + userID,
		Key:    "/guilds/" + guildID + "/bans",
	}
}

func GetGuildPruneCount(guildID, query string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/guilds/" + guildID + "/prune?" + query,
		Key:    "/guilds/" + guildID + "/prune",
	}
}

func BeginGuildPrune(guildID, query string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPost,
		Path:   "/guilds/" + guildID + "/prune?" + query,
		Key:    "/guilds/" + guildID + "/prune",
	}
}

func GetGuildVoiceRegions(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/guilds/" + guildID + "/regions",
		Key:    "/guilds/" + guildID + "/regions",
	}
}

func GetGuildInvites(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/guilds/" + guildID + "/invites",
		Key:    "/guilds/" + guildID + "/invites",
	}
}

func GetGuildIntegrations(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/guilds/" + guildID + "/integrations",
		Key:    "/guilds/" + guildID + "/integrations",
	}
}

func CreateGuildIntegration(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPost,
		Path:   "/guilds/" + guildID + "/integrations",
		Key:    "/guilds/" + guildID + "/integrations",
	}
}

func ModifyGuildIntegration(guildID, integrationID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPatch,
		Path:   "/guilds/" + guildID + "/integrations/" + integrationID,
		Key:    "/guilds/" + guildID + "/integrations",
	}
}

func DeleteGuildIntegration(guildID, integrationID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodDelete,
		Path:   "/guilds/" + guildID + "/integrations/" + integrationID,
		Key:    "/guilds/" + guildID + "/integrations",
	}
}

func SyncGuildIntegration(guildID, integrationID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPost,
		Path:   "/guilds/" + guildID + "/integrations/" + integrationID + "/sync",
		Key:    "/guilds/" + guildID + "/integrations",
	}
}

func GetGuildEmbed(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/guilds/" + guildID + "/embed",
		Key:    "/guilds/" + guildID + "/embed",
	}
}

func ModifyGuildEmbed(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPatch,
		Path:   "/guilds/" + guildID + "/embed",
		Key:    "/guilds/" + guildID + "/embed",
	}
}

func GetGuildVanityURL(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/guilds/" + guildID + "/vanity-url",
		Key:    "/guilds/" + guildID + "/vanity-url",
	}
}
