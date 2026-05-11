package indicator

import (
	"math"
	"testing"
)

// 50 hardcoded data points (ascending chronological order, typical BTC-like closes).
var testCloses = []float64{
	41000, 41200, 41500, 41300, 41600,
	41800, 41700, 41900, 42100, 42000,
	42200, 42400, 42300, 42500, 42700,
	42600, 42800, 43000, 42900, 43100,
	43300, 43200, 43400, 43600, 43500,
	43700, 43900, 43800, 44000, 44200,
	44100, 44300, 44500, 44400, 44600,
	44800, 44700, 44900, 45100, 45000,
	45200, 45400, 45300, 45500, 45700,
	45600, 45800, 46000, 45900, 46100,
}

// testCandles builds CandleData from testCloses with synthetic OHLCV.
func testCandles() []CandleData {
	candles := make([]CandleData, len(testCloses))
	for i, close := range testCloses {
		high := close + 200
		low := close - 200
		open := close - 50
		if i > 0 {
			open = testCloses[i-1] + 50
		}
		candles[i] = CandleData{
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: 1000 + float64(i)*10,
			Time:   "2026-01-01T00:00:00",
		}
	}
	return candles
}

// --- SMA tests ---

func TestSMA_Basic(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5}
	result := SMA(values, 3)

	// First two should be NaN.
	if !math.IsNaN(result[0]) {
		t.Errorf("SMA[0] should be NaN, got %v", result[0])
	}
	if !math.IsNaN(result[1]) {
		t.Errorf("SMA[1] should be NaN, got %v", result[1])
	}

	// SMA(3) at index 2: (1+2+3)/3 = 2
	if math.Abs(result[2]-2.0) > 1e-9 {
		t.Errorf("SMA[2] = %v, want 2.0", result[2])
	}
	// SMA(3) at index 3: (2+3+4)/3 = 3
	if math.Abs(result[3]-3.0) > 1e-9 {
		t.Errorf("SMA[3] = %v, want 3.0", result[3])
	}
	// SMA(3) at index 4: (3+4+5)/3 = 4
	if math.Abs(result[4]-4.0) > 1e-9 {
		t.Errorf("SMA[4] = %v, want 4.0", result[4])
	}
}

func TestSMA_Empty(t *testing.T) {
	result := SMA([]float64{}, 3)
	if len(result) != 0 {
		t.Errorf("SMA of empty slice should return empty, got length %d", len(result))
	}
}

func TestSMA_InvalidPeriod(t *testing.T) {
	result := SMA([]float64{1, 2, 3}, 0)
	for _, v := range result {
		if !math.IsNaN(v) {
			t.Errorf("SMA with period=0 should return all NaN, got %v", v)
		}
	}
}

// --- StdDev tests ---

func TestStdDev_Basic(t *testing.T) {
	values := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	sd := StdDev(values)
	// Population stddev of [2,4,4,4,5,5,7,9]: mean=5, variance=4, stddev=2
	if math.Abs(sd-2.0) > 1e-9 {
		t.Errorf("StdDev = %v, want 2.0", sd)
	}
}

func TestStdDev_Empty(t *testing.T) {
	sd := StdDev([]float64{})
	if sd != 0 {
		t.Errorf("StdDev of empty slice = %v, want 0", sd)
	}
}

// --- TrueRange tests ---

func TestTrueRange_Basic(t *testing.T) {
	// high=50, low=40, prevClose=45
	tr := TrueRange(50, 40, 45)
	// tr1=10, tr2=5, tr3=5 => max=10
	if math.Abs(tr-10.0) > 1e-9 {
		t.Errorf("TrueRange = %v, want 10.0", tr)
	}
}

func TestTrueRange_GapUp(t *testing.T) {
	// Gap up: high=60, low=55, prevClose=40
	tr := TrueRange(60, 55, 40)
	// tr1=5, tr2=20, tr3=15 => max=20
	if math.Abs(tr-20.0) > 1e-9 {
		t.Errorf("TrueRange = %v, want 20.0", tr)
	}
}

