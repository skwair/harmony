package harmony

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sync"
	"time"

	"go.uber.org/atomic"
	"nhooyr.io/websocket"

	"github.com/skwair/harmony/internal/payload"
	"github.com/skwair/harmony/internal/rate"
	"github.com/skwair/harmony/log"
	"github.com/skwair/harmony/voice"
)

const (
	defaultBaseURL        = "https://discordapp.com/api/v6"
	defaultLargeThreshold = 250
)

var (
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

	// This is the name of the bot, used to set the
	// User-Agent when sending HTTP request.
	// It defaults to "Harmony".
	name string

	// Authentication token used to interact with
	// Discord's API.
	token string

	gatewayURL string
	baseURL    string // Base URL of the Discord API.

	// Underlying HTTP client used to call Discord's REST API.
	client *http.Client

	// Rate limiter used to throttle outgoing HTTP requests.
	limiter *rate.Limiter

	// Underlying websocket used to communicate with
	// Discord's real-time API.
	conn    *websocket.Conn
	connRMu sync.Mutex // Read mutex.

	// Accessed atomically, acts as a boolean and is set to 1
	// when the client is connected to the Gateway.
	connected *atomic.Bool

	// Accessed atomically, acts as a boolean and is set to 1
	// when the client is connecting to voice.
	connectingToVoice *atomic.Bool
	// When connectingToVoice is set to 1, some
	// payloads received by the event handler will
	// be sent through this channel.
	voicePayloads chan *payload.Payload

	// See WithLargeThreshold for more information.
	largeThreshold int
	// See WithSharding for more information.
	shard [2]int
	// See WithGuildSubscriptions for more information.
	guildSubscriptions bool

	userID    string
	sessionID string
	// Accessed atomically, sequence number of the last
	// Dispatch event we received from the Gateway.
	sequence *atomic.Int64
	// Accessed atomically, UNIX timestamp in nanoseconds
	// of the last heartbeat acknowledgement.
	lastHeartbeatACK *atomic.Int64
	// Accessed atomically, UNIX timestamp in nanoseconds
	// of the last heartbeat send. Used to calculate RTT.
	lastHeartbeatSend *atomic.Int64

	// Those fields are used for synchronisation between
	// the listen, receive, heartbeat and wait goroutines
	// when the connection to the gateway is up.
	wg    sync.WaitGroup
	error chan error
	stop  chan struct{}

	// Shared context used for sending and receiving websocket
	// payloads. Will be canceled when the client disconnects
	// or an error occurs.
	ctx    context.Context
	cancel context.CancelFunc

	handlersMu sync.RWMutex
	handlers   map[string]handler

	// Backoff strategy used when trying to reconnect to
	// the Gateway after an error.
	backoff backoff

	// If true (the default value), the State
	// will be populated and updated as events
	// are received from the Discord Gateway.
	withStateTracking bool
	State             *State

	// voice connections that were established by
	// this client.
	voiceConnections map[string]*voice.Connection

	logger log.Logger
}

// NewClient creates a new client to work with Discord's API.
// It is meant to be long lived and shared across your application.
// The token is automatically prefixed with "Bot ", which is a requirement
// by Discord for bot users. Automated normal user accounts (generally called
// "self-bots"), are not supported. To customize a Client, refer to available
// ClientOption.
func NewClient(token string, opts ...ClientOption) (*Client, error) {
	if token == "" {
		return nil, errors.New("harmony: a token is mandatory to create a client")
	}

	c := &Client{
		name:               "Harmony",
		token:              "Bot " + token,
		baseURL:            defaultBaseURL,
		client:             http.DefaultClient,
		limiter:            rate.NewLimiter(),
		largeThreshold:     defaultLargeThreshold,
		guildSubscriptions: true,
		handlers:           make(map[string]handler),
		backoff:            defaultBackoff,
		withStateTracking:  true,
		voiceConnections:   make(map[string]*voice.Connection),
		logger:             log.NewStd(os.Stderr, log.LevelError),
		sequence:           atomic.NewInt64(0),
		lastHeartbeatSend:  atomic.NewInt64(0),
		lastHeartbeatACK:   atomic.NewInt64(0),
		connected:          atomic.NewBool(false),
		connectingToVoice:  atomic.NewBool(false),
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.withStateTracking {
		c.State = newState()
	}

	return c, nil
}
