package tui

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	ws "github.com/kyungw00k/upbit/internal/api/websocket"
	"github.com/kyungw00k/upbit/internal/i18n"
)

// Unicode characters for candlestick rendering (cli-candlestick-chart style)
const (
	charVoid           = ' '  // empty space
	charBody           = '┃'  // full body
	charHalfBodyBottom = '╻'  // bottom half body
	charHalfBodyTop    = '╹'  // top half body
	charWick           = '│'  // thin wick
	charTop            = '╽'  // wick-to-body top
	charBottom         = '╿'  // wick-to-body bottom
	charUpperWick      = '╷'  // short upper wick
	charLowerWick      = '╵'  // short lower wick
)

const (
	candleWidth  = 3  // 1 space + 1 body + 1 space
	yAxisWidth   = 12 // right-side price label width
	headerLines  = 3  // title + OHLC info + blank
	footerLines  = 4  // time axis + summary line + quit hint + blank
	volumeHeight = 3  // height of the volume bar pane
)

// CandleData holds OHLC data for a single candle
type CandleData struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
	Time   string // CandleDateTimeKST
}

// CandleModel watch candle TUI model
type CandleModel struct {
	candlesMap map[string][]CandleData // 마켓별 캔들 데이터
	markets    []string                // 마켓 순서
	current    int                     // 현재 선택 인덱스
	interval   string
	width      int
	height     int
	maxCandles int
}

// NewCandleModel creates a new CandleModel
func NewCandleModel() CandleModel {
	return CandleModel{
		candlesMap: make(map[string][]CandleData),
		maxCandles: 40,
	}
}

// NewCandleModelWithData creates a CandleModel pre-loaded with historical candles
func NewCandleModelWithData(market, interval string, initial []CandleData) CandleModel {
	cm := make(map[string][]CandleData)
	cm[market] = initial
	return CandleModel{
		candlesMap: cm,
		markets:    []string{market},
		interval:   interval,
		maxCandles: 40,
	}
}

// NewCandleModelWithMultiData creates a CandleModel pre-loaded with multiple markets' candles
func NewCandleModelWithMultiData(markets []string, interval string, data map[string][]CandleData) CandleModel {
	return CandleModel{
		candlesMap: data,
		markets:    markets,
		interval:   interval,
		maxCandles: 40,
	}
}

// CandleDataFromOHLCV creates CandleData from raw values (for preloading)
func CandleDataFromOHLCV(open, high, low, close, volume float64, time string) CandleData {
	return CandleData{
		Open: open, High: high, Low: low, Close: close,
		Volume: volume, Time: time,
	}
}

func (m CandleModel) Init() tea.Cmd {
	return nil
}

func (m CandleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		m.maxCandles = (m.width - yAxisWidth) / candleWidth
		if m.maxCandles < 1 {
			m.maxCandles = 1
		}
	case WsMsg:
		var c ws.CandleStream
		if err := json.Unmarshal(msg.Data, &c); err == nil && c.Code != "" {
			code := c.Code
			if _, exists := m.candlesMap[code]; !exists {
				m.markets = append(m.markets, code)
				m.candlesMap[code] = make([]CandleData, 0)
			}
			cd := CandleData{
				Open:   c.OpeningPrice,
				High:   c.HighPrice,
				Low:    c.LowPrice,
				Close:  c.TradePrice,
				Volume: c.CandleAccTradeVolume,
				Time:   c.CandleDateTimeKST,
			}
			candles := m.candlesMap[code]
			// Same time → update last candle; new time → append
			if len(candles) > 0 && candles[len(candles)-1].Time == cd.Time {
				candles[len(candles)-1] = cd
			} else {
				candles = append(candles, cd)
			}
			// Trim buffer
			bufSize := m.maxCandles * 2
			if bufSize < 100 {
				bufSize = 100
			}
			if len(candles) > bufSize {
				candles = candles[len(candles)-bufSize:]
			}
			m.candlesMap[code] = candles
		}
	case ErrMsg:
		return m, tea.Quit
	case ServerErrMsg:
		return m, tea.Quit
	}
	return m, nil
}

