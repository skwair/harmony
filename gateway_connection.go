package harmony

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"nhooyr.io/websocket"

	"github.com/skwair/harmony/internal/payload"
)

const (
	gatewayVersion  = 6
	gatewayEncoding = "json"
)

// Connect connects and identifies the client to the Discord Gateway.
func (c *Client) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isConnected() {
		return ErrAlreadyConnected
	}

	c.connecting.Store(true)
	defer c.connecting.Store(false)

	var err error
	// Get the Gateway endpoint if we don't have one cached yet.
	if c.gatewayURL == "" {
		// NOTE: not using GatewayBot here because a Client has no
		// notion of automatic sharding. This is handled at a higher level,
		// when creating a Client with the WithSharding option.
		c.gatewayURL, err = c.Gateway(ctx)
		if err != nil {
			return fmt.Errorf("could not get gateway URL: %w", err)
		}
	}

	// Those fields' lifecycle is tied to a connection, not to the Client,
	// so we need to initialize them each time we attempt a new connection.
	c.voicePayloads = make(chan *payload.Payload)
	c.error = make(chan error)
	c.reportErrorOnce = sync.Once{}
	c.stop = make(chan struct{})

	// This context is bound to the Gateway connection and will be
	// canceled when it is closed.
	c.ctx, c.cancel = context.WithCancel(context.Background())

	// Open the Gateway websocket connection.
	header := make(http.Header)
	header.Add("Accept-Encoding", "zlib")
	gwURL := fmt.Sprintf("%s?v=%d&encoding=%s", c.gatewayURL, gatewayVersion, gatewayEncoding)
	c.logger.Debugf("connecting to the gateway: %s", gwURL)
	c.conn, _, err = websocket.Dial(ctx, gwURL, &websocket.DialOptions{HTTPHeader: header})
	if err != nil {
		return err
	}

	// If any error occurs during the connection process, we
	// should close the underlying websocket connection, so
	// we can try to reconnect later. We should also signal
	// to already started goroutines to stop by closing the
	// stop channel to prevent them from leaking, mark this
	// client as not connected and cancel the connection
	// context.
	defer func() {
		if err != nil {
			_ = c.conn.Close(websocket.StatusInternalError, "failed to establish connection") // Not much we can do about this, maybe log it?
			c.connected.Store(false)
			close(c.stop)
			c.cancel()
		}
	}()

	// The Gateway should first send us a Hello packet defining the heartbeat
	// interval we must use when we connect to the websocket.
	p, err := c.recvPayload()
	if err != nil {
		return fmt.Errorf("could not receive payload from gateway: %w", err)
	}
	if p.Op != 10 {
		return fmt.Errorf("expected Opcode 10 Hello; got Opcode %d", p.Op)
	}

	var hello struct {
		HeartbeatInterval int `json:"heartbeat_interval"`
	}
	if err = json.Unmarshal(p.D, &hello); err != nil {
		return err
	}

	// If the sequence number is 0 and we don't have a
	// session ID, we must identify to the Gateway to
	// create a new session, else this means we have already
	// been connected to the Gateway with this client and
	// we should try to resume a previous connection.
	seq := c.sequence.Load()
	if seq == 0 && c.sessionID == "" {
		c.logger.Debug("identifying to the gateway")
		if err = c.identify(ctx); err != nil {
			return err
		}

		// The Gateway should send us a Ready event if we successfully authenticated.
		if err = c.ready(); err != nil {
			return err
		}
	} else {
		c.logger.Debugf("trying to resume an existing session (seq=%d; sessID=%q)", seq, c.sessionID)
		if err = c.resume(ctx); err != nil {
			return err
		}
		// The Gateway should replay events we missed since we were disconnected
		// and then send us a Resumed payload. All of this is handled by the
		// listenAndHandlePayloads goroutine.
	}

	// From now, we are connected to the Gateway.
	// Start the connection manager, heartbeating
	// and listening for Gateway events.
	c.wg.Add(1)
	go c.wait()

	c.wg.Add(1)
	go c.heartbeat(time.Duration(hello.HeartbeatInterval) * time.Millisecond)

	c.wg.Add(1)
	go c.listenAndHandlePayloads()

	return nil
}

// Disconnect closes the connection to the Discord Gateway.
func (c *Client) Disconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// No-op if we're already disconnected and not trying to reconnect.
	if !c.isConnected() && !c.isReconnecting() {
		return
	}

	// First, try to properly leave all voice channel we are connected to.
	// NOTE: maybe adjust this timeout to the number of voice connections
	// we have. Something like min 10s, max 120s with 1s/conn.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	for guildID := range c.voiceConnections {
		wg.Add(1)

		go func(guildID string) {
			defer wg.Done()

			if err := c.LeaveVoiceChannel(ctx, guildID); err != nil {
				c.logger.Errorf("could not properly disconnect from voice channel: %v", err)
			}
		}(guildID)
	}

	// Wait for all voice connections to be closed.
	wg.Wait()

	// Then, signal the connection manager that we want to disconnect.
	close(c.stop)
	// Properly wait for all goroutines to exit.
	c.wg.Wait()
}

