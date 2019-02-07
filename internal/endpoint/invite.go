package endpoint

import "net/http"

func GetInvite(code, query string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		URL:    "/invites/" + code + "?" + query,
		Key:    "/invites",
	}
}

func DeleteInvite(code string) *Endpoint {
	return &Endpoint{
		Method: http.MethodDelete,
		URL:    "/invites/" + code,
		Key:    "/invites",
	}
}
