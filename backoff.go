package harmony

import (
	"math/rand"
	"time"
)

// backoff configures a backoff strategy.
type backoff struct {
	// baseDelay is the amount of time to wait before retrying after the first
	// failure.
	baseDelay time.Duration
	// maxDelay is the upper bound of backoff delay.
	maxDelay time.Duration
	// factor is applied to the backoff after each attempt.
	factor float64
	// jitter provides a range to randomize backoff delays.
	jitter float64
}

// forAttempt returns the duration to wait for the n-th attempt.
func (b backoff) forAttempt(n int) time.Duration {
	if n == 0 {
		return b.baseDelay
	}

	backoff, max := float64(b.baseDelay), float64(b.maxDelay)
	for backoff < max && n > 0 {
		backoff *= b.factor
		n--
	}

	if backoff > max {
		backoff = max
	}

	backoff *= 1 + b.jitter*(rand.Float64()*2-1)

	if backoff < 0 {
		return 0
	}

	return time.Duration(backoff)
}
