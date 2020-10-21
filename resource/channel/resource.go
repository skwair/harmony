package channel

import (
	"github.com/skwair/harmony/resource"
)

// Resource is a resource that allows to perform various actions on a Discord channel.
// Create one with Client.Channel.
type Resource struct {
	channelID string
	client    resource.RestClient
}

func NewResource(c resource.RestClient, id string) *Resource {
	return &Resource{client: c, channelID: id}
}
