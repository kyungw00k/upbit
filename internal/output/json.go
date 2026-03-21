package output

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// JSONFormatter JSON 출력 포맷터
// tty: pretty print (indent 2), non-tty: compact (한 줄)
type JSONFormatter struct {
	pretty bool
}

// NewJSONFormatter JSONFormatter 생성
func NewJSONFormatter(pretty bool) *JSONFormatter {
	return &JSONFormatter{pretty: pretty}
}

// Format JSON으로 출력
func (f *JSONFormatter) Format(data interface{}) error {
	var b []byte
	var err error

	if f.pretty {
		b, err = json.MarshalIndent(data, "", "  ")
	} else {
		b, err = json.Marshal(data)
	}
	if err != nil {
		return fmt.Errorf("JSON 직렬화 실패: %w", err)
	}

	_, err = fmt.Fprintln(os.Stdout, string(b))
	return err
}

// JSONLFormatter JSONL (JSON Lines) 출력 포맷터
// 슬라이스의 각 항목을 한 줄씩 출력
type JSONLFormatter struct{}

// NewJSONLFormatter JSONLFormatter 생성
func NewJSONLFormatter() *JSONLFormatter {
	return &JSONLFormatter{}
}

// Format JSONL로 출력
func (f *JSONLFormatter) Format(data interface{}) error {
	// 슬라이스인 경우 항목별로 한 줄씩
	rv := reflect.ValueOf(data)
	if rv.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			b, err := json.Marshal(rv.Index(i).Interface())
			if err != nil {
				return fmt.Errorf("JSON 직렬화 실패 (index %d): %w", i, err)
			}
			fmt.Fprintln(os.Stdout, string(b))
		}
		return nil
	}

	// 슬라이스가 아닌 경우 단일 라인 출력
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("JSON 직렬화 실패: %w", err)
	}
	_, err = fmt.Fprintln(os.Stdout, string(b))
	return err
}

// JSONFieldsFormatter --json fields 처리 포맷터
// 지정된 필드만 추출하여 JSON 출력
type JSONFieldsFormatter struct {
	fields []string
	pretty bool
}

// NewJSONFieldsFormatter JSONFieldsFormatter 생성
// fields: 쉼표로 구분된 필드 목록
func NewJSONFieldsFormatter(fields string, pretty bool) *JSONFieldsFormatter {
	parts := strings.Split(fields, ",")
	trimmed := make([]string, 0, len(parts))
	for _, p := range parts {
		if f := strings.TrimSpace(p); f != "" {
			trimmed = append(trimmed, f)
		}
	}
	return &JSONFieldsFormatter{fields: trimmed, pretty: pretty}
}

// Format 지정 필드만 추출하여 JSON 출력
func (f *JSONFieldsFormatter) Format(data interface{}) error {
	extracted := f.extractFields(data)

	var b []byte
	var err error
	if f.pretty {
		b, err = json.MarshalIndent(extracted, "", "  ")
	} else {
		b, err = json.Marshal(extracted)
	}
	if err != nil {
		return fmt.Errorf("JSON 직렬화 실패: %w", err)
	}

	_, err = fmt.Fprintln(os.Stdout, string(b))
	return err
}

// extractFields 데이터에서 지정 필드를 추출
// 슬라이스인 경우 각 항목에서 필드를 추출한 슬라이스 반환
// 맵/구조체인 경우 단일 맵 반환
func (f *JSONFieldsFormatter) extractFields(data interface{}) interface{} {
	rv := reflect.ValueOf(data)
	if rv.Kind() == reflect.Slice {
		result := make([]map[string]interface{}, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			result = append(result, f.extractFromItem(rv.Index(i).Interface()))
		}
		return result
	}
	return f.extractFromItem(data)
}

// extractFromItem 단일 항목에서 지정 필드 추출
func (f *JSONFieldsFormatter) extractFromItem(item interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// map[string]interface{} 처리
	if m, ok := item.(map[string]interface{}); ok {
		for _, field := range f.fields {
			if v, exists := m[field]; exists {
				result[field] = v
			}
		}
		return result
	}

	// struct 처리: JSON 태그 또는 필드명 사용
	// 먼저 JSON으로 직렬화 후 map으로 변환
	b, err := json.Marshal(item)
	if err != nil {
		return result
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return result
	}

	for _, field := range f.fields {
		if v, exists := m[field]; exists {
			result[field] = v
		}
	}

	return result
}
