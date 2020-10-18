package harmony

import (
	"net/http"
	"time"

	"github.com/skwair/harmony/log"
)

// ClientOption is a function that configures a Client.
// It is used in NewClient.
type ClientOption func(*Client)

// WithName sets the name of the client. It will be used to
// set the User-Agent of HTTP requests sent by the Client.
// Defaults to "Harmony".
func WithName(n string) ClientOption {
	return func(c *Client) {
		c.name = n
	}
}

// WithHTTPClient can be used to specify the http.Client to use when making
// HTTP requests to the Discord HTTP API.
// Defaults to http.DefaultClient.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.client = client
	}
}

// WithBaseURL can be used to change de base URL of the API.
// This is used for testing.
// Deprecated.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithSharding allows you to specify a sharding configuration when connecting to the Gateway.
// See https://discord.com/developers/docs/topics/gateway#sharding for more details.
// Defaults to nothing, sharding is not enabled.
func WithSharding(current, total int) ClientOption {
	return func(c *Client) {
		c.shard[0] = current
		c.shard[1] = total
	}
}

// WithGuildSubscriptions allows to set whether the client should identify to the Gateway with
// guild subscription enabled or not. Guild subscriptions are guild member presence updates
// and typing events.
// Defaults to true.
// While not deprecated, Guild Subscriptions have been superseded by Gateway Intents. It is recommended
// to use WithGatewayIntents for better results.
func WithGuildSubscriptions(y bool) ClientOption {
	return func(c *Client) {
		c.guildSubscriptions = y
	}
}

// WithGatewayIntents allows to customize which Gateway Intents the client should subscribe to.
// See https://discord.com/developers/docs/topics/gateway#gateway-intents for more information.
// By default, the client subscribes to all events.
func WithGatewayIntents(i GatewayIntent) ClientOption {
	return func(c *Client) {
		c.intents = i
	}
}

// WithStateTracking allows you to specify whether the client is tracking the state of
// the current connection or not.
// Defaults to true.
func WithStateTracking(y bool) ClientOption {
	return func(c *Client) {
		c.withStateTracking = y
	}
}

// WithLargeThreshold allows you to set the large threshold when connecting to the Gateway.
// This threshold will dictate the number of offline guild members are returned with a guild.
// See: https://discord.com/developers/docs/topics/gateway#request-guild-members for more details.
// Defaults to 250.
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
// Defaults to 1s (baseDelay), 120s (maxDelay), 1.6 (factor), 0.2 (jitter).
func WithBackoffStrategy(baseDelay, maxDelay time.Duration, factor, jitter float64) ClientOption {
	return func(c *Client) {
		c.backoff.baseDelay = baseDelay
		c.backoff.maxDelay = maxDelay
		c.backoff.factor = factor
		c.backoff.jitter = jitter
	}
}

// WithLogger can be used to set the logger used by Harmony.
// Defaults to a standard logger reporting only errors.
// See the log package for more information about logging with Harmony.
func WithLogger(l log.Logger) ClientOption {
	return func(c *Client) {
		c.logger = l
	}
}
