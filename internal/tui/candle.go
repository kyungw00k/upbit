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
	candles    []CandleData
	market     string
	width      int
	height     int
	maxCandles int
}

// NewCandleModel creates a new CandleModel
func NewCandleModel() CandleModel {
	return CandleModel{
		candles:    make([]CandleData, 0),
		maxCandles: 40,
	}
}

// NewCandleModelWithData creates a CandleModel pre-loaded with historical candles
func NewCandleModelWithData(market string, initial []CandleData) CandleModel {
	return CandleModel{
		candles:    initial,
		market:     market,
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
			m.market = c.Code
			cd := CandleData{
				Open:   c.OpeningPrice,
				High:   c.HighPrice,
				Low:    c.LowPrice,
				Close:  c.TradePrice,
				Volume: c.CandleAccTradeVolume,
				Time:   c.CandleDateTimeKST,
			}
			// Same time → update last candle; new time → append
			if len(m.candles) > 0 && m.candles[len(m.candles)-1].Time == cd.Time {
				m.candles[len(m.candles)-1] = cd
			} else {
				m.candles = append(m.candles, cd)
			}
			// Trim old candles beyond buffer (keep 2x maxCandles for safety)
			bufSize := m.maxCandles * 2
			if bufSize < 100 {
				bufSize = 100
			}
			if len(m.candles) > bufSize {
				m.candles = m.candles[len(m.candles)-bufSize:]
			}
		}
	case ErrMsg:
		return m, tea.Quit
	case ServerErrMsg:
		return m, tea.Quit
	}
	return m, nil
}

func (m CandleModel) View() string {
	if len(m.candles) == 0 {
		return "\n  Waiting for data...\n"
	}

	var b strings.Builder

	// Determine visible candles
	chartHeight := m.height - headerLines - footerLines - volumeHeight
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

	visible := m.candles
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

	title := fmt.Sprintf("%s %s",
		i18n.T(i18n.TUICandleTitle),
		StyleTitle.Render(m.market))
	b.WriteString(StyleTitle.Render(title))
	b.WriteString("\n")

	ohlcLine := fmt.Sprintf("  %s:%s  %s:%s  %s:%s  %s:%s  %s:%.4f",
		i18n.T(i18n.WatchOpen), priceStyle.Render(smartPrice(last.Open)),
		i18n.T(i18n.WatchHigh), priceStyle.Render(smartPrice(last.High)),
		i18n.T(i18n.WatchLow), priceStyle.Render(smartPrice(last.Low)),
		i18n.T(i18n.WatchClose), priceStyle.Render(smartPrice(last.Close)),
		i18n.T(i18n.WatchVolume), last.Volume)
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

	// Quit hint
	b.WriteString(StyleHint.Render(i18n.T(i18n.TUIQuitHint)))
	b.WriteString("\n")

	return b.String()
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

	var b strings.Builder

	// Show time labels at intervals
	labelInterval := len(candles) / 5
	if labelInterval < 1 {
		labelInterval = 1
	}

	// Build time axis with spacing
	for i, c := range candles {
		if i%labelInterval == 0 || i == len(candles)-1 {
			// Extract HH:MM from datetime string (format: YYYY-MM-DDTHH:MM:SS)
			timeLabel := extractTime(c.Time)
			// Pad/truncate to candleWidth
			if len(timeLabel) >= candleWidth {
				b.WriteString(StyleHint.Render(timeLabel[:candleWidth]))
			} else {
				b.WriteString(StyleHint.Render(timeLabel))
				for j := len(timeLabel); j < candleWidth; j++ {
					b.WriteByte(' ')
				}
			}
		} else {
			b.WriteString("   ")
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
		"%s: %s  %s: %s  %s: %s  %s: %s  %s: %s  %s: %.2f",
		i18n.T(i18n.WatchClose), smartPrice(last.Close),
		i18n.T(i18n.TUIHighest), smartPrice(highestPrice),
		i18n.T(i18n.TUILowest), smartPrice(lowestPrice),
		i18n.T(i18n.TUIVar), varStyle.Render(fmt.Sprintf("%s%.2f%%", varSign, varPct)),
		i18n.T(i18n.TUIAvg), smartPrice(avgClose),
		i18n.T(i18n.TUICumVol), totalVol,
	)

	return StyleHint.Render(summary)
}
