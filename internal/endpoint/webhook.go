package endpoint

import "net/http"

func CreateWebhook(chID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPost,
		Path:   "/channels/" + chID + "/webhooks",
		Key:    "/channels/" + chID + "/webhooks",
	}
}

func GetChannelWebhooks(chID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/channels/" + chID + "/webhooks",
		Key:    "/channels/" + chID + "/webhooks",
	}
}

func GetGuildWebhooks(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/guilds/" + guildID + "/webhooks",
		Key:    "/guilds/" + guildID + "/webhooks",
	}
}

func GetWebhook(whID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/webhooks/" + whID,
		Key:    "/webhooks/" + whID,
	}
}

func GetWebhookWithToken(whID, token string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/webhooks/" + whID + "/" + token,
		Key:    "/webhooks/" + whID,
	}
}

func ModifyWebhook(whID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPatch,
		Path:   "/webhooks/" + whID,
		Key:    "/webhooks/" + whID,
	}
}

func ModifyWebhookWithToken(whID, token string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPatch,
		Path:   "/webhooks/" + whID + "/" + token,
		Key:    "/webhooks/" + whID,
	}
}

func DeleteWebhook(whID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodDelete,
		Path:   "/webhooks/" + whID,
		Key:    "/webhooks/" + whID,
	}
}

func DeleteWebhookWithToken(whID, token string) *Endpoint {
	return &Endpoint{
		Method: http.MethodDelete,
		Path:   "/webhooks/" + whID + "/" + token,
		Key:    "/webhooks/" + whID,
	}
}

func ExecuteWebhook(whID, token, query string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPost,
		Path:   "/webhooks/" + whID + "/" + token + "?" + query,
		Key:    "/webhooks/" + whID,
	}
}
