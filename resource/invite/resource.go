package invite

import (
	"github.com/skwair/harmony/resource"
)

// Resource is a resource that allows to perform various actions on a Discord guild invite.
// Create one with Client.Invite.
type Resource struct {
	code   string
	client resource.RestClient
}

func NewResource(c resource.RestClient, code string) *Resource {
	return &Resource{client: c, code: code}
}
