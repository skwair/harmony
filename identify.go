package harmony

import (
	"context"
	"runtime"
	"strings"
	"sync/atomic"
)

// identify is used to trigger the initial handshake with the gateway.
type identify struct {
	Token              string            `json:"token"`
	Properties         map[string]string `json:"properties"`
	Compress           bool              `json:"compress,omitempty"`
	LargeThreshold     int               `json:"large_threshold,omitempty"`
	Shard              *[2]int           `json:"shard,omitempty"`
	Presence           *Status           `json:"presence,omitempty"`
	GuildSubscriptions bool              `json:"guild_subscriptions"`
}

// Status is sent by the client to indicate a presence or status update.
type Status struct {
	Since  int       `json:"since"`
	Game   *Activity `json:"game"`
	Status string    `json:"status"`
	AFK    bool      `json:"afk"`
}

// identify sends an Identify payload to the Gateway.
func (c *Client) identify(ctx context.Context) error {
	i := &identify{
		Token: c.token,
		Properties: map[string]string{
			"$os":      strings.Title(runtime.GOOS),
			"$browser": "github.com/skwair/harmony",
		},
		Compress:           true,
		LargeThreshold:     c.largeThreshold,
		GuildSubscriptions: c.guildSubscriptions,
	}

	if c.shard[1] != 0 {
		i.Shard = &[2]int{c.shard[0], c.shard[1]}
	}

	return c.sendPayload(ctx, gatewayOpcodeIdentify, i)
}

// resume is used to replay missed events when a disconnected client resumes.
type resume struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       int64  `json:"seq"`
}

// resume sends a Resume payload to the Gateway.
func (c *Client) resume(ctx context.Context) error {
	r := &resume{
		Token:     c.token,
		SessionID: c.sessionID,
		Seq:       atomic.LoadInt64(&c.sequence),
	}
	return c.sendPayload(ctx, gatewayOpcodeResume, r)
}
