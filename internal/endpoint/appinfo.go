package endpoint

import "net/http"

func GetApplicationInfo() *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		URL:    "/oauth2/applications/@me",
		Key:    "/oauth2/applications/@me",
	}
}
