package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestFloat64_UnmarshalJSON_Number(t *testing.T) {
	input := `{"value": 123.456}`
	var result struct {
		Value Float64 `json:"value"`
	}
	err := json.Unmarshal([]byte(input), &result)
	if err != nil {
		t.Fatalf("파싱 실패: %v", err)
	}
	if float64(result.Value) != 123.456 {
		t.Errorf("기대: 123.456, 실제: %f", float64(result.Value))
	}
}

func TestFloat64_UnmarshalJSON_String(t *testing.T) {
	input := `{"value": "789.012"}`
	var result struct {
		Value Float64 `json:"value"`
	}
	err := json.Unmarshal([]byte(input), &result)
	if err != nil {
		t.Fatalf("문자열 파싱 실패: %v", err)
	}
	if float64(result.Value) != 789.012 {
		t.Errorf("기대: 789.012, 실제: %f", float64(result.Value))
	}
}

func TestFloat64_UnmarshalJSON_EmptyString(t *testing.T) {
	input := `{"value": ""}`
	var result struct {
		Value Float64 `json:"value"`
	}
	err := json.Unmarshal([]byte(input), &result)
	if err != nil {
		t.Fatalf("빈 문자열 파싱 실패: %v", err)
	}
	if float64(result.Value) != 0 {
		t.Errorf("빈 문자열은 0이어야 함, 실제: %f", float64(result.Value))
	}
}

func TestFloat64_UnmarshalJSON_Zero(t *testing.T) {
	input := `{"value": 0}`
	var result struct {
		Value Float64 `json:"value"`
	}
	err := json.Unmarshal([]byte(input), &result)
	if err != nil {
		t.Fatalf("0 파싱 실패: %v", err)
	}
	if float64(result.Value) != 0 {
		t.Errorf("기대: 0, 실제: %f", float64(result.Value))
	}
}

func TestFloat64_UnmarshalJSON_InvalidString(t *testing.T) {
	input := `{"value": "not-a-number"}`
	var result struct {
		Value Float64 `json:"value"`
	}
	err := json.Unmarshal([]byte(input), &result)
	if err == nil {
		t.Error("유효하지 않은 문자열에서 에러가 반환되어야 함")
	}
}

func TestFloat64_String(t *testing.T) {
	tests := []struct {
		input    Float64
		expected string
	}{
		{Float64(123.456), "123.456"},
		{Float64(0), "0"},
		{Float64(-99.9), "-99.9"},
		{Float64(1000000), "1000000"},
		{Float64(0.001), "0.001"},
	}

	for _, tt := range tests {
		got := tt.input.String()
		if got != tt.expected {
			t.Errorf("Float64(%f).String() = %s, 기대: %s", float64(tt.input), got, tt.expected)
		}
	}
}

func TestFloat64_MarshalJSON(t *testing.T) {
	f := Float64(42.5)
	data, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("Marshal 실패: %v", err)
	}
	if string(data) != "42.5" {
		t.Errorf("기대: 42.5, 실제: %s", string(data))
	}
}

func TestFloat64_Float64Value(t *testing.T) {
	f := Float64(99.99)
	if f.Float64Value() != 99.99 {
		t.Errorf("기대: 99.99, 실제: %f", f.Float64Value())
	}
}

func TestTimestamp_UnmarshalJSON_ISO8601(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"RFC3339", `"2024-01-15T10:30:00+09:00"`},
		{"UTC Z", `"2024-01-15T01:30:00Z"`},
		{"No timezone", `"2024-01-15T10:30:00"`},
		{"Space separator", `"2024-01-15 10:30:00"`},
		{"Millisecond Z", `"2024-01-15T01:30:00.000Z"`},
		{"Millisecond no tz", `"2024-01-15T10:30:00.000"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ts Timestamp
			err := json.Unmarshal([]byte(tt.input), &ts)
			if err != nil {
				t.Fatalf("파싱 실패: %v", err)
			}

			parsed := time.Time(ts)
			if parsed.IsZero() {
				t.Error("파싱된 시간이 zero value")
			}
		})
	}
}

func TestTimestamp_UnmarshalJSON_UnixMillis(t *testing.T) {
	// 2024-01-15T01:30:00Z in milliseconds
	ms := int64(1705282200000)
	input := `1705282200000`

	var ts Timestamp
	err := json.Unmarshal([]byte(input), &ts)
	if err != nil {
		t.Fatalf("Unix ms 파싱 실패: %v", err)
	}

	expected := time.UnixMilli(ms).UTC()
	got := time.Time(ts)
	if !got.Equal(expected) {
		t.Errorf("기대: %v, 실제: %v", expected, got)
	}
}

func TestTimestamp_UnmarshalJSON_InvalidFormat(t *testing.T) {
	input := `"not-a-timestamp"`
	var ts Timestamp
	err := json.Unmarshal([]byte(input), &ts)
	if err == nil {
		t.Error("유효하지 않은 포맷에서 에러가 반환되어야 함")
	}
}

func TestTimestamp_MarshalJSON(t *testing.T) {
	ts := Timestamp(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC))
	data, err := json.Marshal(ts)
	if err != nil {
		t.Fatalf("Marshal 실패: %v", err)
	}

	expected := `"2024-01-15T10:30:00Z"`
	if string(data) != expected {
		t.Errorf("기대: %s, 실제: %s", expected, string(data))
	}
}

func TestTimestamp_String(t *testing.T) {
	ts := Timestamp(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC))
	got := ts.String()
	expected := "2024-01-15T10:30:00Z"
	if got != expected {
		t.Errorf("기대: %s, 실제: %s", expected, got)
	}
}

func TestTimestamp_Time(t *testing.T) {
	original := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	ts := Timestamp(original)
	if !ts.Time().Equal(original) {
		t.Errorf("Time() 반환값이 원래 값과 다름")
	}
}
