package indicator

import "math"

// CandleData holds OHLCV data for a single candle.
type CandleData struct {
	Open, High, Low, Close, Volume float64
	Time                           string
}

// Result holds all calculated technical indicators for the latest candle.
type Result struct {
	RSI                                         float64
	MACDLine, MACDSignal, MACDHist              float64
	BBUpper, BBMiddle, BBLower                  float64
	EMA20, EMA50, ATR, VolRatio                 float64
	Trend                                       string // "up", "down", "neutral"
}

// EMA computes the Exponential Moving Average.
// The first value is seeded with SMA of the first `period` elements.
// Returns a slice of the same length as input. Positions before the
// seed window are filled with NaN.
func EMA(values []float64, period int) []float64 {
	n := len(values)
	result := make([]float64, n)
	if period <= 0 || n < period {
		for i := range result {
			result[i] = math.NaN()
		}
		return result
	}

	k := 2.0 / float64(period+1)

	// Seed: SMA of first `period` values.
	var sum float64
	for i := 0; i < period; i++ {
		sum += values[i]
	}
	seed := sum / float64(period)

	for i := 0; i < period-1; i++ {
		result[i] = math.NaN()
	}
	result[period-1] = seed

	for i := period; i < n; i++ {
		result[i] = values[i]*k + result[i-1]*(1-k)
	}
	return result
}

// RSI computes the Relative Strength Index using Wilder's smoothing.
// Returns a slice of the same length as input. The first valid value
// appears at index `period`.
func RSI(closes []float64, period int) []float64 {
	n := len(closes)
	result := make([]float64, n)
	if period <= 0 || n <= period {
		for i := range result {
			result[i] = math.NaN()
		}
		return result
	}

	for i := 0; i < period; i++ {
		result[i] = math.NaN()
	}

	// Initial average gain/loss from first `period` deltas.
	var avgGain, avgLoss float64
	for i := 1; i <= period; i++ {
		delta := closes[i] - closes[i-1]
		if delta > 0 {
			avgGain += delta
		} else {
			avgLoss += math.Abs(delta)
		}
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	if avgLoss == 0 {
		result[period] = 100
	} else {
		rs := avgGain / avgLoss
		result[period] = 100 - 100/(1+rs)
	}

	// Wilder's smoothing for subsequent values.
	for i := period + 1; i < n; i++ {
		delta := closes[i] - closes[i-1]
		var gain, loss float64
		if delta > 0 {
			gain = delta
		} else {
			loss = math.Abs(delta)
		}
		avgGain = (avgGain*float64(period-1) + gain) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + loss) / float64(period)

		if avgLoss == 0 {
			result[i] = 100
		} else {
			rs := avgGain / avgLoss
			result[i] = 100 - 100/(1+rs)
		}
	}

	return result
}

// MACD computes the Moving Average Convergence Divergence.
// Returns three slices (macdLine, signalLine, histogram) of the same length as input.
func MACD(closes []float64, fast, slow, signal int) (macdLine, signalLine, histogram []float64) {
	n := len(closes)
	macdLine = make([]float64, n)
	signalLine = make([]float64, n)
	histogram = make([]float64, n)

	if n == 0 {
		return
	}

	emaFast := EMA(closes, fast)
	emaSlow := EMA(closes, slow)

	// Build MACD line.
	for i := 0; i < n; i++ {
		if math.IsNaN(emaFast[i]) || math.IsNaN(emaSlow[i]) {
			macdLine[i] = math.NaN()
		} else {
			macdLine[i] = emaFast[i] - emaSlow[i]
		}
	}

	// Extract valid MACD values for signal EMA calculation.
	validStart := -1
	for i := 0; i < n; i++ {
		if !math.IsNaN(macdLine[i]) {
			validStart = i
			break
		}
	}

	if validStart < 0 {
		return
	}

	validMACD := make([]float64, 0, n-validStart)
	for i := validStart; i < n; i++ {
		validMACD = append(validMACD, macdLine[i])
	}

	emaSignal := EMA(validMACD, signal)

	for i := 0; i < n; i++ {
		signalLine[i] = math.NaN()
		histogram[i] = math.NaN()
	}

	for i, v := range emaSignal {
		idx := validStart + i
		if idx < n {
			signalLine[idx] = v
			if !math.IsNaN(macdLine[idx]) && !math.IsNaN(v) {
				histogram[idx] = macdLine[idx] - v
			}
		}
	}

	return
}

