package endpoint

import "net/http"

func GetChannel(chID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/channels/" + chID,
		Key:    "/channels/" + chID,
	}
}

func ModifyChannel(chID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPatch,
		Path:   "/channels/" + chID,
		Key:    "/channels/" + chID,
	}
}

func DeleteChannel(chID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodDelete,
		Path:   "/channels/" + chID,
		Key:    "/channels/" + chID,
	}
}

func GetChannelMessages(chID, query string) *Endpoint {
	if query != "" {
		query = "?" + query
	}

	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/channels/" + chID + "/messages" + query,
		Key:    "/channels/" + chID + "/messages",
	}
}

func GetChannelMessage(chID, msgID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/channels/" + chID + "/messages/" + msgID,
		Key:    "/channels/" + chID + "/messages",
	}
}

func EditChannelPermissions(chID, targetID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPut,
		Path:   "/channels/" + chID + "/permissions/" + targetID,
		Key:    "/channels/" + chID + "/permissions",
	}
}

func DeleteChannelPermission(chID, targetID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodDelete,
		Path:   "/channels/" + chID + "/permissions/" + targetID,
		Key:    "/channels/" + chID + "/permissions",
	}
}

func GetChannelInvites(chID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/channels/" + chID + "/invites",
		Key:    "/channels/" + chID + "/invites",
	}
}

func CreateChannelInvite(chID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPost,
		Path:   "/channels/" + chID + "/invites",
		Key:    "/channels/" + chID + "/invites",
	}
}

func GroupDMAddRecipient(chID, recipientID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPut,
		Path:   "/channels/" + chID + "/recipients/" + recipientID,
		Key:    "/channels/" + chID + "/recipients",
	}
}

func GroupDMRemoveRecipient(chID, recipientID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodDelete,
		Path:   "/channels/" + chID + "/recipients/" + recipientID,
		Key:    "/channels/" + chID + "/recipients",
	}
}

func TriggerTypingIndicator(chID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPost,
		Path:   "/channels/" + chID + "/typing",
		Key:    "/channels/" + chID + "/typing",
	}
}
