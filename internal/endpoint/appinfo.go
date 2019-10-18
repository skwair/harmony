package endpoint

import "net/http"

func GetAppInfo() *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		URL:    "/oauth2/applications/@me",
		Key:    "/oauth2/applications/@me",
	}
}
