package indicator

import "math"

// SMA computes the Simple Moving Average over the given period.
// Returns a slice of the same length as input. Positions before a full
// window is available are filled with NaN.
func SMA(values []float64, period int) []float64 {
	n := len(values)
	result := make([]float64, n)
	if period <= 0 || n == 0 {
		for i := range result {
			result[i] = math.NaN()
		}
		return result
	}

	for i := 0; i < n; i++ {
		if i < period-1 {
			result[i] = math.NaN()
			continue
		}
		var sum float64
		for j := i - period + 1; j <= i; j++ {
			sum += values[j]
		}
		result[i] = sum / float64(period)
	}
	return result
}

// StdDev returns the population standard deviation of values.
// Returns 0 for empty input.
func StdDev(values []float64) float64 {
	n := len(values)
	if n == 0 {
		return 0
	}

	var mean float64
	for _, v := range values {
		mean += v
	}
	mean /= float64(n)

	var variance float64
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(n)

	return math.Sqrt(variance)
}

// TrueRange calculates the True Range for a single candle.
func TrueRange(high, low, prevClose float64) float64 {
	tr1 := high - low
	tr2 := math.Abs(high - prevClose)
	tr3 := math.Abs(low - prevClose)
	return math.Max(tr1, math.Max(tr2, tr3))
}
