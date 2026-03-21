package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Float64 JSON 문자열/숫자 양방향 파싱 타입
// Upbit API가 숫자를 문자열로 반환하는 경우 대응
type Float64 float64

// String fmt.Stringer 구현
func (f Float64) String() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

// UnmarshalJSON JSON 숫자 또는 문자열로부터 파싱
func (f *Float64) UnmarshalJSON(data []byte) error {
	// 숫자 우선 시도
	var n float64
	if err := json.Unmarshal(data, &n); err == nil {
		*f = Float64(n)
		return nil
	}

	// 문자열로 시도
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Float64 파싱 실패: %s", string(data))
	}

	// 빈 문자열 처리
	if strings.TrimSpace(s) == "" {
		*f = 0
		return nil
	}

	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("Float64 문자열 파싱 실패: %s", s)
	}
	*f = Float64(n)
	return nil
}

// MarshalJSON JSON 숫자로 직렬화
func (f Float64) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(f))
}

// Float64Value 기본 float64 값 반환
func (f Float64) Float64Value() float64 {
	return float64(f)
}

// Timestamp 다중 시간 포맷 파싱 타입
// ISO 8601, Unix ms 등 다양한 Upbit API 응답 형식 처리
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

// UnmarshalJSON JSON 다중 포맷으로부터 시간 파싱
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// 문자열 포맷 우선 시도
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		for _, format := range timestampFormats {
			if parsed, err := time.Parse(format, s); err == nil {
				*t = Timestamp(parsed)
				return nil
			}
		}
		return fmt.Errorf("Timestamp 포맷 인식 불가: %s", s)
	}

	// Unix 밀리초 타임스탬프 시도
	var ms int64
	if err := json.Unmarshal(data, &ms); err == nil {
		*t = Timestamp(time.UnixMilli(ms).UTC())
		return nil
	}

	return fmt.Errorf("Timestamp 파싱 실패: %s", string(data))
}

// MarshalJSON ISO 8601 (RFC3339) 포맷으로 직렬화
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).UTC().Format(time.RFC3339) + `"`), nil
}

// Time 내부 time.Time 값 반환
func (t Timestamp) Time() time.Time {
	return time.Time(t)
}

// String 사람이 읽기 좋은 형식으로 반환
func (t Timestamp) String() string {
	return time.Time(t).Format(time.RFC3339)
}
