package tui

import (
	"encoding/json"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	ws "github.com/kyungw00k/upbit/internal/api/websocket"
	"github.com/kyungw00k/upbit/internal/i18n"
)

// OrderbookModel watch orderbook TUI 모델
type OrderbookModel struct {
	data   *ws.OrderbookStream
	market string
	width  int
	height int
}

// NewOrderbookModel OrderbookModel 생성
func NewOrderbookModel() OrderbookModel {
	return OrderbookModel{}
}

func (m OrderbookModel) Init() tea.Cmd {
	return nil
}

func (m OrderbookModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		var o ws.OrderbookStream
		if err := json.Unmarshal(msg.Data, &o); err == nil && o.Code != "" {
			m.data = &o
			m.market = o.Code
		}
	case ErrMsg:
		return m, tea.Quit
	case ServerErrMsg:
		return m, tea.Quit
	}
	return m, nil
}

func (m OrderbookModel) View() string {
	if m.data == nil || len(m.data.OrderbookUnits) == 0 {
		return "\n  Waiting for data...\n"
	}

	o := m.data
	units := o.OrderbookUnits
	var b strings.Builder

	// 타이틀: 마켓명 + 총매도/총매수
	title := fmt.Sprintf("%s %s  %s: %.4f  %s: %.4f",
		i18n.T(i18n.TUIOrderbookTitle), StyleTitle.Render(m.market),
		i18n.T(i18n.TUITotalAsk), o.TotalAskSize,
		i18n.T(i18n.TUITotalBid), o.TotalBidSize)
	b.WriteString(StyleTitle.Render(title))
	b.WriteString("\n\n")

	// 최대 잔량 계산 (바 비율용)
	maxSize := 0.0
	for _, u := range units {
		if u.AskSize > maxSize {
			maxSize = u.AskSize
		}
		if u.BidSize > maxSize {
			maxSize = u.BidSize
		}
	}
	if maxSize == 0 {
		maxSize = 1
	}

	// 바 차트 최대 너비
	barMaxWidth := 20
	if m.width > 60 {
		barMaxWidth = (m.width - 40) / 2
		if barMaxWidth > 40 {
			barMaxWidth = 40
		}
	}

	// 터미널 높이에 맞게 스프레드 중심으로 위/아래 동일 개수 표시
	// 사용 가능 줄 = 전체 높이 - 헤더(2) - 스프레드(1) - 하단(2)
	availableLines := m.height - 5
	if availableLines < 4 {
		availableLines = 4
	}
	showPerSide := availableLines / 2
	if showPerSide > len(units) {
		showPerSide = len(units)
	}

	// 매도 호가 (스프레드에 가까운 showPerSide개만, 역순)
	askStart := 0
	if showPerSide < len(units) {
		askStart = showPerSide - 1
	} else {
		askStart = len(units) - 1
	}
	for i := askStart; i >= 0; i-- {
		u := units[i]
		price := smartPrice(u.AskPrice)
		size := fmt.Sprintf("%.4f", u.AskSize)
		ratio := u.AskSize / maxSize
		bar := barString(ratio, barMaxWidth)

		line := fmt.Sprintf("  %s  %s  %s",
			StyleFall.Render(padLeft(price, 16)),
			padLeft(size, 12),
			bar)
		b.WriteString(line)
		b.WriteString("\n")
	}

	// 스프레드
	if len(units) > 0 {
		spread := units[0].AskPrice - units[0].BidPrice
		spreadLine := fmt.Sprintf("  %s  %s",
			padLeft("--- spread ---", 16),
			padLeft(smartPrice(spread), 12))
		b.WriteString(StyleHint.Render(spreadLine))
		b.WriteString("\n")
	}

	// 매수 호가 (스프레드에 가까운 showPerSide개)
	for i := 0; i < showPerSide && i < len(units); i++ {
		u := units[i]
		price := smartPrice(u.BidPrice)
		size := fmt.Sprintf("%.4f", u.BidSize)
		ratio := u.BidSize / maxSize
		bar := barString(ratio, barMaxWidth)

		line := fmt.Sprintf("  %s  %s  %s",
			StyleRise.Render(padLeft(price, 16)),
			padLeft(size, 12),
			bar)
		b.WriteString(line)
		b.WriteString("\n")
	}

	// 하단 힌트
	b.WriteString("\n")
	b.WriteString(StyleHint.Render(i18n.T(i18n.TUIQuitHint)))
	b.WriteString("\n")

	return b.String()
}