// wait waits for an error to happen while connected to the Gateway
// or for a stop signal to be sent.
// If an unexpected error happens while connected to the
// Gateway, this method will automatically try to reconnect.
func (c *Client) wait() {
	defer c.wg.Done()

	c.logger.Debug("starting gateway connection manager")
	defer c.logger.Debug("stopped gateway connection manager")

	var err error
	select {
	// An unexpected error occurred while communicating with the Gateway.
	case err = <-c.error:
		c.onGatewayError(err)

	// User called Client.Disconnect.
	case <-c.stop:
		c.logger.Debug("disconnecting from the gateway")
		c.onDisconnect()
	}

	close(c.voicePayloads)

	c.cancel()
	c.connected.Store(false)

	// If there was an error, try to reconnect depending on its code.
	if shouldReconnect(err) {
		c.reconnectWithBackoff()
	}
}

// Determine whether we should try to reconnect based on the error we got.
// See https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-close-event-codes for more information.
func shouldReconnect(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, errMustReconnect) {
		return true
	}

	switch websocket.CloseStatus(err) {
	case 4001, 4002, 4003, 4004, 4005, 4010, 4011:
		return false
	case 4000, 4007, 4008, 4009:
		return true
	case -1: // Not a websocket error.
		return true
	default: // New (or undocumented?) close status code.
		return true
	}
}

// reconnectWithBackoff attempts to reconnect to the Gateway using the Client's
// backoff strategy.
func (c *Client) reconnectWithBackoff() {
	c.reconnecting.Store(true)
	defer c.reconnecting.Store(false)

	c.logger.Debug("trying to reconnect to the gateway")

	for i := 0; true; i++ {
		// Try to establish a new connection with a 30 seconds timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		if err := c.Connect(ctx); err != nil {
			cancel()

			if !shouldReconnect(err) {
				c.logger.Error("invalid Gateway session, can not recover: %v", err)
				return
			}

			duration := c.backoff.forAttempt(i)
			c.logger.Errorf("failed to reconnect: %v, retrying in %s", err, duration)

			select {
			case <-time.After(duration):
				continue // Make a new connection attempt.
			case <-c.stop:
				// Client called Disconnect(), stop trying to reconnect.
				c.logger.Info("client called Disconnect while trying to reconnect to the gateway, aborting")
				return
			}
		} else {
			cancel()

			// We could reconnect.
			c.logger.Info("successfully reconnected to the gateway")
			return
		}
	}
}

// reportErr reports the first fatal error encountered while connected to
// the Gateway. Calls after the first one are no-ops.
func (c *Client) reportErr(err error) {
	c.reportErrorOnce.Do(func() {
		select {
		case c.error <- err:

		// Discard the error if we are already closing the connection.
		case <-c.stop:
		}
		close(c.error)
	})
}

// onGatewayError is called when an error occurs while the connection to
// the Gateway is up. It closes the underlying websocket connection
// with a 1006 code, logs the error and finally signals to all other
// goroutines (heartbeat, listenAndHandlePayloads, etc.) to stop by
// closing the stop channel.
func (c *Client) onGatewayError(err error) {
	c.logger.Errorf("gateway connection error: %v", err)

	if closeErr := c.conn.Close(websocket.StatusInternalError, "gateway error"); closeErr != nil {
		c.logger.Errorf("could not properly close websocket connection (error): %v", closeErr)
	}

	// If an error occurred while we are establishing the connection,
	// the stop channel will already be closed, so return early.
	if c.isConnecting() {
		return
	}

	close(c.stop)
}

// onDisconnect is called when a normal disconnection happens (the client
// called the Disconnect() method). It closes the underlying websocket
// connection with a 1000 code and resets the session of this Client
// so it can open a new fresh connection by calling Connect() again.
func (c *Client) onDisconnect() {
	if err := c.conn.Close(websocket.StatusNormalClosure, "disconnecting"); err != nil {
		c.logger.Errorf("could not properly close websocket connection: %v", err)
	}
	// Reset the Gateway session so the user gets a new
	// fresh session if reconnecting with the same Client.
	c.resetGatewaySession()
}

// isConnected reports whether the client is currently connected to the Gateway.
func (c *Client) isConnected() bool {
	return c.connected.Load()
}

// isConnecting reports whether the client is currently connecting to the Gateway.
func (c *Client) isConnecting() bool {
	return c.connecting.Load()
}

// isConnectingToVoice reports whether the client is currently connecting to
// a voice server.
func (c *Client) isConnectingToVoice() bool {
	return c.connectingToVoice.Load()
}

// isReconnecting reports whether the client is currently reconnecting to the Gateway.
func (c *Client) isReconnecting() bool {
	return c.reconnecting.Load()
}

// resetGatewaySession resets the session ID as well as the sequence number
// of the Gateway connection.
// After a session reset, a call to Connect will send an Identify payload and
// start a new fresh session, instead of trying to resume an existing session.
func (c *Client) resetGatewaySession() {
	c.sessionID = ""
	c.sequence.Store(0)
}
