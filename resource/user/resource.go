package user

import (
	"github.com/skwair/harmony/resource"
)

// Resource is a resource that allows to perform various actions on a Discord user.
// Create one with Client.User.
type Resource struct {
	userID string
	client resource.RestClient
}

func NewResource(c resource.RestClient, id string) *Resource {
	return &Resource{client: c, userID: id}
}
