package rest

import "encoding/json"

// Payload is a payload that is sent to Discord's REST API.
type Payload struct {
	body        []byte
	contentType string
}

func (p *Payload) hasBody() bool {
	return p != nil && p.body != nil
}

// JSONPayload creates a new Payload from some raw JSON data.
func JSONPayload(body json.RawMessage) *Payload {
	return &Payload{
		body:        body,
		contentType: "application/json",
	}
}

// CustomPayload creates a new custom payload from raw bytes and a given content type.
func CustomPayload(body []byte, contentType string) *Payload {
	return &Payload{
		body:        body,
		contentType: contentType,
	}
}