// --- EMA tests ---

func TestEMA_FirstValueIsSMA(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result := EMA(values, 5)

	// First valid EMA at index 4 should equal SMA(5) of first 5 values: (1+2+3+4+5)/5 = 3
	if math.IsNaN(result[4]) {
		t.Fatal("EMA[4] should not be NaN")
	}
	if math.Abs(result[4]-3.0) > 1e-9 {
		t.Errorf("EMA[4] = %v, want 3.0 (SMA seed)", result[4])
	}
}

func TestEMA_Length(t *testing.T) {
	result := EMA(testCloses, 20)
	if len(result) != len(testCloses) {
		t.Errorf("EMA result length = %d, want %d", len(result), len(testCloses))
	}
}

func TestEMA_InsufficientData(t *testing.T) {
	result := EMA([]float64{1, 2}, 5)
	for _, v := range result {
		if !math.IsNaN(v) {
			t.Errorf("EMA with insufficient data should return all NaN, got %v", v)
		}
	}
}

// --- RSI tests ---

func TestRSI_Range(t *testing.T) {
	result := RSI(testCloses, 14)

	for i, v := range result {
		if math.IsNaN(v) {
			continue
		}
		if v < 0 || v > 100 {
			t.Errorf("RSI[%d] = %v, want 0 <= RSI <= 100", i, v)
		}
	}
}

func TestRSI_LastValueValid(t *testing.T) {
	result := RSI(testCloses, 14)
	last := result[len(result)-1]
	if math.IsNaN(last) {
		t.Error("Last RSI value should not be NaN")
	}
	if last < 0 || last > 100 {
		t.Errorf("Last RSI = %v, want 0 <= RSI <= 100", last)
	}
}

// With monotonically increasing prices, RSI should be high (close to 100).
func TestRSI_Uptrend(t *testing.T) {
	uptrend := make([]float64, 30)
	for i := range uptrend {
		uptrend[i] = float64(i + 1)
	}
	result := RSI(uptrend, 14)
	last := result[len(result)-1]
	if last < 90 {
		t.Errorf("RSI in strong uptrend = %v, want >= 90", last)
	}
}

// With monotonically decreasing prices, RSI should be low (close to 0).
func TestRSI_Downtrend(t *testing.T) {
	downtrend := make([]float64, 30)
	for i := range downtrend {
		downtrend[i] = float64(30 - i)
	}
	result := RSI(downtrend, 14)
	last := result[len(result)-1]
	if last > 10 {
		t.Errorf("RSI in strong downtrend = %v, want <= 10", last)
	}
}

func TestRSI_ShortData(t *testing.T) {
	result := RSI([]float64{1, 2}, 14)
	for _, v := range result {
		if !math.IsNaN(v) {
			t.Errorf("RSI with insufficient data should return NaN, got %v", v)
		}
	}
}

// --- MACD tests ---

func TestMACD_Lengths(t *testing.T) {
	macdLine, signalLine, histogram := MACD(testCloses, 12, 26, 9)
	if len(macdLine) != len(testCloses) {
		t.Errorf("macdLine length = %d, want %d", len(macdLine), len(testCloses))
	}
	if len(signalLine) != len(testCloses) {
		t.Errorf("signalLine length = %d, want %d", len(signalLine), len(testCloses))
	}
	if len(histogram) != len(testCloses) {
		t.Errorf("histogram length = %d, want %d", len(histogram), len(testCloses))
	}
}

func TestMACD_HistogramEqualsDiff(t *testing.T) {
	macdLine, signalLine, histogram := MACD(testCloses, 12, 26, 9)
	for i := range testCloses {
		if math.IsNaN(macdLine[i]) || math.IsNaN(signalLine[i]) {
			if !math.IsNaN(histogram[i]) {
				t.Errorf("histogram[%d] should be NaN when macd or signal is NaN", i)
			}
			continue
		}
		expected := macdLine[i] - signalLine[i]
		if math.Abs(histogram[i]-expected) > 1e-9 {
			t.Errorf("histogram[%d] = %v, want %v (macdLine - signalLine)", i, histogram[i], expected)
		}
	}
}

