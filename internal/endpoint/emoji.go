package endpoint

import "net/http"

func ListGuildEmojis(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		URL:    "/channels/" + guildID + "/emojis",
		Key:    "/channels/" + guildID + "/emojis",
	}
}

func GetGuildEmoji(guildID, emojiID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		URL:    "/channels/" + guildID + "/emojis/" + emojiID,
		Key:    "/channels/" + guildID + "/emojis",
	}
}

func CreateGuildEmoji(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPost,
		URL:    "/channels/" + guildID + "/emojis",
		Key:    "/channels/" + guildID + "/emojis",
	}
}

func ModifyGuildEmoji(guildID, emojiID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPatch,
		URL:    "/channels/" + guildID + "/emojis/" + emojiID,
		Key:    "/channels/" + guildID + "/emojis",
	}
}

func DeleteGuildEmoji(guildID, emojiID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodDelete,
		URL:    "/channels/" + guildID + "/emojis/" + emojiID,
		Key:    "/channels/" + guildID + "/emojis",
	}
}
