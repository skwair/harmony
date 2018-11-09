package endpoint

func GetGuildEmojis(guildID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + guildID + "/emojis",
		Key: "/channels/" + guildID + "/emojis",
	}
}

func GetGuildEmoji(guildID, emojiID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + guildID + "/emojis/" + emojiID,
		Key: "/channels/" + guildID + "/emojis",
	}
}

func CreateGuildEmoji(guildID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + guildID + "/emojis",
		Key: "/channels/" + guildID + "/emojis",
	}
}

func ModifyGuildEmoji(guildID, emojiID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + guildID + "/emojis/" + emojiID,
		Key: "/channels/" + guildID + "/emojis",
	}
}

func DeleteGuildEmoji(guildID, emojiID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + guildID + "/emojis/" + emojiID,
		Key: "/channels/" + guildID + "/emojis",
	}
}
