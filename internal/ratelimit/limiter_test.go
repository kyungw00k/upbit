package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestNewLimiter(t *testing.T) {
	l := NewLimiter()
	if l == nil {
		t.Fatal("NewLimiter가 nil을 반환")
	}

	// 모든 그룹에 대한 버킷이 존재하는지 확인
	groups := []Group{
		GroupMarket, GroupCandle, GroupTicker, GroupOrderbook,
		GroupTrade, GroupDefault, GroupOrder, GroupOrderTest, GroupOrderCancelAll,
	}
	for _, g := range groups {
		if _, ok := l.buckets[g]; !ok {
			t.Errorf("그룹 %s에 대한 버킷이 없음", g)
		}
	}
}

func TestWait_BasicOperation(t *testing.T) {
	l := NewLimiter()

	ctx := context.Background()
	err := l.Wait(ctx, GroupDefault)
	if err != nil {
		t.Fatalf("Wait 실패: %v", err)
	}
}

func TestWait_CancelledContext(t *testing.T) {
	l := NewLimiter()

	// 버킷 토큰을 모두 소진
	b := l.buckets[GroupOrderCancelAll]
	b.mu.Lock()
	b.tokens = 0
	b.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 즉시 취소

	err := l.Wait(ctx, GroupOrderCancelAll)
	if err == nil {
		t.Error("취소된 context에서 에러가 반환되어야 함")
	}
}

func TestWait_UnknownGroup_FallsBackToDefault(t *testing.T) {
	l := NewLimiter()

	ctx := context.Background()
	err := l.Wait(ctx, Group("unknown"))
	if err != nil {
		t.Fatalf("알 수 없는 그룹에서 default로 폴백해야 함: %v", err)
	}
}

func TestGroupFromPath(t *testing.T) {
	tests := []struct {
		path     string
		expected Group
	}{
		{"/candles/minutes/1", GroupCandle},
		{"/candles/days", GroupCandle},
		{"/orders", GroupOrder},
		{"/orders/chance", GroupOrder},
		{"/orders/test", GroupOrderTest},
		{"/orders/batch", GroupOrderCancelAll},
		{"/tickers", GroupTicker},
		{"/ticker", GroupTicker},
		{"/orderbooks", GroupOrderbook},
		{"/trades/ticks", GroupTrade},
		{"/trading_pairs", GroupMarket},
		{"/accounts", GroupDefault},
		{"/deposits", GroupDefault},
		{"/withdraws", GroupDefault},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := GroupFromPath(tt.path)
			if got != tt.expected {
				t.Errorf("GroupFromPath(%q) = %s, 기대: %s", tt.path, got, tt.expected)
			}
		})
	}
}

func TestGroupValidate(t *testing.T) {
	// 유효한 그룹
	validGroups := []Group{
		GroupMarket, GroupCandle, GroupTicker, GroupOrderbook,
		GroupTrade, GroupDefault, GroupOrder, GroupOrderTest, GroupOrderCancelAll,
	}
	for _, g := range validGroups {
		if err := g.Validate(); err != nil {
			t.Errorf("유효한 그룹 %s에서 에러: %v", g, err)
		}
	}

	// 유효하지 않은 그룹
	if err := Group("invalid").Validate(); err == nil {
		t.Error("유효하지 않은 그룹에서 에러가 반환되어야 함")
	}
}

func TestUpdateFromHeader(t *testing.T) {
	l := NewLimiter()

	// sec=2 이하이면 긴급 제한 (rate=1.0)
	l.UpdateFromHeader("group=market; min=599; sec=1", GroupMarket)
	b := l.buckets[GroupMarket]
	b.mu.Lock()
	rate := b.rate
	b.mu.Unlock()
	if rate != 1.0 {
		t.Errorf("sec=1일 때 rate=1.0 기대, 실제: %f", rate)
	}

	// sec > 2이면 해당 값으로 조절
	l2 := NewLimiter()
	l2.UpdateFromHeader("group=ticker; min=500; sec=8", GroupTicker)
	b2 := l2.buckets[GroupTicker]
	b2.mu.Lock()
	rate2 := b2.rate
	b2.mu.Unlock()
	if rate2 != 8.0 {
		t.Errorf("sec=8일 때 rate=8.0 기대, 실제: %f", rate2)
	}
}

func TestUpdateFromHeader_Empty(t *testing.T) {
	l := NewLimiter()
	// 빈 헤더는 아무 변경도 하지 않음 (패닉 없이)
	l.UpdateFromHeader("", GroupMarket)
}

func TestBucketWait_Timeout(t *testing.T) {
	l := NewLimiter()

	// OrderCancelAll은 rate=0.5 (2초에 1회)
	b := l.buckets[GroupOrderCancelAll]
	b.mu.Lock()
	b.tokens = 0
	b.lastTime = time.Now()
	b.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := l.Wait(ctx, GroupOrderCancelAll)
	if err == nil {
		t.Error("타임아웃 시 에러가 반환되어야 함")
	}
}
