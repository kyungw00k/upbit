package tui

import (
	"strings"
	"testing"
)

// --- charForCell tests ---

func TestCharForCell_Void(t *testing.T) {
	// 캔들 범위(lowSR=5, highSR=10) 밖의 서브 로우 → charVoid
	ch := charForCell(3, 2, 10, 5, 7, 6)
	if ch != charVoid {
		t.Errorf("expected charVoid, got %q", ch)
	}
}

func TestCharForCell_Body(t *testing.T) {
	// 두 서브 로우 모두 몸통 내부 → charBody
	// bodyBot=4, bodyTop=8, highSR=10, lowSR=2
	ch := charForCell(7, 6, 10, 2, 8, 4)
	if ch != charBody {
		t.Errorf("expected charBody, got %q", ch)
	}
}

func TestCharForCell_Wick(t *testing.T) {
	// 두 서브 로우 모두 심지(wick) 영역 → charWick
	// highSR=10, lowSR=0, bodyBot=5, bodyTop=7
	// topSub=9, botSub=8: 둘 다 wick 영역(7 < 8 <= 10)
	ch := charForCell(9, 8, 10, 0, 7, 5)
	if ch != charWick {
		t.Errorf("expected charWick, got %q", ch)
	}
}

func TestCharForCell_HalfBodyTop(t *testing.T) {
	// topSub가 몸통, botSub가 범위 밖 → charHalfBodyTop
	// highSR=10, lowSR=5, bodyBot=6, bodyTop=9
	// topSub=7 (body), botSub=4 (range 밖)
	ch := charForCell(7, 4, 10, 5, 9, 6)
	if ch != charHalfBodyTop {
		t.Errorf("expected charHalfBodyTop, got %q", ch)
	}
}

func TestCharForCell_HalfBodyBottom(t *testing.T) {
	// botSub가 몸통, topSub가 범위 밖 → charHalfBodyBottom
	// highSR=10, lowSR=5, bodyBot=5, bodyTop=7
	// topSub=4 (범위 밖), botSub=5 (bodyBot)
	ch := charForCell(4, 5, 10, 5, 7, 5)
	if ch != charHalfBodyBottom {
		t.Errorf("expected charHalfBodyBottom, got %q", ch)
	}
}

// --- RenderTabBar tests ---

func TestRenderTabBar_Single(t *testing.T) {
	markets := []string{"KRW-BTC"}
	result := RenderTabBar(markets, 0)
	// 단일 마켓: 선택 상태이므로 "[KRW-BTC]" 포함
	if !strings.Contains(result, "KRW-BTC") {
		t.Errorf("expected KRW-BTC in tab bar, got %q", result)
	}
	if !strings.Contains(result, "[") || !strings.Contains(result, "]") {
		t.Errorf("expected brackets for selected market, got %q", result)
	}
}

func TestRenderTabBar_MultipleSelected(t *testing.T) {
	markets := []string{"KRW-BTC", "KRW-ETH", "KRW-XRP"}
	result := RenderTabBar(markets, 1)
	// 선택된 마켓에 [] 표시
	if !strings.Contains(result, "[KRW-ETH]") {
		t.Errorf("expected [KRW-ETH] in tab bar, got %q", result)
	}
	// 비선택 마켓은 [] 없이
	if strings.Contains(result, "[KRW-BTC]") {
		t.Errorf("expected KRW-BTC without brackets, got %q", result)
	}
	if strings.Contains(result, "[KRW-XRP]") {
		t.Errorf("expected KRW-XRP without brackets, got %q", result)
	}
}

func TestRenderTabBar_FirstSelected(t *testing.T) {
	markets := []string{"KRW-BTC", "KRW-ETH"}
	result := RenderTabBar(markets, 0)
	if !strings.Contains(result, "[KRW-BTC]") {
		t.Errorf("expected [KRW-BTC] selected, got %q", result)
	}
	if strings.Contains(result, "[KRW-ETH]") {
		t.Errorf("expected KRW-ETH without brackets, got %q", result)
	}
}

// --- smartPrice tests ---

func TestSmartPrice_LargeNumber(t *testing.T) {
	// >= 100 → 정수
	result := smartPrice(106000000)
	if result != "106000000" {
		t.Errorf("expected 106000000, got %q", result)
	}
}

func TestSmartPrice_Medium(t *testing.T) {
	// >= 1 → 소수 2자리
	result := smartPrice(12.345)
	if result != "12.35" {
		t.Errorf("expected 12.35, got %q", result)
	}
}

func TestSmartPrice_Small(t *testing.T) {
	// >= 0.01 → 소수 4자리
	result := smartPrice(0.0512)
	if result != "0.0512" {
		t.Errorf("expected 0.0512, got %q", result)
	}
}

func TestSmartPrice_VerySmall(t *testing.T) {
	// < 0.01 → 소수 8자리
	result := smartPrice(0.00002045)
	if result != "0.00002045" {
		t.Errorf("expected 0.00002045, got %q", result)
	}
}

func TestSmartPrice_Zero(t *testing.T) {
	result := smartPrice(0)
	if result != "0" {
		t.Errorf("expected 0, got %q", result)
	}
}

func TestSmartPrice_Exactly100(t *testing.T) {
	// 정확히 100 → 정수
	result := smartPrice(100)
	if result != "100" {
		t.Errorf("expected 100, got %q", result)
	}
}

// --- TruncateToHeight tests ---

func TestTruncateToHeight_WithinLimit(t *testing.T) {
	input := "line1\nline2\nline3"
	result := TruncateToHeight(input, 5)
	if result != input {
		t.Errorf("expected unchanged, got %q", result)
	}
}

func TestTruncateToHeight_ExceedsLimit(t *testing.T) {
	input := "line1\nline2\nline3\nline4\nline5"
	result := TruncateToHeight(input, 3)
	lines := strings.Split(result, "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d: %q", len(lines), result)
	}
	if lines[0] != "line1" || lines[1] != "line2" || lines[2] != "line3" {
		t.Errorf("unexpected lines: %q", result)
	}
}

func TestTruncateToHeight_ZeroHeight(t *testing.T) {
	input := "line1\nline2"
	result := TruncateToHeight(input, 0)
	// height <= 0 → 원본 그대로
	if result != input {
		t.Errorf("expected unchanged for zero height, got %q", result)
	}
}

// --- CandleDataFromOHLCV tests ---

func TestCandleDataFromOHLCV(t *testing.T) {
	cd := CandleDataFromOHLCV(100.0, 120.0, 90.0, 110.0, 5.5, "2024-01-01T12:00:00")
	if cd.Open != 100.0 {
		t.Errorf("Open: expected 100.0, got %v", cd.Open)
	}
	if cd.High != 120.0 {
		t.Errorf("High: expected 120.0, got %v", cd.High)
	}
	if cd.Low != 90.0 {
		t.Errorf("Low: expected 90.0, got %v", cd.Low)
	}
	if cd.Close != 110.0 {
		t.Errorf("Close: expected 110.0, got %v", cd.Close)
	}
	if cd.Volume != 5.5 {
		t.Errorf("Volume: expected 5.5, got %v", cd.Volume)
	}
	if cd.Time != "2024-01-01T12:00:00" {
		t.Errorf("Time: expected 2024-01-01T12:00:00, got %q", cd.Time)
	}
}
