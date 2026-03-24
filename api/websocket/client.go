// Package websocket provides a WebSocket client for Upbit's real-time streaming API.
// See https://docs.upbit.com/reference/websocket for API documentation.
package websocket

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/kyungw00k/upbit/api"
)

const (
	// PublicURL is the WebSocket endpoint for market data (Quotation).
	PublicURL = "wss://api.upbit.com/websocket/v1"
	// PrivateURL is the WebSocket endpoint for account and order data (Exchange).
	PrivateURL = "wss://api.upbit.com/websocket/v1/private"

	defaultPingInterval = 30 * time.Second
	defaultPongWait     = 10 * time.Second
	defaultMaxRetries   = 5

	// readWait accounts for Upbit's 120-second idle timeout with headroom.
	readWait = 120 * time.Second
)

// Option is a functional option for WSClient.
type Option func(*WSClient)

// WithAuth sets authentication credentials for private channel access.
func WithAuth(accessKey, secretKey string) Option {
	return func(c *WSClient) {
		c.accessKey = accessKey
		c.secretKey = secretKey
	}
}

// WithPingInterval sets the interval between ping messages.
func WithPingInterval(d time.Duration) Option {
	return func(c *WSClient) {
		c.pingInterval = d
	}
}

// WithMaxRetries sets the maximum number of reconnection attempts.
func WithMaxRetries(n int) Option {
	return func(c *WSClient) {
		c.maxRetries = n
	}
}

// WSClient is a WebSocket client for Upbit streaming.
type WSClient struct {
	url          string
	accessKey    string
	secretKey    string
	pingInterval time.Duration
	maxRetries   int

	conn      *websocket.Conn
	mu        sync.Mutex
	closed    bool
	closeCh   chan struct{}
	closeOnce sync.Once
	pingDone  chan struct{}
	pingWg    sync.WaitGroup

	// lastSubMsg stores the last subscription message for auto-reconnect resubscription.
	lastSubMsg []byte
}

