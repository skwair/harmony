package guild

import (
	"github.com/skwair/harmony/resource"
)

// Resource is a resource that allows to perform various actions on a Discord guild.
// Create one with Client.Guild.
type Resource struct {
	guildID string
	client  resource.RestClient
}

func NewResource(c resource.RestClient, id string) *Resource {
	return &Resource{client: c, guildID: id}
}