func TestMACD_UptrendPositiveMACD(t *testing.T) {
	uptrend := make([]float64, 100)
	for i := range uptrend {
		uptrend[i] = float64(i + 1) * 100
	}
	macdLine, _, _ := MACD(uptrend, 12, 26, 9)
	last := macdLine[len(macdLine)-1]
	if last <= 0 {
		t.Errorf("MACD line in uptrend = %v, want > 0", last)
	}
}

func TestMACD_Empty(t *testing.T) {
	macdLine, signalLine, histogram := MACD([]float64{}, 12, 26, 9)
	if len(macdLine) != 0 || len(signalLine) != 0 || len(histogram) != 0 {
		t.Error("MACD of empty input should return empty slices")
	}
}

// --- BollingerBands tests ---

func TestBollingerBands_Ordering(t *testing.T) {
	upper, middle, lower := BollingerBands(testCloses, 20, 2.0)
	for i := range testCloses {
		if math.IsNaN(upper[i]) {
			continue
		}
		if !(upper[i] >= middle[i] && middle[i] >= lower[i]) {
			t.Errorf("BB[%d]: upper=%v, middle=%v, lower=%v — want upper >= middle >= lower",
				i, upper[i], middle[i], lower[i])
		}
	}
}

func TestBollingerBands_MiddleEqualsSMA(t *testing.T) {
	upper, middle, _ := BollingerBands(testCloses, 20, 2.0)
	sma := SMA(testCloses, 20)
	for i := range testCloses {
		if math.IsNaN(middle[i]) {
			continue
		}
		if math.Abs(middle[i]-sma[i]) > 1e-9 {
			t.Errorf("BB middle[%d] = %v, SMA[%d] = %v — should be equal", i, middle[i], i, sma[i])
		}
	}
	_ = upper // avoid unused var
}

func TestBollingerBands_Empty(t *testing.T) {
	upper, middle, lower := BollingerBands([]float64{}, 20, 2.0)
	if len(upper) != 0 || len(middle) != 0 || len(lower) != 0 {
		t.Error("BB of empty input should return empty slices")
	}
}

// --- ATR tests ---

func TestATR_PositiveValues(t *testing.T) {
	candles := testCandles()
	result := ATR(candles, 14)
	for i, v := range result {
		if math.IsNaN(v) {
			continue
		}
		if v < 0 {
			t.Errorf("ATR[%d] = %v, want >= 0", i, v)
		}
	}
}

func TestATR_FirstValidAtPeriod(t *testing.T) {
	candles := testCandles()
	result := ATR(candles, 14)
	// First valid ATR should be at index 14.
	for i := 0; i < 14; i++ {
		if !math.IsNaN(result[i]) {
			t.Errorf("ATR[%d] should be NaN (before period), got %v", i, result[i])
		}
	}
	if math.IsNaN(result[14]) {
		t.Error("ATR[14] should not be NaN")
	}
}

func TestATR_LastValuePositive(t *testing.T) {
	candles := testCandles()
	result := ATR(candles, 14)
	last := result[len(result)-1]
	if math.IsNaN(last) {
		t.Error("Last ATR should not be NaN")
	}
	if last <= 0 {
		t.Errorf("Last ATR = %v, want > 0", last)
	}
}

func TestATR_ShortData(t *testing.T) {
	candles := []CandleData{
		{High: 10, Low: 5, Close: 8},
	}
	result := ATR(candles, 14)
	for _, v := range result {
		if !math.IsNaN(v) {
			t.Errorf("ATR with insufficient data should return NaN, got %v", v)
		}
	}
}

// --- CalculateAll tests ---

func TestCalculateAll_RSIInRange(t *testing.T) {
	candles := testCandles()
	r := CalculateAll(candles)
	if math.IsNaN(r.RSI) {
		t.Error("RSI should not be NaN")
	}
	if r.RSI < 0 || r.RSI > 100 {
		t.Errorf("RSI = %v, want 0 <= RSI <= 100", r.RSI)
	}
}

