package harmony

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/skwair/harmony/internal/endpoint"
)

// requestPayload is a payload that is sent to Discord's REST API.
type requestPayload struct {
	body        []byte
	contentType string
}

// jsonPayload creates a new requestPayload from some raw JSON data.
func jsonPayload(body json.RawMessage) *requestPayload {
	return &requestPayload{
		body:        body,
		contentType: "application/json",
	}
}

// customPayload creates a new custom payload from raw bytes and a given content type.
func customPayload(body []byte, contentType string) *requestPayload {
	return &requestPayload{
		body:        body,
		contentType: contentType,
	}
}

// doReq is used to request Discord's HTTP endpoints.
// If you need more control over headers you send, use doReqWithHeader directly.
func (c *Client) doReq(ctx context.Context, e *endpoint.Endpoint, p *requestPayload) (*http.Response, error) {
	return c.doReqWithHeader(ctx, e, p, nil)
}

// doReqWithHeader sends an HTTP request and returns the response given an endpoint
// an optional payload and some headers. It adds the required Authorization header,
// Content-Type based on the given payload and also sets the User-Agent.
// It also takes care of rate limiting, using the client's built in rate limiter.
func (c *Client) doReqWithHeader(ctx context.Context, e *endpoint.Endpoint, p *requestPayload, h http.Header) (*http.Response, error) {
	var (
		err error
		req *http.Request
	)
	if p != nil && p.body != nil {
		req, err = http.NewRequest(e.Method, c.baseURL+e.URL, bytes.NewReader(p.body))
	} else {
		req, err = http.NewRequest(e.Method, c.baseURL+e.URL, nil)
	}
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	// Add custom headers provided. This has to be done
	// before adding other mandatory headers to make
	// sure they are not overridden.
	for k, vs := range h {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	// Add the Content-Type header accordingly to the payload's body, if any.
	if p != nil && p.body != nil {
		req.Header.Set("Content-Type", p.contentType)
	}
	// Add the Authorization header.
	req.Header.Set("Authorization", c.token)
	// Finally, set the User-Agent header.
	ua := fmt.Sprintf("%s (github.com/skwair/harmony, %s)", c.name, version)
	req.Header.Set("User-Agent", ua)

	c.limiter.Wait(e.Key)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	c.limiter.Update(e.Key, resp.Header)

	// We are being rate limited, rate limiter has been updated
	// and will wait before sending future requests, but we must
	// try and resend this one since it was rejected.
	// NOTE: this should never happen since we now wait
	// before sending requests.
	if resp.StatusCode == http.StatusTooManyRequests {
		return c.doReqWithHeader(ctx, e, p, h)
	}

	return resp, nil
}

// rateLimit is the JSON body Discord sends when we are rate limited.
type rateLimit struct {
	Message    string `json:"message"`
	RetryAfter int    `json:"retry_after"`
	Global     bool   `json:"global"`
}

// doReqNoAuth is used to request endpoints that do not need authentication.
// If you need more control over headers you send, use doReqNoAuthWithHeader directly.
func doReqNoAuth(ctx context.Context, e *endpoint.Endpoint, p *requestPayload) (*http.Response, error) {
	return doReqNoAuthWithHeader(ctx, e, p, nil)
}

// doReqNoAuth is used to request endpoints that do not need authentication. It is
// like doReqWithHeader otherwise, except for rate limiting where it is more likely
// to result in 429's if abused.
func doReqNoAuthWithHeader(ctx context.Context, e *endpoint.Endpoint, p *requestPayload, h http.Header) (*http.Response, error) {
	var (
		err error
		req *http.Request
	)
	if p != nil && p.body != nil {
		req, err = http.NewRequest(e.Method, defaultBaseURL+e.URL, bytes.NewReader(p.body))
	} else {
		req, err = http.NewRequest(e.Method, defaultBaseURL+e.URL, nil)
	}
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	for k, vs := range h {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	if p != nil && p.body != nil {
		h.Set("Content-Type", p.contentType)
	}
	ua := fmt.Sprintf("%s (github.com/skwair/harmony, %s", "Harmony", version)
	req.Header.Set("User-Agent", ua)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// We are being rate limited, wait a bit and resend the request.
	// NOTE: maybe use HTTP headers (if set) instead of having to
	// parse some JSON.
	if resp.StatusCode == http.StatusTooManyRequests {
		var r rateLimit
		if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
			return nil, err
		}
		time.Sleep(time.Millisecond * time.Duration(r.RetryAfter))
		return doReqNoAuthWithHeader(ctx, e, p, h)
	}

	return resp, nil
}