func (m CandleModel) View() string {
	multiMarket := len(m.markets) > 1

	// 현재 마켓 결정
	var currentMarket string
	var candles []CandleData
	if len(m.markets) > 0 {
		currentMarket = m.markets[m.current]
		candles = m.candlesMap[currentMarket]
	}

	if len(candles) == 0 {
		return "\n  Waiting for data...\n"
	}

	var b strings.Builder

	// 탭 바 (복수 마켓일 때만)
	if multiMarket {
		b.WriteString(RenderTabBar(m.markets, m.current))
		b.WriteString("\n\n")
	}

	// Determine visible candles
	extraLines := 0
	if multiMarket {
		extraLines = 2 // 탭 바 + 빈 줄
	}
	chartHeight := m.height - headerLines - footerLines - volumeHeight - extraLines
	if chartHeight < 5 {
		chartHeight = 5
	}
	maxCandles := m.maxCandles
	if maxCandles < 1 {
		maxCandles = (m.width - yAxisWidth) / candleWidth
		if maxCandles < 1 {
			maxCandles = 1
		}
	}

	visible := candles
	if len(visible) > maxCandles {
		visible = visible[len(visible)-maxCandles:]
	}

	last := visible[len(visible)-1]

	// Header: market + latest OHLC
	direction := "EVEN"
	if last.Close > last.Open {
		direction = "RISE"
	} else if last.Close < last.Open {
		direction = "FALL"
	}
	priceStyle := PriceStyle(direction)

	intervalLabel := m.interval
	if intervalLabel == "" {
		intervalLabel = "?"
	}
	title := fmt.Sprintf("%s %s [%s]",
		i18n.T(i18n.TUICandleTitle),
		StyleTitle.Render(currentMarket),
		intervalLabel)
	b.WriteString(StyleTitle.Render(title))
	b.WriteString("\n")

	ohlcLine := fmt.Sprintf("  %s:%s  %s:%s  %s:%s  %s:%s  %s:%s",
		i18n.T(i18n.WatchOpen), priceStyle.Render(smartPrice(last.Open)),
		i18n.T(i18n.WatchHigh), priceStyle.Render(smartPrice(last.High)),
		i18n.T(i18n.WatchLow), priceStyle.Render(smartPrice(last.Low)),
		i18n.T(i18n.WatchClose), priceStyle.Render(smartPrice(last.Close)),
		i18n.T(i18n.WatchVolume), smartVolume(last.Volume))
	b.WriteString(ohlcLine)
	b.WriteString("\n\n")

	// Render candle chart
	chart := renderCandleChart(visible, chartHeight)
	b.WriteString(chart)

	// Render volume bar pane
	b.WriteString(renderVolumePane(visible, volumeHeight))

	// Footer: time labels
	b.WriteString(renderTimeAxis(visible, maxCandles))
	b.WriteString("\n")

	// Summary line
	b.WriteString(renderSummaryLine(visible))
	b.WriteString("\n")

	// Hint
	if multiMarket {
		b.WriteString(StyleHint.Render(i18n.T(i18n.TUITabHint)))
	} else {
		b.WriteString(StyleHint.Render(i18n.T(i18n.TUIQuitHint)))
	}
	b.WriteString("\n")

	return TruncateToHeight(b.String(), m.height)
}

// renderCandleChart renders the candlestick chart area
func renderCandleChart(candles []CandleData, chartHeight int) string {
	if len(candles) == 0 || chartHeight < 1 {
		return ""
	}

	// Calculate global min/max prices
	minPrice := math.MaxFloat64
	maxPrice := -math.MaxFloat64
	for _, c := range candles {
		if c.Low < minPrice {
			minPrice = c.Low
		}
		if c.High > maxPrice {
			maxPrice = c.High
		}
	}

	// Add small padding to avoid edge artifacts
	priceRange := maxPrice - minPrice
	if priceRange == 0 {
		priceRange = 1
		minPrice -= 0.5
		maxPrice += 0.5
	}

	// Use half-character resolution: each row represents 2 sub-rows
	// This doubles the effective vertical resolution
	subRows := chartHeight * 2

	// Map price to sub-row (0 = bottom, subRows-1 = top)
	priceToSubRow := func(price float64) int {
		ratio := (price - minPrice) / priceRange
		sr := int(math.Round(ratio * float64(subRows-1)))
		if sr < 0 {
			sr = 0
		}
		if sr >= subRows {
			sr = subRows - 1
		}
		return sr
	}

	// Pre-compute each candle's sub-row positions
	type candlePos struct {
		highSR   int
		lowSR    int
		bodyTop  int
		bodyBot  int
		isRise   bool
		isFlat   bool
	}
	positions := make([]candlePos, len(candles))
	for i, c := range candles {
		pos := candlePos{
			highSR: priceToSubRow(c.High),
			lowSR:  priceToSubRow(c.Low),
		}
		if c.Close >= c.Open {
			pos.bodyTop = priceToSubRow(c.Close)
			pos.bodyBot = priceToSubRow(c.Open)
			pos.isRise = true
		} else {
			pos.bodyTop = priceToSubRow(c.Open)
			pos.bodyBot = priceToSubRow(c.Close)
			pos.isRise = false
		}
		if c.Close == c.Open {
			pos.isFlat = true
		}
		positions[i] = pos
	}

	// Y-axis label positions (5-8 labels distributed evenly)
	numLabels := 6
	if chartHeight < 10 {
		numLabels = 4
	}
	labelRows := make(map[int]string, numLabels)
	for i := 0; i < numLabels; i++ {
		row := i * (chartHeight - 1) / (numLabels - 1)
		price := maxPrice - float64(row)/float64(chartHeight-1)*priceRange
		labelRows[row] = smartPrice(price)
	}

	// Render row by row (top to bottom)
	var b strings.Builder
	for row := 0; row < chartHeight; row++ {
		// Each display row corresponds to sub-rows: topSub and botSub
		topSub := subRows - 1 - row*2
		botSub := topSub - 1
		if botSub < 0 {
			botSub = 0
		}

		// Render each candle column
		for _, pos := range positions {
			ch := charForCell(topSub, botSub, pos.highSR, pos.lowSR, pos.bodyTop, pos.bodyBot)

			var style lipgloss.Style
			if pos.isFlat {
				style = StyleFlat
			} else if pos.isRise {
				style = StyleRise
			} else {
				style = StyleFall
			}

			if ch == charVoid {
				b.WriteString("   ")
			} else {
				b.WriteByte(' ')
				b.WriteString(style.Render(string(ch)))
				b.WriteByte(' ')
			}
		}

		// Y-axis label
		if label, ok := labelRows[row]; ok {
			b.WriteString(" " + padLeft(label, yAxisWidth-1))
		}

		b.WriteByte('\n')
	}

	return b.String()
}

