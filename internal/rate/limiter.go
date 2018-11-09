package rate

import (
	"net/http"
	"sync"
)

// Limiter holds a collection of buckets to track global and per-route
// rate limits. Create one with NewLimiter.
type Limiter struct {
	mu sync.Mutex

	global  *bucket
	buckets map[string]*bucket
}

// NewLimiter returns an initialized and ready to use Limiter.
func NewLimiter() *Limiter {
	return &Limiter{
		global:  &bucket{},
		buckets: make(map[string]*bucket),
	}
}

// Wait waits for a request to be theoretically safe to be sent (meaning it should
// not result in a 429 TO MANY REQUESTS) given the requested endpoint's key.
func (r *Limiter) Wait(key string) {
	r.mu.Lock()

	if r.global.enabled {
		r.mu.Unlock()
		r.global.lockAndWait()
	} else {
		_, ok := r.buckets[key]
		if !ok {
			r.buckets[key] = &bucket{}
		}
		b := r.buckets[key]
		r.mu.Unlock()
		b.lockAndWait()
	}
}

// Update updates the rate limit for an endpoint given its key.
func (r *Limiter) Update(key string, header http.Header) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.global.enabled {
		// We were globally rate limited and still are.
		if global := header.Get("X-RateLimit-Global"); global != "" {
			r.global.updateAndUnlock(header)
		} else { // We were globally rate limited but not anymore.
			r.global.unlock()
			r.buckets[key].update(header)
		}
	} else {
		// We were not globally rate limited but now we are.
		if global := header.Get("X-RateLimit-Global"); global != "" {
			r.global.update(header)
			r.buckets[key].unlock()
		} else { // We were not globally rate limited and still are not.
			r.buckets[key].updateAndUnlock(header)
		}
	}
}
