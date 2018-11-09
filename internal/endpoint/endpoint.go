package endpoint

// Endpoint is a Discord's HTTP endpoint along with its key, used for rate limiting.
type Endpoint struct {
	URL string
	Key string
}
