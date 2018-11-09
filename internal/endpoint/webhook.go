package endpoint

func CreateWebhook(chID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + chID + "/webhooks",
		Key: "/channels/" + chID + "/webhooks",
	}
}

func GetChannelWebhooks(chID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + chID + "/webhooks",
		Key: "/channels/" + chID + "/webhooks",
	}
}

func GetGuildWebhooks(guildID string) *Endpoint {
	return &Endpoint{
		URL: "/guilds/" + guildID + "/webhooks",
		Key: "/guilds/" + guildID + "/webhooks",
	}
}

func GetWebhook(whID string) *Endpoint {
	return &Endpoint{
		URL: "/webhooks/" + whID,
		Key: "/webhooks/" + whID,
	}
}

func GetWebhookWithToken(whID, token string) *Endpoint {
	return &Endpoint{
		URL: "/webhooks/" + whID + "/" + token,
		Key: "/webhooks/" + whID,
	}
}

func ModifyWebhook(whID string) *Endpoint {
	return &Endpoint{
		URL: "/webhooks/" + whID,
		Key: "/webhooks/" + whID,
	}
}

func ModifyWebhookWithToken(whID, token string) *Endpoint {
	return &Endpoint{
		URL: "/webhooks/" + whID + "/" + token,
		Key: "/webhooks/" + whID,
	}
}

func DeleteWebhook(whID string) *Endpoint {
	return &Endpoint{
		URL: "/webhooks/" + whID,
		Key: "/webhooks/" + whID,
	}
}

func DeleteWebhookWithToken(whID, token string) *Endpoint {
	return &Endpoint{
		URL: "/webhooks/" + whID + "/" + token,
		Key: "/webhooks/" + whID,
	}
}

func ExecuteWebhook(whID, token, query string) *Endpoint {
	return &Endpoint{
		URL: "/webhooks/" + whID + "/" + token + "?" + query,
		Key: "webhooks/" + whID,
	}
}