// charForCell determines the character for a cell given sub-row positions.
// topSub and botSub are the two sub-rows this display row covers.
func charForCell(topSub, botSub, highSR, lowSR, bodyTop, bodyBot int) rune {
	// Check if either sub-row is in the candle's range
	topInRange := topSub >= lowSR && topSub <= highSR
	botInRange := botSub >= lowSR && botSub <= highSR

	if !topInRange && !botInRange {
		return charVoid
	}

	topInBody := topSub >= bodyBot && topSub <= bodyTop
	botInBody := botSub >= bodyBot && botSub <= bodyTop
	topInWick := topInRange && !topInBody
	botInWick := botInRange && !botInBody

	// Both sub-rows in body
	if topInBody && botInBody {
		return charBody
	}

	// Only top in body
	if topInBody && !botInBody {
		if botInWick {
			return charBottom // body-to-wick transition at bottom
		}
		return charHalfBodyTop
	}

	// Only bottom in body
	if !topInBody && botInBody {
		if topInWick {
			return charTop // wick-to-body transition at top
		}
		return charHalfBodyBottom
	}

	// Both sub-rows are wick
	if topInWick && botInWick {
		return charWick
	}

	// Only top is wick
	if topInWick && !botInWick {
		return charLowerWick // top sub-row only → appears as lower part of cell
	}

	// Only bottom is wick
	if !topInWick && botInWick {
		return charUpperWick
	}

	// Doji / single line case: body is zero-height
	if bodyTop == bodyBot {
		if topSub == bodyTop || botSub == bodyTop {
			if topInRange && botInRange {
				return charWick
			}
			if topInRange {
				return charLowerWick
			}
			return charUpperWick
		}
	}

	return charVoid
}

// renderTimeAxis renders the time labels below the chart
func renderTimeAxis(candles []CandleData, maxCandles int) string {
	if len(candles) == 0 {
		return ""
	}

	// 시간 라벨은 "HH:MM" (5글자), 캔들 폭은 3글자
	// 라벨이 겹치지 않도록 최소 2캔들(6글자) 간격으로 표시
	const labelWidth = 5
	minGap := (labelWidth + candleWidth - 1) / candleWidth // 캔들 단위 최소 간격 = 2
	if minGap < 2 {
		minGap = 2
	}

	labelInterval := len(candles) / 5
	if labelInterval < minGap {
		labelInterval = minGap
	}

	// 라벨 위치 결정
	labelAt := make(map[int]string)
	for i, c := range candles {
		if i%labelInterval == 0 {
			labelAt[i] = extractTime(c.Time)
		}
	}
	// 마지막 캔들도 항상 표시 (겹치지 않으면)
	lastIdx := len(candles) - 1
	if _, exists := labelAt[lastIdx]; !exists {
		// 이전 라벨과 충분한 거리?
		tooClose := false
		for idx := range labelAt {
			if lastIdx-idx < minGap {
				tooClose = true
				break
			}
		}
		if !tooClose {
			labelAt[lastIdx] = extractTime(candles[lastIdx].Time)
		}
	}

	// 렌더링: 각 캔들 위치에 라벨 또는 공백
	var b strings.Builder
	skip := 0
	for i := range candles {
		if skip > 0 {
			skip--
			continue
		}
		if label, ok := labelAt[i]; ok {
			b.WriteString(StyleHint.Render(label))
			// 라벨이 차지하는 추가 캔들 수 = (labelWidth - candleWidth) / candleWidth
			extraCandles := (labelWidth - candleWidth) / candleWidth
			if extraCandles < 0 {
				extraCandles = 0
			}
			// 라벨 후 다음 캔들까지 패딩
			remaining := candleWidth*(extraCandles+1) - labelWidth
			for j := 0; j < remaining; j++ {
				b.WriteByte(' ')
			}
			skip = extraCandles
		} else {
			for j := 0; j < candleWidth; j++ {
				b.WriteByte(' ')
			}
		}
	}

	return b.String()
}

