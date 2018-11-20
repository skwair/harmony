package harmony

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/skwair/harmony/internal/rate"
)

const (
	defaultBaseURL        = "https://discordapp.com/api/v6"
	defaultLargeThreshold = 250
)

var (
	// ErrNoTokenProvided is returned by NewClient when neither a user token nor a bot token is provided.
	ErrNoTokenProvided = errors.New("no token provided, use WithToken or WithBotToken when creating a new client")

	// defaultErrorHandler is the default handle that is called when an error occurs when connected to the Gateway.
	defaultErrorHandler = func(err error) { log.Println("gateway client error:", err) }

	// defaultBackoff is the backoff strategy used by default when trying to reconnect to the Gateway.
	defaultBackoff = backoff{
		baseDelay: 1 * time.Second,
		maxDelay:  120 * time.Second,
		factor:    1.6,
		jitter:    0.2,
	}
)

// Client is used to communicate with Discord's API.
// To start receiving events from the Gateway with a
// Client, you first need to call its Connect method.
type Client struct {
	// General lock for long operations that should
	// not happen concurrently like Connect or Disconnect.
	mu sync.Mutex

	token string
	// Whether this is a regular user or a bot user.
	// This is set depending on which of WithToken
	// or WithBotToken client option is used when
	// creating the client.
	bot bool

	gatewayURL string
	baseURL    string // Base URL of the Discord API.

	// Underlying HTTP client used to call Discord's API.
	client *http.Client

	// Underlying websocket used to communicate with
	// Discord's real-time API.
	conn    *websocket.Conn
	connRMu sync.Mutex // Read mutex.
	connWMu sync.Mutex // Write mutex.

	// Accessed atomically, acts as a boolean and is set to 1
	// when the client is connected to the Gateway.
	connected int32

	// Accessed atomically, acts as a boolean and is set to 1
	// when the client is connecting to voice.
	connectingToVoice int32
	// When connectingToVoice is set to 1, some
	// payloads received by the event handler will
	// be sent through this channel.
	voicePayloads chan *payload

	// See WithLargeThreshold for more information.
	largeThreshold int
	// See WithSharding for more information.
	shard [2]int

	userID    string
	sessionID string
	// Accessed atomically, sequence number of the last
	// Dispatch event we received from the Gateway.
	sequence int64
	// Accessed atomically, UNIX timestamp in nanoseconds
	// of the last heartbeat acknowledgement.
	lastHeartbeatACK int64
	// Accessed atomically, UNIX timestamp in nanoseconds
	// of the last heartbeat send. Used to calculate RTT.
	lastHeartbeatSend int64

	// Those fields are used for synchronisation between
	// the listen, receive and heartbeat goroutines when
	// the connection to the gateway is up.
	wg    sync.WaitGroup
	error chan error
	stop  chan struct{}

	handlersMu   sync.RWMutex
	handlers     map[string]handler
	errorHandler func(error)

	limiter *rate.Limiter

	// Used when trying to reconnect to
	// the Gateway after an error.
	backoff backoff

	// If true (the default value), the State
	// will be populated and updated as events
	// are received from the Discord Gateway.
	withStateTracking bool
	State             *State
}

// NewClient creates a new client to work with Discord's API.
// It is meant to be long lived and shared across your application.
func NewClient(opts ...ClientOption) (*Client, error) {
	c := &Client{
		baseURL:           defaultBaseURL,
		client:            http.DefaultClient,
		largeThreshold:    defaultLargeThreshold,
		errorHandler:      defaultErrorHandler,
		handlers:          make(map[string]handler),
		limiter:           rate.NewLimiter(),
		backoff:           defaultBackoff,
		withStateTracking: true,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.token == "" {
		return nil, ErrNoTokenProvided
	}

	if c.withStateTracking {
		c.State = newState()
	}

	return c, nil
}