// NewWSClient creates a new WSClient with the given URL and options.
func NewWSClient(url string, opts ...Option) *WSClient {
	c := &WSClient{
		url:          url,
		pingInterval: defaultPingInterval,
		maxRetries:   defaultMaxRetries,
		closeCh:      make(chan struct{}),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Connect establishes a WebSocket connection.
func (c *WSClient) Connect(ctx context.Context) error {
	return c.connectWithRetry(ctx)
}

// connectWithRetry reconnects with exponential backoff.
func (c *WSClient) connectWithRetry(ctx context.Context) error {
	var lastErr error
	for i := 0; i <= c.maxRetries; i++ {
		if i > 0 {
			// exponential backoff: 1s, 2s, 4s...
			backoff := time.Duration(math.Pow(2, float64(i-1))) * time.Second
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}

		err := c.dial(ctx)
		if err == nil {
			return nil
		}
		lastErr = err
	}
	return fmt.Errorf("WebSocket connection failed after %d attempts: %w", c.maxRetries+1, lastErr)
}

// dial establishes the underlying WebSocket connection.
func (c *WSClient) dial(ctx context.Context) error {
	dialer := websocket.DefaultDialer
	header := http.Header{}

	// Add authorization header for private channel access.
	if c.accessKey != "" && c.secretKey != "" {
		token, err := api.GenerateToken(c.accessKey, c.secretKey, nil)
		if err != nil {
			return fmt.Errorf("failed to generate JWT token: %w", err)
		}
		header.Set("Authorization", "Bearer "+token)
	}

	conn, _, err := dialer.DialContext(ctx, c.url, header)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.conn = conn
	// S-1: Reset closeCh, closeOnce, and closed flag on reconnect to prevent ping
	// interruption and double-close panics.
	c.closeCh = make(chan struct{})
	c.closeOnce = sync.Once{}
	c.closed = false
	// C-1: Recreate pingDone channel on each reconnect to prevent close panics.
	c.pingDone = make(chan struct{})
	c.mu.Unlock()

	// C-3: Set initial ReadDeadline to detect the 120-second idle timeout.
	_ = conn.SetReadDeadline(time.Now().Add(readWait))

	// Pong handler: renew ReadDeadline on each pong received.
	conn.SetPongHandler(func(appData string) error {
		return conn.SetReadDeadline(time.Now().Add(readWait))
	})

	// W-3: Track pingLoop goroutine with WaitGroup.
	c.pingWg.Add(1)
	go c.pingLoop()

	return nil
}

// pingLoop sends periodic ping messages to keep the connection alive.
func (c *WSClient) pingLoop() {
	defer c.pingWg.Done()

	ticker := time.NewTicker(c.pingInterval)
	defer ticker.Stop()

	c.mu.Lock()
	done := c.pingDone
	c.mu.Unlock()

	defer close(done)

	for {
		select {
		case <-c.closeCh:
			return
		case <-ticker.C:
			c.mu.Lock()
			if c.conn == nil || c.closed {
				c.mu.Unlock()
				return
			}
			err := c.conn.WriteControl(
				websocket.PingMessage,
				nil,
				time.Now().Add(defaultPongWait),
			)
			c.mu.Unlock()
			if err != nil {
				return
			}
		}
	}
}

// Subscribe sends a subscription message and stores it for reconnect resubscription.
func (c *WSClient) Subscribe(msg []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return fmt.Errorf("WebSocket is not connected")
	}

	// Store the message for resubscription on reconnect.
	c.lastSubMsg = make([]byte, len(msg))
	copy(c.lastSubMsg, msg)

	return c.conn.WriteMessage(websocket.TextMessage, msg)
}

// ReadMessage receives a message from the WebSocket connection.
// C-2: On error, attempts reconnection and resubscription before resuming reads.
func (c *WSClient) ReadMessage() (messageType int, data []byte, err error) {
	// W-4: Read conn under mutex protection.
	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()

	if conn == nil {
		return 0, nil, fmt.Errorf("WebSocket is not connected")
	}
	return conn.ReadMessage()
}

// ReadMessageWithReconnect receives messages with automatic reconnection.
// Returns immediately if ctx is cancelled. On connection error, reconnects and
// resubscribes before continuing to receive messages.
func (c *WSClient) ReadMessageWithReconnect(ctx context.Context) (int, []byte, error) {
	for {
		if ctx.Err() != nil {
			return 0, nil, ctx.Err()
		}

		c.mu.Lock()
		conn := c.conn
		c.mu.Unlock()

		if conn == nil {
			return 0, nil, fmt.Errorf("WebSocket is not connected")
		}

		mt, data, err := conn.ReadMessage()
		if err == nil {
			return mt, data, nil
		}

		// Context cancelled by the caller.
		if ctx.Err() != nil {
			return 0, nil, ctx.Err()
		}

		// Do not reconnect on normal closure.
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			return 0, nil, err
		}

		// Do not reconnect if Close() was explicitly called.
		c.mu.Lock()
		if c.closed {
			c.mu.Unlock()
			return 0, nil, err
		}
		c.mu.Unlock()

		// Attempt reconnection.
		reconnErr := c.reconnect(ctx)
		if reconnErr != nil {
			return 0, nil, fmt.Errorf("reconnect failed: %w (original error: %v)", reconnErr, err)
		}

		// Reconnect succeeded — loop back to read again.
	}
}

// reconnect closes the current connection and establishes a new one with resubscription.
func (c *WSClient) reconnect(ctx context.Context) error {
	// Clean up the existing connection.
	c.cleanupConn()

	// Reconnect.
	if err := c.connectWithRetry(ctx); err != nil {
		return err
	}

	// Resubscribe.
	c.mu.Lock()
	subMsg := c.lastSubMsg
	c.mu.Unlock()

	if len(subMsg) > 0 {
		c.mu.Lock()
		conn := c.conn
		c.mu.Unlock()
		if conn != nil {
			if err := conn.WriteMessage(websocket.TextMessage, subMsg); err != nil {
				return fmt.Errorf("resubscribe failed: %w", err)
			}
		}
	}

	return nil
}

// cleanupConn closes the current connection and waits for pingLoop to exit,
// without setting the closed flag (unlike Close).
func (c *WSClient) cleanupConn() {
	c.mu.Lock()
	conn := c.conn
	c.conn = nil
	c.mu.Unlock()

	if conn != nil {
		_ = conn.Close()
	}

	// W-3: Wait for pingLoop goroutine to exit.
	c.pingWg.Wait()
}

// Close gracefully shuts down the WebSocket connection.
func (c *WSClient) Close() error {
	var closeErr error

	c.closeOnce.Do(func() {
		c.mu.Lock()
		c.closed = true
		close(c.closeCh)
		conn := c.conn
		c.conn = nil
		c.mu.Unlock()

		// W-3: Wait for pingLoop goroutine to exit.
		c.pingWg.Wait()

		if conn != nil {
			// Graceful close: send CloseMessage before closing the connection.
			_ = conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			closeErr = conn.Close()
		}
	})

	return closeErr
}
