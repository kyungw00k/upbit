package tui

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	ws "github.com/kyungw00k/upbit/internal/api/websocket"
	"github.com/kyungw00k/upbit/internal/i18n"
)

var kstLoc = time.FixedZone("KST", 9*60*60)

// TradeModel watch trade TUI 모델
type TradeModel struct {
	trades   []ws.TradeStream // 최근 N건 유지 (ring buffer 대신 슬라이스)
	maxLines int
	width    int
	height   int
}

// NewTradeModel TradeModel 생성
func NewTradeModel() TradeModel {
	return TradeModel{
		trades:   make([]ws.TradeStream, 0),
		maxLines: 30, // 기본값, WindowSizeMsg로 조정
	}
}

func (m TradeModel) Init() tea.Cmd {
	return nil
}

func (m TradeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// 헤더(2줄) + 힌트(2줄)를 빼고 체결 표시 가능 행 수
		m.maxLines = msg.Height - 4
		if m.maxLines < 5 {
			m.maxLines = 5
		}
		// 기존 데이터가 maxLines 초과하면 자름
		if len(m.trades) > m.maxLines {
			m.trades = m.trades[:m.maxLines]
		}
	case WsMsg:
		var t ws.TradeStream
		if err := json.Unmarshal(msg.Data, &t); err == nil && t.Code != "" {
			// 새 체결을 맨 앞에 추가
			m.trades = append([]ws.TradeStream{t}, m.trades...)
			if len(m.trades) > m.maxLines {
				m.trades = m.trades[:m.maxLines]
			}
		}
	case ErrMsg:
		return m, tea.Quit
	case ServerErrMsg:
		return m, tea.Quit
	}
	return m, nil
}

func (m TradeModel) View() string {
	if len(m.trades) == 0 {
		return "\n  Waiting for data...\n"
	}

	var b strings.Builder

	// 타이틀
	b.WriteString(StyleTitle.Render(i18n.T(i18n.TUITradeTitle)))
	b.WriteString("\n")

	// 각 체결 행
	for _, t := range m.trades {
		side := i18n.T(i18n.WatchBuy)
		style := StyleRise
		if t.AskBid == "ASK" {
			side = i18n.T(i18n.WatchSell)
			style = StyleFall
		}

		ts := time.UnixMilli(t.TradeTimestamp).In(kstLoc).Format("15:04:05")
		price := smartPrice(t.TradePrice)
		volume := fmt.Sprintf("%.8f", t.TradeVolume)

		line := fmt.Sprintf("  %s  %s  %s  %s  %s",
			padRight(t.Code, 14),
			style.Render(padRight(side, 4)),
			style.Render(padLeft(price, 16)),
			padLeft(volume, 16),
			StyleHint.Render(ts))
		b.WriteString(line)
		b.WriteString("\n")
	}

	// 하단 힌트
	b.WriteString("\n")
	b.WriteString(StyleHint.Render(i18n.T(i18n.TUIQuitHint)))
	b.WriteString("\n")

	return TruncateToHeight(b.String(), m.height)
}
