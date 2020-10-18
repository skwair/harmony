package endpoint

import "net/http"

func Gateway() *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/gateway",
		Key:    "/gateway",
	}
}

func GatewayBot() *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		Path:   "/gateway/bot",
		Key:    "/gateway/bot",
	}
}
