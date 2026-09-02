package websocket

import (
	"context"
	"testing"
	"time"
)

func TestMsgLimiter_PerSecBudget(t *testing.T) {
	l := newMsgLimiter(5, 0)

	for i := 0; i < 5; i++ {
		if !l.allow() {
			t.Fatalf("버스트 내 %d번째 메시지가 거부됨", i+1)
		}
	}
	if l.allow() {
		t.Error("초당 5회 초과 메시지가 허용됨")
	}
}

func TestMsgLimiter_PerMinBudget(t *testing.T) {
	l := newMsgLimiter(0, 3)

	for i := 0; i < 3; i++ {
		if !l.allow() {
			t.Fatalf("분당 한도 내 %d번째 메시지가 거부됨", i+1)
		}
	}
	if l.allow() {
		t.Error("분당 3회 초과 메시지가 허용됨")
	}
}

func TestMsgLimiter_Refill(t *testing.T) {
	l := newMsgLimiter(1, 0)

	if !l.allow() {
		t.Fatal("첫 메시지가 거부됨")
	}
	if l.allow() {
		t.Fatal("초당 1회 초과 허용됨")
	}

	// 1초 + 여유 후 토큰 1개 리필
	time.Sleep(1100 * time.Millisecond)
	if !l.allow() {
		t.Error("리필 후 메시지가 거부됨")
	}
}

func TestMsgLimiter_Reset(t *testing.T) {
	l := newMsgLimiter(2, 2)

	l.allow()
	l.allow()
	if l.allow() {
		t.Fatal("예산 소진 후 허용됨")
	}

	l.reset()
	for i := 0; i < 2; i++ {
		if !l.allow() {
			t.Fatalf("reset 후 %d번째 메시지가 거부됨", i+1)
		}
	}
}

func TestMsgLimiter_WaitCancelledContext(t *testing.T) {
	l := newMsgLimiter(1, 0)
	l.allow() // 예산 소진

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	if err := l.wait(ctx); ctx.Err() == nil || err == nil {
		t.Errorf("취소된 컨텍스트에서 wait가 에러를 반환해야 함: err=%v ctxErr=%v", err, ctx.Err())
	}
}

func TestNewWSClient_DefaultsSendLimiter(t *testing.T) {
	c := NewWSClient(PublicURL)
	if c.sendLimiter == nil {
		t.Fatal("sendLimiter가 초기화되지 않음")
	}
	if c.sendLimiter.perSec != maxMessagesPerSec {
		t.Errorf("초당 한도 = %f, 기대: %f", c.sendLimiter.perSec, maxMessagesPerSec)
	}
	if c.sendLimiter.perMin != maxMessagesPerMin {
		t.Errorf("분당 한도 = %f, 기대: %f", c.sendLimiter.perMin, maxMessagesPerMin)
	}
}
