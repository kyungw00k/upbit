package output

import (
	"testing"
)

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{999, "999"},
		{1000, "1,000"},
		{1234567, "1,234,567"},
		{-1234567, "-1,234,567"},
		{1234.5678, "1,234.5678"},
		{-1234.5678, "-1,234.5678"},
		{0.001, "0.001"},
		{100000000, "100,000,000"},
	}

	for _, tt := range tests {
		got := formatNumber(tt.input)
		if got != tt.expected {
			t.Errorf("formatNumber(%v) = %s, 기대: %s", tt.input, got, tt.expected)
		}
	}
}

func TestFormatInt(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{-1, "-1"},
		{999, "999"},
		{1000, "1,000"},
		{999999, "999,999"},
		{1000000, "1,000,000"},
		{-1000000, "-1,000,000"},
	}

	for _, tt := range tests {
		got := formatInt(tt.input)
		if got != tt.expected {
			t.Errorf("formatInt(%d) = %s, 기대: %s", tt.input, got, tt.expected)
		}
	}
}

func TestFormatPercent(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0.05, "+5.00%"},
		{-0.03, "-3.00%"},
		{0, "0.00%"},
		{0.1234, "+12.34%"},
		{-0.001, "-0.10%"},
		{1.0, "+100.00%"},
	}

	for _, tt := range tests {
		got := FormatPercent(tt.input)
		if got != tt.expected {
			t.Errorf("FormatPercent(%v) = %s, 기대: %s", tt.input, got, tt.expected)
		}
	}
}

func TestFormatValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"nil", nil, ""},
		{"float64", float64(1234.5), "1,234.5"},
		{"int", 1000, "1,000"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"string", "hello", "hello"},
		{"array", []interface{}{"a", "b"}, "a, b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatValue(tt.input)
			if got != tt.expected {
				t.Errorf("formatValue(%v) = %s, 기대: %s", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFormatColumnValue(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		format   string
		expected string
	}{
		{"nil", nil, "number", ""},
		{"number format", float64(1234), "number", "1,234"},
		{"percent format", float64(0.05), "percent", "+5.00%"},
		{"string format", "test", "string", "test"},
		{"date format", "2024-01-15", "date", "2024-01-15"},
		{"datetime format RFC3339", "2024-01-15T01:30:00+09:00", "datetime", "2024-01-15 01:30:00"},
		{"default format", float64(999), "", "999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatColumnValue(tt.value, tt.format)
			if got != tt.expected {
				t.Errorf("formatColumnValue(%v, %q) = %s, 기대: %s", tt.value, tt.format, got, tt.expected)
			}
		})
	}
}

func TestFormatDateTimeKST(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"RFC3339 KST", "2024-01-15T10:30:00+09:00", "2024-01-15 10:30:00"},
		{"UTC Z", "2024-01-15T01:30:00Z", "2024-01-15 10:30:00"},
		{"No timezone", "2024-01-15T10:30:00", "2024-01-15 19:30:00"},
		{"Space separator", "2024-01-15 10:30:00", "2024-01-15 19:30:00"},
		{"Non-string", 12345, "12345"},
		{"Unparseable", "not-a-date", "not-a-date"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDateTimeKST(tt.input)
			if got != tt.expected {
				t.Errorf("formatDateTimeKST(%v) = %s, 기대: %s", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFormatTimeKST(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"UTC 01:30", "01:30:00", "10:30:00"},
		{"UTC 15:00", "15:00:00", "00:00:00"},
		{"Non-string", 123, "123"},
		{"Invalid format", "invalid", "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatTimeKST(tt.input)
			if got != tt.expected {
				t.Errorf("formatTimeKST(%v) = %s, 기대: %s", tt.input, got, tt.expected)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected float64
		ok       bool
	}{
		{"float64", float64(1.5), 1.5, true},
		{"int", int(42), 42, true},
		{"int64", int64(100), 100, true},
		{"string", "3.14", 3.14, true},
		{"invalid string", "abc", 0, false},
		{"bool", true, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := toFloat64(tt.input)
			if ok != tt.ok {
				t.Errorf("toFloat64(%v) ok = %v, 기대: %v", tt.input, ok, tt.ok)
			}
			if ok && got != tt.expected {
				t.Errorf("toFloat64(%v) = %f, 기대: %f", tt.input, got, tt.expected)
			}
		})
	}
}

func TestToMap(t *testing.T) {
	type sample struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	m := toMap(sample{Name: "test", Value: 42})
	if m == nil {
		t.Fatal("toMap이 nil을 반환")
	}
	if m["name"] != "test" {
		t.Errorf("name 기대: test, 실제: %v", m["name"])
	}
	// JSON unmarshal은 숫자를 float64로 변환
	if m["value"] != float64(42) {
		t.Errorf("value 기대: 42, 실제: %v", m["value"])
	}
}

func TestFormatNumberPublic(t *testing.T) {
	got := FormatNumberPublic(1234567.89)
	expected := "1,234,567.89"
	if got != expected {
		t.Errorf("FormatNumberPublic(1234567.89) = %s, 기대: %s", got, expected)
	}
}
