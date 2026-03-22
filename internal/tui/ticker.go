package tui

import (
	"encoding/json"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"

	ws "github.com/kyungw00k/upbit/internal/api/websocket"
	"github.com/kyungw00k/upbit/internal/i18n"
)

// TickerModel watch ticker TUI 모델
type TickerModel struct {
	data    map[string]ws.TickerStream // 마켓코드 -> 최신 데이터
	markets []string                   // 표시 순서 (수신 순)
	width   int
	height  int
}

// NewTickerModel TickerModel 생성
func NewTickerModel() TickerModel {
	return TickerModel{
		data:    make(map[string]ws.TickerStream),
		markets: make([]string, 0),
	}
}

func (m TickerModel) Init() tea.Cmd {
	return nil
}

func (m TickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case WsMsg:
		var t ws.TickerStream
		if err := json.Unmarshal(msg.Data, &t); err == nil && t.Code != "" {
			if _, exists := m.data[t.Code]; !exists {
				m.markets = append(m.markets, t.Code)
			}
			m.data[t.Code] = t
		}
	case ErrMsg:
		return m, tea.Quit
	case ServerErrMsg:
		return m, tea.Quit
	}
	return m, nil
}

func (m TickerModel) View() string {
	if len(m.data) == 0 {
		return "\n  Waiting for data...\n"
	}

	var b strings.Builder

	// 헤더
	header := i18n.T(i18n.TUITickerHeader)
	b.WriteString(StyleHeader.Render(header))
	b.WriteString("\n")

	// 각 마켓별 행
	for _, code := range m.markets {
		t, ok := m.data[code]
		if !ok {
			continue
		}
		style := PriceStyle(t.Change)

		arrow := changeArrow(t.Change)
		price := smartPrice(t.TradePrice)
		changeRate := fmt.Sprintf("%s%.2f%%", signPrefix(t.SignedChangeRate), t.SignedChangeRate*100)
		volume := fmt.Sprintf("%.4f", t.AccTradeVolume24h)

		// runewidth 기반 정렬
		codePad := padRight(code, 14)
		pricePad := padLeft(price, 16)
		changePad := padLeft(changeRate, 12)

		line := fmt.Sprintf("%s %s %s  %s  %s",
			codePad, arrow, style.Render(pricePad), style.Render(changePad), volume)

		b.WriteString(line)
		b.WriteString("\n")
	}

	// 하단 힌트
	b.WriteString("\n")
	b.WriteString(StyleHint.Render(i18n.T(i18n.TUIQuitHint)))
	b.WriteString("\n")

	return b.String()
}

// --- 유틸 함수 (watch_streams.go에서 이동) ---

func changeArrow(change string) string {
	switch change {
	case "RISE":
		return StyleRise.Render("▲")
	case "FALL":
		return StyleFall.Render("▼")
	default:
		return StyleFlat.Render("-")
	}
}

func signPrefix(rate float64) string {
	if rate > 0 {
		return "+"
	}
	return ""
}

func smartPrice(price float64) string {
	if price == 0 {
		return "0"
	}
	abs := price
	if abs < 0 {
		abs = -abs
	}
	switch {
	case abs >= 100:
		return fmt.Sprintf("%.0f", price)
	case abs >= 1:
		return fmt.Sprintf("%.2f", price)
	case abs >= 0.01:
		return fmt.Sprintf("%.4f", price)
	default:
		return fmt.Sprintf("%.8f", price)
	}
}

// padRight runewidth 기반 오른쪽 패딩
func padRight(s string, width int) string {
	sw := runewidth.StringWidth(s)
	if sw >= width {
		return s
	}
	return s + strings.Repeat(" ", width-sw)
}

// padLeft runewidth 기반 왼쪽 패딩
func padLeft(s string, width int) string {
	sw := runewidth.StringWidth(s)
	if sw >= width {
		return s
	}
	return strings.Repeat(" ", width-sw) + s
}

// barString 바 차트 문자열 생성
func barString(ratio float64, maxWidth int) string {
	if maxWidth <= 0 {
		maxWidth = 20
	}
	n := int(ratio * float64(maxWidth))
	if n < 0 {
		n = 0
	}
	if n > maxWidth {
		n = maxWidth
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Render(strings.Repeat("█", n))
}
