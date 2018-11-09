package endpoint

func GetChannel(chID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + chID,
		Key: "/channels/" + chID,
	}
}

func ModifyChannel(chID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + chID,
		Key: "/channels/" + chID,
	}
}

func DeleteChannel(chID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + chID,
		Key: "/channels/" + chID,
	}
}

func GetChannelMessages(chID, query string) *Endpoint {
	if query != "" {
		query = "?" + query
	}

	return &Endpoint{
		URL: "/channels/" + chID + "/messages" + query,
		Key: "/channels/" + chID + "/messages",
	}
}

func GetChannelMessage(chID, msgID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + chID + "/messages/" + msgID,
		Key: "/channels/" + chID + "/messages",
	}
}

func UpdateChannelPermissions(chID, targetID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + chID + "/permissions/" + targetID,
		Key: "/channels/" + chID + "/permissions",
	}
}

func DeleteChannelPermission(chID, targetID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + chID + "/permissions/" + targetID,
		Key: "/channels/" + chID + "/permissions",
	}
}

func GetChannelInvites(chID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + chID + "/invites",
		Key: "/channels/" + chID + "/invites",
	}
}

func CreateChannelInvite(chID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + chID + "/invites",
		Key: "/channels/" + chID + "/invites",
	}
}

func AddRecipient(chID, recipientID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + chID + "/recipients/" + recipientID,
		Key: "/channels/" + chID + "/recipients",
	}
}

func RemoveRecipient(chID, recipientID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + chID + "/recipients/" + recipientID,
		Key: "/channels/" + chID + "/recipients",
	}
}

func TriggerTyping(chID string) *Endpoint {
	return &Endpoint{
		URL: "/channels/" + chID + "/typing",
		Key: "/channels/" + chID + "/typing",
	}
}
