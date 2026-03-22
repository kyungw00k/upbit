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

	// 1. 전체 데이터를 순수 텍스트로 준비
	type rowData struct {
		code       string
		arrow      string
		price      string
		changeRate string
		volume     string
		change     string // 색상 판단용
	}
	rows := make([]rowData, 0, len(m.markets))
	for _, code := range m.markets {
		t, ok := m.data[code]
		if !ok {
			continue
		}
		rows = append(rows, rowData{
			code:       code,
			arrow:      changeArrowPlain(t.Change),
			price:      smartPrice(t.TradePrice),
			changeRate: fmt.Sprintf("%s%.2f%%", signPrefix(t.SignedChangeRate), t.SignedChangeRate*100),
			volume:     fmt.Sprintf("%.4f", t.AccTradeVolume24h),
			change:     t.Change,
		})
	}

	// 2. 각 컬럼 최대 폭 계산
	headers := [5]string{"Market", "  ", "Price", "Change", "Volume(24h)"}
	if i18n.IsKorean() {
		headers = [5]string{"마켓", "  ", "현재가", "변동률", "거래량(24h)"}
	}
	widths := [5]int{
		runewidth.StringWidth(headers[0]),
		2,
		runewidth.StringWidth(headers[2]),
		runewidth.StringWidth(headers[3]),
		runewidth.StringWidth(headers[4]),
	}
	for _, r := range rows {
		if w := runewidth.StringWidth(r.code); w > widths[0] { widths[0] = w }
		if w := runewidth.StringWidth(r.price); w > widths[2] { widths[2] = w }
		if w := runewidth.StringWidth(r.changeRate); w > widths[3] { widths[3] = w }
		if w := runewidth.StringWidth(r.volume); w > widths[4] { widths[4] = w }
	}

	var b strings.Builder

	// 3. 헤더 출력
	headerLine := fmt.Sprintf("%s  %s  %s  %s  %s",
		padRight(headers[0], widths[0]),
		padRight(headers[1], widths[1]),
		padLeft(headers[2], widths[2]),
		padLeft(headers[3], widths[3]),
		padLeft(headers[4], widths[4]),
	)
	b.WriteString(StyleHeader.Render(headerLine))
	b.WriteString("\n")

	// 4. 데이터 행 출력 (패딩 후 색상)
	for _, r := range rows {
		style := PriceStyle(r.change)
		arrow := changeArrow(r.change)

		line := fmt.Sprintf("%s  %s  %s  %s  %s",
			padRight(r.code, widths[0]),
			arrow,
			style.Render(padLeft(r.price, widths[2])),
			style.Render(padLeft(r.changeRate, widths[3])),
			padLeft(r.volume, widths[4]),
		)
		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(StyleHint.Render(i18n.T(i18n.TUIQuitHint)))
	b.WriteString("\n")

	return b.String()
}

// changeArrowPlain 색상 없는 화살표 (폭 계산용)
func changeArrowPlain(change string) string {
	switch change {
	case "RISE":
		return "▲"
	case "FALL":
		return "▼"
	default:
		return "-"
	}
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