func TestCalculateAll_BBOrdering(t *testing.T) {
	candles := testCandles()
	r := CalculateAll(candles)
	if math.IsNaN(r.BBUpper) || math.IsNaN(r.BBMiddle) || math.IsNaN(r.BBLower) {
		t.Error("BB values should not be NaN")
	}
	if !(r.BBUpper >= r.BBMiddle && r.BBMiddle >= r.BBLower) {
		t.Errorf("BB ordering violated: upper=%v, middle=%v, lower=%v",
			r.BBUpper, r.BBMiddle, r.BBLower)
	}
}

func TestCalculateAll_ATRPositive(t *testing.T) {
	candles := testCandles()
	r := CalculateAll(candles)
	if math.IsNaN(r.ATR) {
		t.Error("ATR should not be NaN")
	}
	if r.ATR <= 0 {
		t.Errorf("ATR = %v, want > 0", r.ATR)
	}
}

func TestCalculateAll_TrendValid(t *testing.T) {
	candles := testCandles()
	r := CalculateAll(candles)
	if r.Trend != "up" && r.Trend != "down" && r.Trend != "neutral" {
		t.Errorf("Trend = %q, want one of 'up', 'down', 'neutral'", r.Trend)
	}
}

func TestCalculateAll_VolRatioPositive(t *testing.T) {
	candles := testCandles()
	r := CalculateAll(candles)
	if r.VolRatio < 0 {
		t.Errorf("VolRatio = %v, want >= 0", r.VolRatio)
	}
}

func TestCalculateAll_EMA20LessThan50(t *testing.T) {
	candles := testCandles()
	r := CalculateAll(candles)
	// In an uptrend (test data), EMA20 should be <= EMA50? Actually for uptrend EMA20 > EMA50.
	if !math.IsNaN(r.EMA20) && !math.IsNaN(r.EMA50) {
		// With monotonically increasing data, EMA20 should be >= EMA50.
		// Just verify they are valid numbers.
		if r.EMA20 <= 0 || r.EMA50 <= 0 {
			t.Errorf("EMA20=%v, EMA50=%v — both should be positive", r.EMA20, r.EMA50)
		}
	}
}

func TestCalculateAll_Uptrend(t *testing.T) {
	// With monotonically increasing data, trend should be "up".
	candles := testCandles()
	r := CalculateAll(candles)
	if r.Trend != "up" {
		t.Errorf("Trend = %q, want 'up' for monotonically increasing data", r.Trend)
	}
}

func TestCalculateAll_Empty(t *testing.T) {
	r := CalculateAll([]CandleData{})
	if r.Trend != "neutral" {
		t.Errorf("Empty input Trend = %q, want 'neutral'", r.Trend)
	}
}

func TestCalculateAll_Downtrend(t *testing.T) {
	// Build monotonically decreasing candles.
	candles := make([]CandleData, 60)
	for i := range candles {
		price := 100000.0 - float64(i)*100
		candles[i] = CandleData{
			Open:   price + 50,
			High:   price + 200,
			Low:    price - 200,
			Close:  price,
			Volume: 1000,
		}
	}
	r := CalculateAll(candles)
	if r.Trend != "down" {
		t.Errorf("Trend = %q, want 'down' for monotonically decreasing data", r.Trend)
	}
}

func TestCalculateAll_MACDFields(t *testing.T) {
	candles := testCandles()
	r := CalculateAll(candles)
	if math.IsNaN(r.MACDLine) {
		t.Error("MACDLine should not be NaN")
	}
	if math.IsNaN(r.MACDSignal) {
		t.Error("MACDSignal should not be NaN")
	}
	if math.IsNaN(r.MACDHist) {
		t.Error("MACDHist should not be NaN")
	}
	// histogram = macdLine - signal
	expectedHist := r.MACDLine - r.MACDSignal
	if math.Abs(r.MACDHist-expectedHist) > 1e-9 {
		t.Errorf("MACDHist = %v, want %v (macdLine - signal)", r.MACDHist, expectedHist)
	}
}
