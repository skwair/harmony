package resource

import (
	"context"
	"net/http"

	"github.com/skwair/harmony/internal/endpoint"
	"github.com/skwair/harmony/internal/rest"
)

type RestClient interface {
	Do(ctx context.Context, e *endpoint.Endpoint, p *rest.Payload) (*http.Response, error)
	DoWithHeader(ctx context.Context, e *endpoint.Endpoint, p *rest.Payload, h http.Header) (*http.Response, error)
}
