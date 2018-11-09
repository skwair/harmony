package endpoint

func GetRoles(guildID string) *Endpoint {
	return &Endpoint{
		URL: "/guilds/" + guildID + "/roles",
		Key: "/guilds/" + guildID + "/roles",
	}
}

func CreateRole(guildID string) *Endpoint {
	return &Endpoint{
		URL: "/guilds/" + guildID + "/roles",
		Key: "/guilds/" + guildID + "/roles",
	}
}

func ModifyRolePositions(guildID string) *Endpoint {
	return &Endpoint{
		URL: "/guilds/" + guildID + "/roles",
		Key: "/guilds/" + guildID + "/roles",
	}
}

func ModifyRole(guildID, roleID string) *Endpoint {
	return &Endpoint{
		URL: "/guilds/" + guildID + "/roles/" + roleID,
		Key: "/guilds/" + guildID + "/roles",
	}
}

func DeleteRole(guildID, roleID string) *Endpoint {
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
