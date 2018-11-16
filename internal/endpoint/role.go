package endpoint

func GetGuildRoles(guildID string) *Endpoint {
	return &Endpoint{
		URL: "/guilds/" + guildID + "/roles",
		Key: "/guilds/" + guildID + "/roles",
	}
}

func CreateGuildRole(guildID string) *Endpoint {
	return &Endpoint{
		URL: "/guilds/" + guildID + "/roles",
		Key: "/guilds/" + guildID + "/roles",
	}
}

func ModifyGuildRolePositions(guildID string) *Endpoint {
	return &Endpoint{
		URL: "/guilds/" + guildID + "/roles",
		Key: "/guilds/" + guildID + "/roles",
	}
}

func ModifyGuildRole(guildID, roleID string) *Endpoint {
	return &Endpoint{
		URL: "/guilds/" + guildID + "/roles/" + roleID,
		Key: "/guilds/" + guildID + "/roles",
	}
}

func DeleteGuildRole(guildID, roleID string) *Endpoint {
	return &Endpoint{
		URL: "/guilds/" + guildID + "/roles/" + roleID,
		Key: "/guilds/" + guildID + "/roles",
	}
}

func AddGuildMemberRole(guildID, userID, roleID string) *Endpoint {
	return &Endpoint{
		URL: "/guilds/" + guildID + "/members/" + userID + "/roles/" + roleID,
		Key: "/guilds/" + guildID + "/members",
	}
}

func RemoveGuildMemberRole(guildID, userID, roleID string) *Endpoint {
	return &Endpoint{
		URL: "/guilds/" + guildID + "/members/" + userID + "/roles/" + roleID,
		Key: "/guilds/" + guildID + "/members",
	}
}
