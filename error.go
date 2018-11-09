package discord

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
)

var (
	// ErrGatewayNotConnected is returned when the client is not connected to the Gateway.
	ErrGatewayNotConnected = errors.New("gateway is not connected")
	// ErrNoFileProvided is returned by SendFiles when no files are provided.
	ErrNoFileProvided = errors.New("no file provided")
)

// APIError is an error returned by the Discord REST API.
type APIError struct {
	HTTPCode int    `json:"http_code"`
	Code     int    `json:"code"`
	Message  string `json:"message"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("%d %s: %s (code %d)",
		e.HTTPCode,
		http.StatusText(e.HTTPCode),
		e.Message,
		e.Code,
	)
}

// apiError is a helper function that extracts an API error from
// an HTTP response and returns it as an APIError.
func apiError(resp *http.Response) error {
	apiErr := APIError{HTTPCode: resp.StatusCode}
	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
		return err
	}
	return apiErr
}

func isConnectionClosed(err error) bool {
	if e, ok := err.(*net.OpError); ok {
		// Ugly but : https://github.com/golang/go/issues/4373
		if e.Err.Error() == "use of closed network connection" {
			return true
		}
	}
	return false
}
