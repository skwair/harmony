package endpoint

import "net/http"

func GetVoiceRegions() *Endpoint {
	return &Endpoint{
		Method: http.MethodGet,
		URL:    "/voice/regions",
		Key:    "/voice/regions",
	}
}
