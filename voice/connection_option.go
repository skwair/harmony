package voice

import "github.com/skwair/harmony/log"

// ConnectionOption is a function that configures a Connection.
// It is used in Connect.
type ConnectionOption func(*Connection)

// WithLogger can be used to set the logger used by this connection.
// Defaults to a standard logger reporting only errors.
// See the log package for more information about logging with Harmony.
func WithLogger(l log.Logger) ConnectionOption {
	return func(c *Connection) {
		c.logger = l
	}
}
