package rate

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

// bucket implements a leaky bucket.
type bucket struct {
	mu sync.Mutex

	// Whether this bucket is enabled or not. If disabled (the default),
	// lockAndWait will always return immediately. This is used for endpoints
	// that do not have a specific rate limit (besides the global rate limit).
	enabled bool

	// Maximum number of tokens this bucket can hold.
	limit int
	// Remaining tokens in the bucket.
	remaining int
	// Unix timestamp for when this bucket
	// refills to its maximum capacity.
	reset int64
}

// lockAndWait locks the bucket, returning immediately after if the bucket is disabled.
// If it is enabled, il will decrement the remaining tokens in the bucket by one if there
// is at least one token remaining, else it will wait for the bucket to refill before
// doing so.
func (b *bucket) lockAndWait() {
	b.mu.Lock()

	if !b.enabled {
		return
	}

	// Reset time is in the past, refill the bucket to its maximum capacity.
	if time.Unix(b.reset, 0).Before(time.Now()) {
		b.remaining = b.limit
	}

	if b.remaining == 0 {
		// We are out of tokens in this bucket, wait until it refills.
		time.Sleep(time.Unix(b.reset, 0).Sub(time.Now()))
		b.remaining = b.limit
	}

	b.remaining--
}

// update updates the bucket by parsing the given HTTP header.
// If none is present, the bucket is disabled.
// NOTE: errors are discarded for now, consider returning them.
func (b *bucket) update(header http.Header) {
	var set bool

	if limit := header.Get("X-RateLimit-Limit"); limit != "" {
		l, _ := strconv.ParseInt(limit, 10, 64)
		b.limit = int(l)
		set = true
	}

	if remaining := header.Get("X-RateLimit-Remaining"); remaining != "" {
		r, _ := strconv.ParseInt(remaining, 10, 64)
		b.remaining = int(r)
		set = true
	}

	if reset := header.Get("X-RateLimit-Reset"); reset != "" {
		r, _ := strconv.ParseInt(reset, 10, 64)
		b.reset = r
		set = true
	}

	b.enabled = set
}

// unlock unlocks the bucket without updating it.
func (b *bucket) unlock() {
	b.mu.Unlock()
}

// updateAndUnlock calls update with the given header, then unlocks the bucket.
func (b *bucket) updateAndUnlock(header http.Header) {
	b.update(header)

	b.unlock()
}
