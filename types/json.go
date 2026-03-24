package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// KSTLoc is the Korea Standard Time zone (UTC+9).
var KSTLoc = time.FixedZone("KST", 9*60*60)

// Float64 is a type that parses both JSON strings and numbers.
// It handles cases where the Upbit API returns numeric values as strings.
type Float64 float64

// String implements fmt.Stringer.
func (f Float64) String() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

// UnmarshalJSON parses a JSON number or string into Float64.
func (f *Float64) UnmarshalJSON(data []byte) error {
	// Try numeric first.
	var n float64
	if err := json.Unmarshal(data, &n); err == nil {
		*f = Float64(n)
		return nil
	}

	// Fall back to string.
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Float64 parse failed: %s", string(data))
	}

	// Handle empty string.
	if strings.TrimSpace(s) == "" {
		*f = 0
		return nil
	}

	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("Float64 string parse failed: %s", s)
	}
	*f = Float64(n)
	return nil
}

// MarshalJSON serializes Float64 as a JSON number.
func (f Float64) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(f))
}

// Float64Value returns the underlying float64 value.
func (f Float64) Float64Value() float64 {
	return float64(f)
}

// Timestamp is a type that parses multiple time formats.
// It handles ISO 8601, Unix milliseconds, and other formats returned by the Upbit API.
type Timestamp time.Time

var timestampFormats = []string{
	time.RFC3339,
	"2006-01-02T15:04:05+09:00",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05.000Z",
	"2006-01-02T15:04:05.000",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
}

// UnmarshalJSON parses a JSON timestamp using multiple supported formats.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Try string formats first.
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		for _, format := range timestampFormats {
			if parsed, err := time.Parse(format, s); err == nil {
				*t = Timestamp(parsed)
				return nil
			}
		}
		return fmt.Errorf("Timestamp format not recognized: %s", s)
	}

	// Try Unix millisecond timestamp.
	var ms int64
	if err := json.Unmarshal(data, &ms); err == nil {
		*t = Timestamp(time.UnixMilli(ms).UTC())
		return nil
	}

	return fmt.Errorf("Timestamp parse failed: %s", string(data))
}

// MarshalJSON serializes Timestamp as an ISO 8601 (RFC3339) string.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).UTC().Format(time.RFC3339) + `"`), nil
}

// Time returns the underlying time.Time value.
func (t Timestamp) Time() time.Time {
	return time.Time(t)
}

// String returns the timestamp in a human-readable RFC3339 format.
func (t Timestamp) String() string {
	return time.Time(t).Format(time.RFC3339)
}
