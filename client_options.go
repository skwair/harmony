package discord

import (
	"net/http"
	"time"
)

// ClientOption is a function that configures a Client.
// It is used in NewClient.
type ClientOption func(*Client)

// WithToken sets the token for a user client. Every call to
// NewClient must include this option or the WithBotToken
// option if the client is a bot instead of a regular user.
func WithToken(token string) ClientOption {
	return func(c *Client) {
		c.token = token
	}
}

// WithBotToken sets the token for a bot client. Every call to
// NewClient must include this option or the WithToken option
// if the client is a regular user instead of a bot.
func WithBotToken(token string) ClientOption {
	return func(c *Client) {
		c.token = "Bot " + token
		c.bot = true
	}
}

// WithHTTPClient can be used to specify the http.Client to use when making
// HTTP requests to the Discord HTTP API.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.client = client
	}
}

// WithBaseURL can be used to change de base URL of the API.
// This is used for testing.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithErrorHandler allows you to specify a custom error handler function
// that will be called whenever an error occurs while the connection
// to the Gateway is up.
func WithErrorHandler(h func(error)) ClientOption {
	return func(c *Client) {
		c.errorHandler = h
	}
}

// WithSharding allows you to specify a sharding configuration when connecting to the Gateway.
// See https://discordapp.com/developers/docs/topics/gateway#sharding for more details.
func WithSharding(current, total int) ClientOption {
	return func(c *Client) {
		c.shard[0] = current
		c.shard[1] = total
	}
}

// WithStateTracking allows you to specify whether the client is tracking the state of
// the current connection or not.
func WithStateTracking(y bool) ClientOption {
	return func(c *Client) {
		c.withStateTracking = y
	}
}

// WithLargeThreshold allows you to set the large threshold when connecting to the Gateway.
// This threshold will dictate the number of offline guild members are returned with a guild.
// See: https://discordapp.com/developers/docs/topics/gateway#request-guild-members for more details.
func WithLargeThreshold(t int) ClientOption {
	return func(c *Client) {
		if t > 250 {
			t = 250
		}
		if t < 0 {
			t = 0
		}
		c.largeThreshold = t
	}
}

// WithBackoffStrategy allows you to customize the backoff strategy used when trying
// to reconnect to the Discord Gateway after an error occurred (such as a network
// failure).
func WithBackoffStrategy(baseDelay, maxDelay time.Duration, factor, jitter float64) ClientOption {
	return func(c *Client) {
		c.backoff.baseDelay = baseDelay
		c.backoff.maxDelay = maxDelay
		c.backoff.factor = factor
		c.backoff.jitter = jitter
	}
}
