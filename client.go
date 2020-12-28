package harmony

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/skwair/harmony/discord"
	"github.com/skwair/harmony/internal/backoff"
	"github.com/skwair/harmony/internal/payload"
	"github.com/skwair/harmony/internal/rest"
	"github.com/skwair/harmony/log"
	"github.com/skwair/harmony/voice"
	"go.uber.org/atomic"
	"nhooyr.io/websocket"
)

const defaultLargeThreshold = 250

var (
	// defaultBackoff is the backoff strategy used by default when trying to reconnect to the Gateway.
	defaultBackoff = backoff.NewExponential(1*time.Second, 120*time.Second, 1.6, 0.2)
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
	// This is the initial status of the bot that is sent
	// when identifying to the Gateway. It can be customized
	// with WithInitialBotStatus and can be later modified
	// with SetBotStatus.
	initialBotStatus *discord.BotStatus

	// Authentication token used to interact with
	// Discord's API.
	token string

	gatewayURL string

	// Underlying HTTP client used to call Discord's REST API.
	httpClient *http.Client
	restClient *rest.Client

	// Underlying websocket used to communicate with
	// Discord's real-time API.
	conn    *websocket.Conn
	connRMu sync.Mutex // Read mutex.

	// Whether the client is currently connecting to the Gateway.
	connecting *atomic.Bool
	// Whether the client is currently connected to the Gateway.
	connected *atomic.Bool
	// Whether the client is currently connecting to a voice server.
	connectingToVoice *atomic.Bool
	// Whether the client is currently reconnecting to the Gateway.
	reconnecting *atomic.Bool

	// When connectingToVoice is true, some
	// payloads received by the event handler will
	// be sent through this channel.
	voicePayloads chan *payload.Payload

	// See WithLargeThreshold for more information.
	largeThreshold int
	// See WithSharding for more information.
	shard [2]int
	// See WithGuildSubscriptions for more information.
	guildSubscriptions bool
	// See WithGatewayIntents for more information.
	intents discord.GatewayIntent

	userID    string
	sessionID string

	// Sequence number of the last Dispatch event
	// we received from the Gateway.
	sequence *atomic.Int64
	// UNIX timestamp in nanoseconds of the last
	// heartbeat acknowledgement.
	lastHeartbeatAck *atomic.Int64
	// UNIX timestamp in nanoseconds of the last
	// heartbeat sent.
	lastHeartbeatSent *atomic.Int64

	// wg keeps track of all goroutines necessary to
	// maintain a connection to the Gateway.
	wg sync.WaitGroup
	// The first fatal error encountered while connected
	// to the Gateway will be sent through this channel.
	error chan error
	// Used to ensure we only report the first error.
	reportErrorOnce sync.Once

	// Closing this channel will gracefully shutdown the
	// Gateway connection.
	stop chan struct{}

	// Shared context used for sending and receiving websocket
	// payloads. Will be canceled when the client disconnects
	// or an error occurs.
	ctx    context.Context
	cancel context.CancelFunc

	// Registered event handlers for this Client.
	handlersMu sync.RWMutex
	handlers   map[string]handler

	// Exponential strategy used when trying to reconnect to
	// the Gateway after an error.
	backoff *backoff.Exponential

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
		httpClient:         http.DefaultClient,
		largeThreshold:     defaultLargeThreshold,
		guildSubscriptions: true,
		intents:            discord.GatewayIntentUnprivileged,
		handlers:           make(map[string]handler),
		backoff:            defaultBackoff,
		withStateTracking:  true,
		voiceConnections:   make(map[string]*voice.Connection),
		logger:             log.NewStd(os.Stderr, log.LevelInfo),
		sequence:           atomic.NewInt64(0),
		lastHeartbeatSent:  atomic.NewInt64(0),
		lastHeartbeatAck:   atomic.NewInt64(0),
		connected:          atomic.NewBool(false),
		connecting:         atomic.NewBool(false),
		connectingToVoice:  atomic.NewBool(false),
		reconnecting:       atomic.NewBool(false),
	}

	for _, opt := range opts {
		opt(c)
	}

	c.restClient = rest.NewClient(
		c.httpClient,
		c.token,
		c.name,
		c.logger,
	)

	if c.withStateTracking {
		c.State = newState()
	}

	return c, nil
}
