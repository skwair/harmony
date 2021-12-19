package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/skwair/harmony/internal/endpoint"
	"github.com/skwair/harmony/internal/rate"
	"github.com/skwair/harmony/log"
	"github.com/skwair/harmony/version"
)

var baseURL = "https://discord.com/api/v" + version.REST()

// rateLimitResp is the JSON body Discord sends when we are rate limited.
type rateLimitResp struct {
	Message    string  `json:"message"`
	RetryAfter float64 `json:"retry_after"`
	Global     bool    `json:"global"`
}

// Client is a client that can make HTTP requests to Discord's REST API.
type Client struct {
	httpClient *http.Client
	limiter    *rate.Limiter
	token      string
	name       string
	logger     log.Logger
}

// NewClient returns a new REST Client.
func NewClient(httpClient *http.Client, token, name string, logger log.Logger) *Client {
	return &Client{
		httpClient: httpClient,
		limiter:    rate.NewLimiter(),
		token:      token,
		name:       name,
		logger:     logger,
	}
}

// Do is used to request Discord's HTTP endpoints.
// If you need more control over headers you send, use DoWithHeader directly.
func (c *Client) Do(ctx context.Context, e *endpoint.Endpoint, p *Payload) (*http.Response, error) {
	return c.DoWithHeader(ctx, e, p, nil)
}

// DoWithHeader sends an HTTP request and returns the response given an endpoint
// an optional payload and some headers. It adds the required Authorization header,
// Content-Type based on the given payload and also sets the User-Agent.
// It also takes care of rate limiting, using the client's built in rate limiter.
func (c *Client) DoWithHeader(ctx context.Context, e *endpoint.Endpoint, p *Payload, h http.Header) (*http.Response, error) {
	var (
		err error
		req *http.Request
	)
	if p.hasBody() {
		req, err = http.NewRequestWithContext(ctx, e.Method, baseURL+e.Path, bytes.NewReader(p.body))
	} else {
		req, err = http.NewRequestWithContext(ctx, e.Method, baseURL+e.Path, nil)
	}
	if err != nil {
		return nil, err
	}

	// Add custom headers provided. This has to be done
	// before adding other mandatory headers to make
	// sure they are not overridden.
	for k, vs := range h {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	// Add the Content-Type header accordingly to the payload's body, if any.
	if p.hasBody() {
		req.Header.Set("Content-Type", p.contentType)
	}
	// Add the Authorization header.
	req.Header.Set("Authorization", c.token)
	// Finally, set the User-Agent header.
	ua := fmt.Sprintf("%s (github.com/skwair/harmony, %s)", c.name, version.Module())
	req.Header.Set("User-Agent", ua)

	c.limiter.Wait(e.Key)

	if c.logger.Level() == log.LevelDebug {
		b, _ := httputil.DumpRequestOut(req, true)
		c.logger.Debug("--> ", string(b))
	}

	before := time.Now()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if c.logger.Level() == log.LevelDebug {
		b, _ := httputil.DumpResponse(resp, true)
		c.logger.Debug("<-- ", time.Since(before), "\n", string(b))
	}

	c.limiter.Update(e.Key, resp.Header)

	// Make sure we agree on time with the server, otherwise rate limit would be inaccurate.
	date, err := http.ParseTime(resp.Header.Get("Date"))
	if err != nil {
		return nil, fmt.Errorf("could not parse date header: %w", err)
	}

	now := time.Now()

	// Only print the warning if the request took less than one second, otherwise it
	// could just be a very high network latency but not a time desynchronization.
	// NOTE: these values probably need some tweaking.
	if now.Sub(before) < time.Second &&
		(now.Before(date.Add(-1500*time.Millisecond)) ||
			now.After(date.Add(1500*time.Millisecond))) {
		c.logger.Warnf("time desynchronization detected (server UTC time: %s, local UTC time: %s), rate limit will be inaccurate and you may encounter 429s, consider using NTP to synchronize time", date.UTC(), now.Round(time.Second).UTC())
	}

	// We are being rate limited, rate limiter has been updated
	// and will wait before sending future requests, but we must
	// try and resend this one since it was rejected.
	// NOTE: this should never happen as long as our time is in
	// sync with Discord servers since we wait before sending requests.
	// Still, keep this check to prevent spamming in the event where
	// a time desynchronization happens.
	if resp.StatusCode == http.StatusTooManyRequests {
		var r rateLimitResp
		if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
			return nil, err
		}

		time.Sleep(time.Millisecond * time.Duration(r.RetryAfter*100))

		return c.DoWithHeader(ctx, e, p, h)
	}

	return resp, nil
}

// Do is used to request endpoints that do not need authentication.
// If you need more control over headers you send, use DoWithHeader directly.
func Do(ctx context.Context, e *endpoint.Endpoint, p *Payload) (*http.Response, error) {
	return DoWithHeader(ctx, e, p, http.Header{})
}

// DoWithHeader is used to request endpoints that do not need authentication. It is
// like Client.DoWithHeader otherwise, except for rate limiting where it is more likely
// to result in 429's if abused.
func DoWithHeader(ctx context.Context, e *endpoint.Endpoint, p *Payload, h http.Header) (*http.Response, error) {
	var (
		err error
		req *http.Request
	)
	if p.hasBody() {
		req, err = http.NewRequest(e.Method, baseURL+e.Path, bytes.NewReader(p.body))
	} else {
		req, err = http.NewRequest(e.Method, baseURL+e.Path, nil)
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
	if p.hasBody() {
		req.Header.Set("Content-Type", p.contentType)
	}
	ua := fmt.Sprintf("%s (github.com/skwair/harmony, %s)", "Harmony", version.Module())
	req.Header.Set("User-Agent", ua)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// We are being rate limited, wait a bit and resend the request.
	// NOTE: maybe use HTTP headers (if set) instead of having to
	// parse some JSON.
	if resp.StatusCode == http.StatusTooManyRequests {
		var r rateLimitResp
		if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
			return nil, err
		}

		time.Sleep(time.Millisecond * time.Duration(r.RetryAfter*100))

		return DoWithHeader(ctx, e, p, h)
	}

	return resp, nil
}

// ReasonHeader returns an HTTP header with the Audit Log reason set to r.
func ReasonHeader(r string) http.Header {
	h := http.Header{}

	if r != "" {
		h.Set("X-Audit-Log-Reason", r)
	}

	return h
}
