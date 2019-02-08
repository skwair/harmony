package endpoint

import "net/http"

func CreateReaction(chID, msgID, emoji string) *Endpoint {
	return &Endpoint{
		Method: http.MethodPut,
		URL:    "/channels/" + chID + "/messages/" + msgID + "/reactions/" + emoji + "/@me",
		Key:    "/channels/" + chID + "/messages",
	}
}

func DeleteOwnReaction(chID, msgID, emoji string) *Endpoint {
	return &Endpoint{
		Method: http.MethodDelete,
		URL:    "/channels/" + chID + "/messages/" + msgID + "/reactions/" + emoji + "/@me",
		Key:    "/channels/" + chID + "/messages",
	}
}

func DeleteUserReaction(chID, msgID, userID, emoji string) *Endpoint {
	return &Endpoint{
		Method: http.MethodDelete,
		URL:    "/channels/" + chID + "/messages/" + msgID + "/reactions/" + emoji + "/" + userID,
		Key:    "/channels/" + chID + "/messages",
	}
}

func DeleteAllReactions(chID, msgID string) *Endpoint {
	return &Endpoint{
		Method: http.MethodDelete,
		URL:    "/channels/" + chID + "/messages/" + msgID + "/reactions",
		Key:    "/channels/" + chID + "/messages",
	}
}

func GetReactions(chID, msgID, emoji, query string) *Endpoint {
	if query != "" {
		query = "?" + query
	}

	return &Endpoint{
		Method: http.MethodGet,
		URL:    "/channels/" + chID + "/messages/" + msgID + "/reactions/" + emoji + query,
		Key:    "/channels/" + chID + "/messages",
	}
}
