package harmony

import (
	"context"
	"runtime"
	"strings"

	"github.com/skwair/harmony/discord"
)

// identify is used to trigger the initial handshake with the gateway.
type identify struct {
	Token              string                `json:"token"`
	Properties         map[string]string     `json:"properties"`
	Compress           bool                  `json:"compress,omitempty"`
	LargeThreshold     int                   `json:"large_threshold,omitempty"`
	Shard              *[2]int               `json:"shard,omitempty"`
	Presence           *discord.BotStatus    `json:"presence,omitempty"`
	GuildSubscriptions bool                  `json:"guild_subscriptions"`
	Intents            discord.GatewayIntent `json:"intents"`
}

// identify sends an Identify payload to the Gateway.
func (c *Client) identify(ctx context.Context) error {
	i := &identify{
		Token: c.token,
		Properties: map[string]string{
			"$os":      strings.Title(runtime.GOOS),
			"$browser": "github.com/skwair/harmony",
		},
		Presence:           c.initialBotStatus,
		Compress:           true,
		LargeThreshold:     c.largeThreshold,
		GuildSubscriptions: c.guildSubscriptions,
		Intents:            c.intents,
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
		Seq:       c.sequence.Load(),
	}
	return c.sendPayload(ctx, gatewayOpcodeResume, r)
}