// extractTime extracts HH:MM from a KST datetime string
func extractTime(datetime string) string {
	// Format: "2024-01-01T15:30:00" → "15:30"
	if len(datetime) >= 16 {
		return datetime[11:16]
	}
	if len(datetime) >= 5 {
		return datetime[len(datetime)-5:]
	}
	return datetime
}

// renderVolumePane renders a volume bar chart pane below the candle chart.
// height specifies the number of display rows for the pane.
func renderVolumePane(candles []CandleData, height int) string {
	if len(candles) == 0 || height < 1 {
		return ""
	}

	// Find max volume
	maxVol := 0.0
	for _, c := range candles {
		if c.Volume > maxVol {
			maxVol = c.Volume
		}
	}
	if maxVol == 0 {
		maxVol = 1
	}

	// Use half-character resolution: each row represents 2 sub-rows
	subRows := height * 2

	var b strings.Builder
	for row := 0; row < height; row++ {
		// topSub: sub-row index from top (0 = top of pane, subRows-1 = bottom)
		// We render top-to-bottom, so row 0 is the topmost row.
		// A bar of ratio r fills sub-rows from bottom up to int(r * subRows).
		for _, c := range candles {
			ratio := c.Volume / maxVol
			filledSubRows := int(math.Round(ratio * float64(subRows)))
			if filledSubRows > subRows {
				filledSubRows = subRows
			}

			// The current display row covers sub-rows (from bottom):
			// botSubFromBottom = (height - 1 - row) * 2
			// topSubFromBottom = botSubFromBottom + 1
			botSubFromBottom := (height - 1 - row) * 2
			topSubFromBottom := botSubFromBottom + 1

			topFilled := topSubFromBottom < filledSubRows
			botFilled := botSubFromBottom < filledSubRows

			var style lipgloss.Style
			if c.Close >= c.Open {
				style = StyleRise
			} else {
				style = StyleFall
			}

			var ch string
			switch {
			case topFilled && botFilled:
				ch = style.Render("█")
			case !topFilled && botFilled:
				ch = style.Render("▄")
			default:
				ch = " "
			}

			b.WriteByte(' ')
			b.WriteString(ch)
			b.WriteByte(' ')
		}
		b.WriteByte('\n')
	}

	return b.String()
}

// renderSummaryLine renders a one-line summary of visible candles.
func renderSummaryLine(candles []CandleData) string {
	if len(candles) == 0 {
		return ""
	}

	last := candles[len(candles)-1]
	first := candles[0]

	// Highest / Lowest across all visible candles
	highestPrice := -math.MaxFloat64
	lowestPrice := math.MaxFloat64
	totalVol := 0.0
	sumClose := 0.0
	for _, c := range candles {
		if c.High > highestPrice {
			highestPrice = c.High
		}
		if c.Low < lowestPrice {
			lowestPrice = c.Low
		}
		totalVol += c.Volume
		sumClose += c.Close
	}
	avgClose := sumClose / float64(len(candles))

	// Var: (last close - first open) / first open * 100
	varPct := 0.0
	if first.Open != 0 {
		varPct = (last.Close - first.Open) / first.Open * 100
	}
	varSign := "+"
	varStyle := StyleRise
	if varPct < 0 {
		varSign = ""
		varStyle = StyleFall
	} else if varPct == 0 {
		varStyle = StyleFlat
	}

	summary := fmt.Sprintf(
		"%s: %s  %s: %s  %s: %s  %s: %s  %s: %s  %s: %s",
		i18n.T(i18n.WatchClose), smartPrice(last.Close),
		i18n.T(i18n.TUIHighest), smartPrice(highestPrice),
		i18n.T(i18n.TUILowest), smartPrice(lowestPrice),
		i18n.T(i18n.TUIVar), varStyle.Render(fmt.Sprintf("%s%.2f%%", varSign, varPct)),
		i18n.T(i18n.TUIAvg), smartPrice(avgClose),
		i18n.T(i18n.TUICumVol), smartVolume(totalVol),
	)

	return StyleHint.Render(summary)
}
