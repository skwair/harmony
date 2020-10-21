package endpoint

import "net/http"

func GetApplication() *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/oauth2/applications/@me",
		Key:    "/oauth2/applications/@me",
	}
}
