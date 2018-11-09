package endpoint

func CreateMessage(chID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + chID + "/messages",
		Key: "/channels/" + chID + "/messages",
	}
}

func EditMessage(chID, msgID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + chID + "/messages/" + msgID,
		Key: "/channels/" + chID + "/messages",
	}
}

func DeleteMessage(chID, msgID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + chID + "/messages/" + msgID,
		// Deleting messages falls under a separate, higher rate limit.
		// This is why the HTTP verb is present in this key.
		Key: "DELETE /channels/" + chID + "/messages",
	}
}

func BulkDeleteMessage(chID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + chID + "/messages/bulk-delete",
		Key: "/channels/" + chID + "/messages/bulk-delete",
	}
}

func GetPins(chID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + chID + "/pins",
		Key: "/channels/" + chID + "/pins",
	}
}

func PinMessage(chID, msgID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + chID + "/pins/" + msgID,
		Key: "/channels/" + chID + "/pins",
	}
}

func UnpinMessage(chID, msgID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + chID + "/pins/" + msgID,
		Key: "/channels/" + chID + "/pins",
	}
}
