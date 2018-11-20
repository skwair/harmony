package harmony

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const (
	gatewayVersion  = 6
	gatewayEncoding = "json"
)

var (
	// ErrAlreadyConnected is returned by Connect when a connection to the Gateway already exists.
	ErrAlreadyConnected = errors.New("already connected to the Gateway")
)

// Connect connects and identifies the client to the Discord gateway.
func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isConnected() {
		return ErrAlreadyConnected
	}

	var err error
	// Get the Gateway endpoint if we don't have one cached yet.
	if c.gatewayURL == "" {
		// NOTE: not using GatewayBot here because a Client has no
		// notion of automatic sharding. This is handled at a higher level,
		// when creating a Client with the WithSharding option.
		c.gatewayURL, err = c.Gateway()
		if err != nil {
			return err
		}
	}

	c.voicePayloads = make(chan *payload)
	c.error = make(chan error)
	c.stop = make(chan struct{})

	header := http.Header{}
	header.Add("Accept-Encoding", "zlib")
	gwURL := fmt.Sprintf("%s?v=%d&encoding=%s", c.gatewayURL, gatewayVersion, gatewayEncoding)
	c.conn, _, err = websocket.DefaultDialer.Dial(gwURL, header)
	if err != nil {
		return err
	}

	// If any error occurs during the connection process, we
	// should close the underlying websocket connection, so
	// we can try to reconnect later. We should also signal
	// to already started goroutines to stop by closing the
	// stop channel to prevent them from leaking and mark
	// this client as not connected.
	defer func() {
		if err != nil {
			c.conn.Close()
			atomic.StoreInt32(&c.connected, 0)
			close(c.stop)
		}
	}()

	// The Gateway should send us a Hello packet defining the heartbeat
	// interval when we connect to the websocket.
	p, err := c.recvPayload()
	if err != nil {
		return fmt.Errorf("could not receive payload from gateway: %v", err)
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
	// create a new session, else we should try to resume it.
	seq := atomic.LoadInt64(&c.sequence)
	if seq == 0 && c.sessionID == "" {
		if err = c.identify(); err != nil {
			return err
		}
		// The Gateway should send us a Ready event if we successfully authenticated.
		if err = c.ready(); err != nil {
			return err
		}
	} else {
		if err = c.resume(); err != nil {
			return err
		}
		// The Gateway should replay events we missed since we were disconnected
		// and then send us a Resumed payload. All of this is handled by the listen
		// goroutine.
		// NOTE: maybe we should reconnect to voice if we had active connections here.
	}

	// From now, we are connected to the Gateway.
	// Start heartbeating and listening for Gateway events.
	c.wg.Add(3) // listen starts an additional goroutine.
	go c.heartbeat(time.Duration(hello.HeartbeatInterval) * time.Millisecond)
	go c.listen()

	go c.wait() // wait does not count in the waitgroup.

	return nil
}

// Disconnect closes the connection to the Discord Gateway.
func (c *Client) Disconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected() {
		return
	}

	close(c.stop)
	c.wg.Wait()
}

// wait waits for an error or a stop signal to be sent.
func (c *Client) wait() {
	var err error
	select {
	// An error occurred while communicating with the Gateway.
	case err = <-c.error:
		c.onGatewayError(err)

	case <-c.stop:
		c.onDisconnect()
	}

	c.conn.Close()
	atomic.StoreInt32(&c.connected, 0)

	// If there was an error, try to reconnect.
	if err != nil {
		for i := 0; true; i++ {
			if err = c.Connect(); err != nil {
				duration := c.backoff.forAttempt(i)
				c.errorHandler(fmt.Errorf("failed to reconnect: %v, retrying in %s", err, duration))
				select {
				case <-time.After(duration):
				case <-c.stop:
					// Client called Disconnect(), stop trying to reconnect.
					return
				}
			} else {
				// We could reconnect.
				return
			}
		}
	}
}

// onGatewayError is called when an error occurs while the connection to
// the Gateway is up. It closes the underlying websocket connection
// with a 1006 code, calls the registered error handler and finally
// signals to all other goroutines (heartbeat, listen, etc.) to stop.
func (c *Client) onGatewayError(err error) {
	if err := c.conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseAbnormalClosure, ""),
		time.Now().Add(time.Second*10),
	); err != nil {
		c.errorHandler(fmt.Errorf("could not properly close websocket: %v", err))
	}
	c.errorHandler(err)
	close(c.stop)
}

// onDisconnect is called when a normal disconnection happens (the client
// called the Disconnect() method). It closes the underlying websocket
// connection with a 1000 code and resets the session of this Client
// so it can open a new fresh connection by calling Connect() again.
func (c *Client) onDisconnect() {
	if err := c.conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(time.Second*10),
	); err != nil {
		c.errorHandler(fmt.Errorf("could not properly close websocket: %v", err))
	}
	// Reset the sequence number and the session ID so
	// the user gets a new fresh session if reconnecting
	// with the same client.
	atomic.StoreInt64(&c.sequence, 0)
	c.sessionID = ""
}

// isConnected reports whether the client is currently connected to the Gateway.
func (c *Client) isConnected() bool {
	return atomic.LoadInt32(&c.connected) == 1
}

// isConnectingToVoice reports whether the client is currently connecting to
// a voice server.
func (c *Client) isConnectingToVoice() bool {
	return atomic.LoadInt32(&c.connectingToVoice) == 1
}