// BollingerBands computes the Bollinger Bands.
// Returns three slices (upper, middle, lower) of the same length as input.
func BollingerBands(closes []float64, period int, stdDev float64) (upper, middle, lower []float64) {
	n := len(closes)
	upper = make([]float64, n)
	middle = make([]float64, n)
	lower = make([]float64, n)

	if period <= 0 || n == 0 {
		for i := 0; i < n; i++ {
			upper[i] = math.NaN()
			middle[i] = math.NaN()
			lower[i] = math.NaN()
		}
		return
	}

	sma := SMA(closes, period)

	for i := 0; i < n; i++ {
		if math.IsNaN(sma[i]) {
			upper[i] = math.NaN()
			middle[i] = math.NaN()
			lower[i] = math.NaN()
			continue
		}

		// Population standard deviation of the window.
		start := i - period + 1
		windowMean := sma[i]
		var variance float64
		for j := start; j <= i; j++ {
			diff := closes[j] - windowMean
			variance += diff * diff
		}
		variance /= float64(period)
		sd := math.Sqrt(variance)

		middle[i] = sma[i]
		upper[i] = sma[i] + stdDev*sd
		lower[i] = sma[i] - stdDev*sd
	}

	return
}

// ATR computes the Average True Range using Wilder's smoothing.
// Returns a slice of the same length as candles. The first valid value
// appears at index `period`.
func ATR(candles []CandleData, period int) []float64 {
	n := len(candles)
	result := make([]float64, n)
	if period <= 0 || n <= 1 {
		for i := range result {
			result[i] = math.NaN()
		}
		return result
	}

	// Compute True Range series.
	tr := make([]float64, n)
	tr[0] = candles[0].High - candles[0].Low
	for i := 1; i < n; i++ {
		tr[i] = TrueRange(candles[i].High, candles[i].Low, candles[i-1].Close)
	}

	for i := 0; i < period; i++ {
		result[i] = math.NaN()
	}

	if n <= period {
		return result
	}

	// Initial ATR: simple average of first `period` TR values.
	var sum float64
	for i := 0; i < period; i++ {
		sum += tr[i]
	}
	result[period] = sum / float64(period)

	// Wilder's smoothing.
	for i := period + 1; i < n; i++ {
		result[i] = (result[i-1]*float64(period-1) + tr[i]) / float64(period)
	}

	return result
}

// lastValid returns the last non-NaN value from a slice, or NaN if none exists.
func lastValid(values []float64) float64 {
	for i := len(values) - 1; i >= 0; i-- {
		if !math.IsNaN(values[i]) {
			return values[i]
		}
	}
	return math.NaN()
}

// CalculateAll computes all technical indicators from candle data and returns
// a single Result for the latest candle.
//
// Default parameters:
//   - RSI: period=14
//   - MACD: fast=12, slow=26, signal=9
//   - Bollinger Bands: period=20, stdDev=2.0
//   - ATR: period=14
//   - Trend: based on EMA20 and EMA50 alignment with close
//   - VolRatio: last volume / 20-period average volume
func CalculateAll(candles []CandleData) Result {
	n := len(candles)
	if n == 0 {
		return Result{Trend: "neutral"}
	}

	closes := make([]float64, n)
	volumes := make([]float64, n)
	for i, c := range candles {
		closes[i] = c.Close
		volumes[i] = c.Volume
	}

	lastClose := closes[n-1]

	// RSI 14
	rsiValues := RSI(closes, 14)
	rsi := lastValid(rsiValues)

	// MACD 12, 26, 9
	macdLine, macdSignal, macdHist := MACD(closes, 12, 26, 9)
	ml := lastValid(macdLine)
	ms := lastValid(macdSignal)
	mh := lastValid(macdHist)

	// Bollinger Bands 20, 2.0
	bbUpper, bbMiddle, bbLower := BollingerBands(closes, 20, 2.0)
	bu := lastValid(bbUpper)
	bm := lastValid(bbMiddle)
	bl := lastValid(bbLower)

	// EMA 20 and EMA 50
	ema20Values := EMA(closes, 20)
	ema50Values := EMA(closes, 50)
	e20 := lastValid(ema20Values)
	e50 := lastValid(ema50Values)

	// ATR 14
	atrValues := ATR(candles, 14)
	atr := lastValid(atrValues)

	// VolRatio: last candle volume / average of last 20 volumes.
	var volRatio float64
	if n >= 20 {
		var volSum float64
		for i := n - 20; i < n; i++ {
			volSum += volumes[i]
		}
		avgVol := volSum / 20
		if avgVol > 0 {
			volRatio = volumes[n-1] / avgVol
		}
	}

	// Trend determination.
	trend := "neutral"
	if !math.IsNaN(e20) && !math.IsNaN(e50) {
		if lastClose > e20 && e20 > e50 {
			trend = "up"
		} else if lastClose < e20 && e20 < e50 {
			trend = "down"
		}
	}

	return Result{
		RSI:        rsi,
		MACDLine:   ml,
		MACDSignal: ms,
		MACDHist:   mh,
		BBUpper:    bu,
		BBMiddle:   bm,
		BBLower:    bl,
		EMA20:      e20,
		EMA50:      e50,
		ATR:        atr,
		VolRatio:   volRatio,
		Trend:      trend,
	}
}
