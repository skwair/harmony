package discord

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	// ErrGatewayNotConnected is returned when the client is not connected to the Gateway.
	ErrGatewayNotConnected = errors.New("gateway is not connected")
	// ErrGatewayAlreadyConnected is returned by Connect when a connection to the Gateway already exists.
	ErrGatewayAlreadyConnected = errors.New("already connected to the Gateway")
	// ErrInvalidMessageSend is returned by Send when no files are provided.
	ErrInvalidMessageSend = errors.New("no content, embed nor file provided")
	// ErrAlreadyConnectedToVoice is returned when trying to join a voice channel in
	// a guild where you are already have an active voice connection.
	ErrAlreadyConnectedToVoice = errors.New("already connected to a voice channel in this guild, consider using the SwitchVoiceChannel method")
	// ErrNotConnectedToVoice is returned when trying to switch to a different voice
	// channel in a guild where you are not yet connected to a voice channel.
	ErrNotConnectedToVoice = errors.New("not connected to a voice channel in this guild, use the JoinVoiceChannel method first")
	// ErrNotCurrentUser is returned for user endpoints used with an ID different than "@me".
	ErrNotCurrentUser = errors.New("endpoint only available for current user (@me)")
)

// APIError is a generic error returned by the Discord HTTP API.
type APIError struct {
	HTTPCode int      `json:"http_code"`
	Code     int      `json:"code"`
	Message  string   `json:"message"`
	Misc     []string `json:"_misc"`
}

// Error implements the error interface.
func (e APIError) Error() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("%d %s:", e.HTTPCode, http.StatusText(e.HTTPCode)))

	if e.Message != "" {
		s.WriteString(fmt.Sprintf(" %s", e.Message))
	}

	if e.Code != 0 {
		s.WriteString(fmt.Sprintf(" (code: %d)", e.Code))
	}

	var i int
	for _, m := range e.Misc {
		if i > 0 {
			s.WriteRune(',')
		}

		s.WriteString(fmt.Sprintf(" %s", m))
		i++
	}
	return s.String()
}

// ValidationError is a validation error returned by the Discord HTTP API
// when it receives invalid parameters.
type ValidationError struct {
	HTTPCode int
	Errors   map[string]json.RawMessage `json:"errors"`
}

// Error implements the error interface.
func (e ValidationError) Error() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("%d %s:", e.HTTPCode, http.StatusText(e.HTTPCode)))

	var i int
	for key, errs := range e.Errors {
		if i > 0 {
			s.WriteRune(',')
		}

		s.WriteString(fmt.Sprintf(" field %q: %v", key, string(errs)))
		i++
	}
	return s.String()
}

// NewAPIError is a helper function that extracts an API error from
// an HTTP response and returns it as a generic APIError or a ValidationError.
func NewAPIError(resp *http.Response) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusBadRequest {
		validationErr := &ValidationError{HTTPCode: resp.StatusCode}
		if err = json.Unmarshal(b, &validationErr); err != nil {
			return err
		}
		return validationErr
	}

	apiErr := &APIError{HTTPCode: resp.StatusCode}
	if err = json.Unmarshal(b, &apiErr); err != nil {
		return err
	}

	return apiErr
}
