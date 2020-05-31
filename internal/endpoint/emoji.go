package endpoint

import "net/http"

func ListGuildEmojis(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/channels/" + guildID + "/emojis",
		Key:    "/channels/" + guildID + "/emojis",
	}
}

func GetGuildEmoji(guildID, emojiID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/channels/" + guildID + "/emojis/" + emojiID,
		Key:    "/channels/" + guildID + "/emojis",
	}
}

func CreateGuildEmoji(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPost,
		Path:   "/channels/" + guildID + "/emojis",
		Key:    "/channels/" + guildID + "/emojis",
	}
}

func ModifyGuildEmoji(guildID, emojiID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPatch,
		Path:   "/channels/" + guildID + "/emojis/" + emojiID,
		Key:    "/channels/" + guildID + "/emojis",
	}
}

func DeleteGuildEmoji(guildID, emojiID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodDelete,
		Path:   "/channels/" + guildID + "/emojis/" + emojiID,
		Key:    "/channels/" + guildID + "/emojis",
	}
}
