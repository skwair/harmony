package harmony

import (
	"github.com/skwair/harmony/resource/channel"
	"github.com/skwair/harmony/resource/guild"
	"github.com/skwair/harmony/resource/invite"
	"github.com/skwair/harmony/resource/user"
	"github.com/skwair/harmony/resource/webhook"
)

// Guild returns a new guild resource to manage the guild with the given ID.
func (c *Client) Guild(id string) *guild.Resource {
	return guild.NewResource(c.restClient, id)
}

// Channel returns a new channel resource to manage the channel with the given ID.
func (c *Client) Channel(id string) *channel.Resource {
	return channel.NewResource(c.restClient, id)
}

// User returns a new user resource to manage the user with the given ID.
// Note that most methods on this resource are only available for the current
// user (@me).
func (c *Client) User(id string) *user.Resource {
	return user.NewResource(c.restClient, id)
}

// Webhook returns a new webhook resource to manage the webhook with the given ID.
func (c *Client) Webhook(id string) *webhook.Resource {
	return webhook.NewResource(c.restClient, id)
}

// Invite returns a new invite resource to manage the invite with the given code.
func (c *Client) Invite(code string) *invite.Resource {
	return invite.NewResource(c.restClient, code)
}
