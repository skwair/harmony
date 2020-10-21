package webhook

import (
	"github.com/skwair/harmony/resource"
)

// Resource is a resource that allows to perform various actions on a Discord webhook.
// Create one with Client.Webhook.
type Resource struct {
	webhookID string
	client    resource.RestClient
}

func NewResource(c resource.RestClient, id string) *Resource {
	return &Resource{client: c, webhookID: id}
}
