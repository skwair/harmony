package endpoint

import "net/http"

func GetInvite(code, query string) *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/invites/" + code + "?" + query,
		Key:    "/invites",
	}
}

func DeleteInvite(code string) *Endpoint {
	return &Endpoint{
		Method: http.MethodDelete,
		Path:   "/invites/" + code,
		Key:    "/invites",
	}
}
