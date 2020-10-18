package endpoint

// Endpoint represent a single REST endpoint exposed by Discord's API. It
// consists of an HTTP method, a path as well as a key, used for rate limiting.
type Endpoint struct {
	Method string
	Path   string
	Key    string
}
