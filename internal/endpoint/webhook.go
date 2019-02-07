package endpoint

import "net/http"

func CreateWebhook(chID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPost,
		URL:    "/channels/" + chID + "/webhooks",
		Key:    "/channels/" + chID + "/webhooks",
	}
}

func GetChannelWebhooks(chID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		URL:    "/channels/" + chID + "/webhooks",
		Key:    "/channels/" + chID + "/webhooks",
	}
}

func GetGuildWebhooks(guildID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		URL:    "/guilds/" + guildID + "/webhooks",
		Key:    "/guilds/" + guildID + "/webhooks",
	}
}

func GetWebhook(whID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		URL:    "/webhooks/" + whID,
		Key:    "/webhooks/" + whID,
	}
}

func GetWebhookWithToken(whID, token string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		URL:    "/webhooks/" + whID + "/" + token,
		Key:    "/webhooks/" + whID,
	}
}

func ModifyWebhook(whID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPatch,
		URL:    "/webhooks/" + whID,
		Key:    "/webhooks/" + whID,
	}
}

func ModifyWebhookWithToken(whID, token string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPatch,
		URL:    "/webhooks/" + whID + "/" + token,
		Key:    "/webhooks/" + whID,
	}
}

func DeleteWebhook(whID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodDelete,
		URL:    "/webhooks/" + whID,
		Key:    "/webhooks/" + whID,
	}
}

func DeleteWebhookWithToken(whID, token string) *Endpoint {
	return &Endpoint{
		Method: http.MethodDelete,
		URL:    "/webhooks/" + whID + "/" + token,
		Key:    "/webhooks/" + whID,
	}
}

func ExecuteWebhook(whID, token, query string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPost,
		URL:    "/webhooks/" + whID + "/" + token + "?" + query,
		Key:    "webhooks/" + whID,
	}
}
