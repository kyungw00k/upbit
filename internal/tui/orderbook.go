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
	dataMap map[string]*ws.OrderbookStream // 마켓별 데이터
	markets []string                       // 마켓 순서 (수신 순)
	current int                            // 현재 선택 인덱스
	width   int
	height  int
}

// NewOrderbookModel OrderbookModel 생성
func NewOrderbookModel() OrderbookModel {
	return OrderbookModel{
		dataMap: make(map[string]*ws.OrderbookStream),
	}
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
		case "tab", "right":
			if len(m.markets) > 0 {
				m.current = (m.current + 1) % len(m.markets)
			}
		case "shift+tab", "left":
			if len(m.markets) > 0 {
				m.current = (m.current - 1 + len(m.markets)) % len(m.markets)
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case WsMsg:
		var o ws.OrderbookStream
		if err := json.Unmarshal(msg.Data, &o); err == nil && o.Code != "" {
			if _, exists := m.dataMap[o.Code]; !exists {
				m.markets = append(m.markets, o.Code)
			}
			m.dataMap[o.Code] = &o
		}
	case ErrMsg:
		return m, tea.Quit
	case ServerErrMsg:
		return m, tea.Quit
	}
	return m, nil
}

func (m OrderbookModel) View() string {
	multiMarket := len(m.markets) > 1

	// 현재 마켓 결정
	var currentMarket string
	var data *ws.OrderbookStream
	if len(m.markets) > 0 {
		currentMarket = m.markets[m.current]
		data = m.dataMap[currentMarket]
	}

	if data == nil || len(data.OrderbookUnits) == 0 {
		return "\n  Waiting for data...\n"
	}

	o := data
	units := o.OrderbookUnits
	var b strings.Builder

	// 탭 바 (복수 마켓일 때만)
	if multiMarket {
		b.WriteString(RenderTabBar(m.markets, m.current))
		b.WriteString("\n\n")
	}

	// 타이틀: 마켓명 + 총매도/총매수
	title := fmt.Sprintf("%s %s  %s: %s  %s: %s",
		i18n.T(i18n.TUIOrderbookTitle), StyleTitle.Render(currentMarket),
		i18n.T(i18n.TUITotalAsk), smartVolume(o.TotalAskSize),
		i18n.T(i18n.TUITotalBid), smartVolume(o.TotalBidSize))
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
	// 고정 줄: 탭바(복수 마켓시 2줄) + 헤더(2) + 스프레드(1) + 하단(2) = 5 or 7줄
	fixedLines := 5
	if multiMarket {
		fixedLines = 7
	}
	availableLines := m.height - fixedLines
	if availableLines < 2 {
		availableLines = 2
	}
	showPerSide := availableLines / 2
	if showPerSide < 1 {
		showPerSide = 1
	}
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
		size := smartVolume(u.AskSize)
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
		size := smartVolume(u.BidSize)
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
	if multiMarket {
		b.WriteString(StyleHint.Render(i18n.T(i18n.TUITabHint)))
	} else {
		b.WriteString(StyleHint.Render(i18n.T(i18n.TUIQuitHint)))
	}
	b.WriteString("\n")

	return TruncateToHeight(b.String(), m.height)
}
