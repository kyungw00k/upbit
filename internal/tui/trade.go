package tui

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	ws "github.com/kyungw00k/upbit/api/websocket"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/types"
)

// TradeModel watch trade TUI 모델
type TradeModel struct {
	tradesMap map[string][]ws.TradeStream // 마켓별 체결 데이터
	markets   []string
	current   int
	maxLines  int
	width     int
	height    int
}

// NewTradeModel TradeModel 생성
func NewTradeModel() TradeModel {
	return TradeModel{
		tradesMap: make(map[string][]ws.TradeStream),
		maxLines:  30, // 기본값, WindowSizeMsg로 조정
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
		// 헤더(2줄) + 힌트(2줄)를 빼고 체결 표시 가능 행 수
		m.maxLines = msg.Height - 4
		if m.maxLines < 5 {
			m.maxLines = 5
		}
		// 기존 데이터가 maxLines 초과하면 자름
		for code, trades := range m.tradesMap {
			if len(trades) > m.maxLines {
				m.tradesMap[code] = trades[:m.maxLines]
			}
		}
	case WsMsg:
		var t ws.TradeStream
		if err := json.Unmarshal(msg.Data, &t); err == nil && t.Code != "" {
			if _, exists := m.tradesMap[t.Code]; !exists {
				m.markets = append(m.markets, t.Code)
				m.tradesMap[t.Code] = make([]ws.TradeStream, 0)
			}
			// 새 체결을 맨 앞에 추가
			m.tradesMap[t.Code] = append([]ws.TradeStream{t}, m.tradesMap[t.Code]...)
			if len(m.tradesMap[t.Code]) > m.maxLines {
				m.tradesMap[t.Code] = m.tradesMap[t.Code][:m.maxLines]
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
	multiMarket := len(m.markets) > 1

	// 현재 마켓 결정
	var currentMarket string
	var trades []ws.TradeStream
	if len(m.markets) > 0 {
		currentMarket = m.markets[m.current]
		trades = m.tradesMap[currentMarket]
	}

	if len(trades) == 0 {
		return "\n  Waiting for data...\n"
	}

	var b strings.Builder

	// 탭 바 (복수 마켓일 때만)
	if multiMarket {
		b.WriteString(RenderTabBar(m.markets, m.current))
		b.WriteString("\n\n")
	}

	// 타이틀
	b.WriteString(StyleTitle.Render(i18n.T(i18n.TUITradeTitle)))
	b.WriteString("\n")

	// 각 체결 행
	for _, t := range trades {
		side := i18n.T(i18n.WatchBuy)
		style := StyleRise
		if t.AskBid == "ASK" {
			side = i18n.T(i18n.WatchSell)
			style = StyleFall
		}

		ts := time.UnixMilli(t.TradeTimestamp).In(types.KSTLoc).Format("15:04:05")
		price := smartPrice(t.TradePrice)
		volume := smartVolume(t.TradeVolume)

		line := fmt.Sprintf("  %s  %s  %s  %s  %s",
			padRight(currentMarket, 14),
			style.Render(padRight(side, 4)),
			style.Render(padLeft(price, 16)),
			padLeft(volume, 16),
			StyleHint.Render(ts))
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
