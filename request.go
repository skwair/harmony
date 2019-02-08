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

// doReq calls doReqWithHeader with the Content-Type header set to "application/json"
// if the provided body is not nil.
// If you need more control over headers you send, use doReqWithHeader directly.
func (c *Client) doReq(ctx context.Context, e *endpoint.Endpoint, body []byte) (*http.Response, error) {
	h := http.Header{}
	if body != nil {
		h.Set("Content-Type", "application/json")
	}

	return c.doReqWithHeader(ctx, e, body, h)
}

// doReqWithHeader sends an HTTP request and returns the response given an endpoint
// an optional body and some headers. It adds the required Authorization header and
// also sets the User-Agent.
// It takes care of rate limiting, using the client's built in rate limiter.
func (c *Client) doReqWithHeader(ctx context.Context, e *endpoint.Endpoint, body []byte, h http.Header) (*http.Response, error) {
	req, err := http.NewRequest(e.Method, c.baseURL+e.URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	// Merge h into req.Header, then set the Authorization
	// and User-Agent header.
	for k, vs := range h {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	req.Header.Set("Authorization", c.token)
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
	if resp.StatusCode == http.StatusTooManyRequests {
		return c.doReqWithHeader(ctx, e, body, h)
	}

	return resp, nil
}

// rateLimit is the JSON body Discord sends when we are rate limited.
type rateLimit struct {
	Message    string `json:"message"`
	RetryAfter int    `json:"retry_after"`
	Global     bool   `json:"global"`
}

// doReqNoAuth is used to request endpoints that do not need authentication. It sets
// the Content-Type to "application/json" if the body is not nil.
// If you need more control over headers you send, use doReqNoAuthWithHeader directly.
func doReqNoAuth(ctx context.Context, e *endpoint.Endpoint, body []byte) (*http.Response, error) {
	h := http.Header{}
	if body != nil {
		h.Set("Content-Type", "application/json")
	}

	return doReqNoAuthWithHeader(ctx, e, body, h)
}

// doReqNoAuth is used to request endpoints that do not need authentication. It is
// like doReqWithHeader otherwise, except for rate limiting where it is more likely
// to result in 429's if abused.
func doReqNoAuthWithHeader(ctx context.Context, e *endpoint.Endpoint, body []byte, h http.Header) (*http.Response, error) {
	req, err := http.NewRequest(e.Method, defaultBaseURL+e.URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	for k, vs := range h {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
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
		return doReqNoAuth(ctx, e, body)
	}

	return resp, nil
}
