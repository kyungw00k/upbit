package websocket

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/kyungw00k/upbit/internal/api"
)

const (
	// PublicURL 시세(Quotation) WebSocket 엔드포인트
	PublicURL = "wss://api.upbit.com/websocket/v1"
	// PrivateURL 자산 및 주문(Exchange) WebSocket 엔드포인트
	PrivateURL = "wss://api.upbit.com/websocket/v1/private"

	defaultPingInterval = 30 * time.Second
	defaultPongWait     = 10 * time.Second
	defaultMaxRetries   = 5

	// Upbit WebSocket 120초 idle timeout 대비 여유분 포함
	readWait = 120 * time.Second
)

// Option WSClient 옵션 함수
type Option func(*WSClient)

// WithAuth 인증 정보 설정 (Private 채널용)
func WithAuth(accessKey, secretKey string) Option {
	return func(c *WSClient) {
		c.accessKey = accessKey
		c.secretKey = secretKey
	}
}

// WithPingInterval ping 전송 주기 설정
func WithPingInterval(d time.Duration) Option {
	return func(c *WSClient) {
		c.pingInterval = d
	}
}

// WithMaxRetries 최대 재연결 시도 횟수 설정
func WithMaxRetries(n int) Option {
	return func(c *WSClient) {
		c.maxRetries = n
	}
}

// WSClient WebSocket 클라이언트
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

	// 자동 재연결용: 마지막 구독 메시지 보관
	lastSubMsg []byte
}

// NewWSClient WSClient 생성
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

// Connect WebSocket 연결
func (c *WSClient) Connect(ctx context.Context) error {
	return c.connectWithRetry(ctx)
}

// connectWithRetry exponential backoff 재연결
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
	return fmt.Errorf("WebSocket 연결 실패 (%d회 시도): %w", c.maxRetries+1, lastErr)
}

// dial 실제 WebSocket 연결 수립
func (c *WSClient) dial(ctx context.Context) error {
	dialer := websocket.DefaultDialer
	header := http.Header{}

	// Private 채널 인증 헤더 추가
	if c.accessKey != "" && c.secretKey != "" {
		token, err := api.GenerateToken(c.accessKey, c.secretKey, nil)
		if err != nil {
			return fmt.Errorf("JWT 토큰 생성 실패: %w", err)
		}
		header.Set("Authorization", "Bearer "+token)
	}

	conn, _, err := dialer.DialContext(ctx, c.url, header)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.conn = conn
	// S-1: 재연결 시 closeCh, closeOnce, closed 플래그를 리셋하여 ping 중단/double close panic 방지
	c.closeCh = make(chan struct{})
	c.closeOnce = sync.Once{}
	c.closed = false
	// C-1: 재연결 시 pingDone 채널을 매번 새로 생성하여 close panic 방지
	c.pingDone = make(chan struct{})
	c.mu.Unlock()

	// C-3: 초기 ReadDeadline 설정 — 120초 idle timeout 감지
	_ = conn.SetReadDeadline(time.Now().Add(readWait))

	// pong 핸들러: pong 수신 시 ReadDeadline 갱신
	conn.SetPongHandler(func(appData string) error {
		return conn.SetReadDeadline(time.Now().Add(readWait))
	})

	// W-3: WaitGroup으로 pingLoop goroutine 추적
	c.pingWg.Add(1)
	go c.pingLoop()

	return nil
}

// pingLoop 주기적 ping 전송
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

// Subscribe 구독 메시지 전송
func (c *WSClient) Subscribe(msg []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return fmt.Errorf("WebSocket이 연결되지 않았습니다")
	}

	// 재연결 시 재구독을 위해 메시지 보관
	c.lastSubMsg = make([]byte, len(msg))
	copy(c.lastSubMsg, msg)

	return c.conn.WriteMessage(websocket.TextMessage, msg)
}

// ReadMessage 메시지 수신 (자동 재연결 포함)
// C-2: 에러 발생 시 재연결을 시도하고, 재연결 성공 시 재구독 후 계속 수신
func (c *WSClient) ReadMessage() (messageType int, data []byte, err error) {
	// W-4: conn을 mutex로 보호하여 읽기
	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()

	if conn == nil {
		return 0, nil, fmt.Errorf("WebSocket이 연결되지 않았습니다")
	}
	return conn.ReadMessage()
}

// ReadMessageWithReconnect 자동 재연결이 포함된 메시지 수신
// ctx가 취소되면 즉시 반환. 연결 에러 시 재연결 후 재구독하고 계속 수신 시도.
func (c *WSClient) ReadMessageWithReconnect(ctx context.Context) (int, []byte, error) {
	for {
		if ctx.Err() != nil {
			return 0, nil, ctx.Err()
		}

		c.mu.Lock()
		conn := c.conn
		c.mu.Unlock()

		if conn == nil {
			return 0, nil, fmt.Errorf("WebSocket이 연결되지 않았습니다")
		}

		mt, data, err := conn.ReadMessage()
		if err == nil {
			return mt, data, nil
		}

		// 사용자가 종료를 요청한 경우
		if ctx.Err() != nil {
			return 0, nil, ctx.Err()
		}

		// 정상 종료 메시지인 경우 재연결하지 않음
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			return 0, nil, err
		}

		// 사용자가 명시적으로 Close()를 호출한 경우
		c.mu.Lock()
		if c.closed {
			c.mu.Unlock()
			return 0, nil, err
		}
		c.mu.Unlock()

		// 재연결 시도
		reconnErr := c.reconnect(ctx)
		if reconnErr != nil {
			return 0, nil, fmt.Errorf("재연결 실패: %w (원본 에러: %v)", reconnErr, err)
		}

		// 재연결 성공 — 루프 상단으로 돌아가 다시 ReadMessage
	}
}

// reconnect 기존 연결을 정리하고 새로 연결 + 재구독
func (c *WSClient) reconnect(ctx context.Context) error {
	// 기존 연결 정리
	c.cleanupConn()

	// 재연결
	if err := c.connectWithRetry(ctx); err != nil {
		return err
	}

	// 재구독
	c.mu.Lock()
	subMsg := c.lastSubMsg
	c.mu.Unlock()

	if len(subMsg) > 0 {
		c.mu.Lock()
		conn := c.conn
		c.mu.Unlock()
		if conn != nil {
			if err := conn.WriteMessage(websocket.TextMessage, subMsg); err != nil {
				return fmt.Errorf("재구독 실패: %w", err)
			}
		}
	}

	return nil
}

// cleanupConn 기존 연결과 pingLoop를 정리 (Close와 달리 closed 플래그를 세팅하지 않음)
func (c *WSClient) cleanupConn() {
	c.mu.Lock()
	conn := c.conn
	c.conn = nil
	c.mu.Unlock()

	if conn != nil {
		_ = conn.Close()
	}

	// W-3: pingLoop goroutine 종료 대기
	c.pingWg.Wait()
}

// Close 연결 종료
func (c *WSClient) Close() error {
	var closeErr error

	c.closeOnce.Do(func() {
		c.mu.Lock()
		c.closed = true
		close(c.closeCh)
		conn := c.conn
		c.conn = nil
		c.mu.Unlock()

		// W-3: pingLoop goroutine 종료 대기
		c.pingWg.Wait()

		if conn != nil {
			// graceful close: CloseMessage 전송 후 연결 닫기
			_ = conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			closeErr = conn.Close()
		}
	})

	return closeErr
}
