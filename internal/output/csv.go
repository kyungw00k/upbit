package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
)

// CSVFormatter CSV 출력 포맷터
type CSVFormatter struct {
	writer *csv.Writer
}

// NewCSVFormatter CSVFormatter 생성
func NewCSVFormatter() *CSVFormatter {
	return &CSVFormatter{
		writer: csv.NewWriter(os.Stdout),
	}
}

// Format CSV로 출력
func (f *CSVFormatter) Format(data interface{}) error {
	rv := reflect.ValueOf(data)

	// 포인터 역참조
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Slice:
		return f.formatSlice(rv)
	default:
		return f.formatSingle(data)
	}
}

// formatSlice 슬라이스를 CSV로 출력 (헤더 행 + 데이터 행)
func (f *CSVFormatter) formatSlice(rv reflect.Value) error {
	if rv.Len() == 0{
		return nil
	}

	first := rv.Index(0).Interface()
	headers, jsonKeys := extractColumns(first)

	if len(headers) == 0 {
		return nil
	}

	// 헤더 행 출력
	if err := f.writer.Write(headers); err != nil {
		return fmt.Errorf("CSV 헤더 출력 실패: %w", err)
	}

	// 데이터 행 출력
	for i := 0; i < rv.Len(); i++ {
		row := rv.Index(i).Interface()
		values := extractValues(row, jsonKeys)
		strValues := make([]string, len(values))
		for j, v := range values {
			strValues[j] = csvFormatValue(v)
		}
		if err := f.writer.Write(strValues); err != nil {
			return fmt.Errorf("CSV 행 출력 실패 (index %d): %w", i, err)
		}
	}

	f.writer.Flush()
	return f.writer.Error()
}

// csvFormatValue CSV용 값 포맷 (과학적 표기법 방지)
func csvFormatValue(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case float64:
		if val == math.Trunc(val) && math.Abs(val) < 1e18 {
			return strconv.FormatInt(int64(val), 10)
		}
		return strconv.FormatFloat(val, 'f', -1, 64)
	case string:
		return val
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// formatSingle 단일 항목을 CSV (키, 값) 형태로 출력
func (f *CSVFormatter) formatSingle(data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("직렬화 실패: %w", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return fmt.Errorf("역직렬화 실패: %w", err)
	}

	// 헤더
	headers := []string{"key", "value"}
	if err := f.writer.Write(headers); err != nil {
		return fmt.Errorf("CSV 헤더 출력 실패: %w", err)
	}

	for k, v := range m {
		row := []string{k, csvFormatValue(v)}
		if err := f.writer.Write(row); err != nil {
			return fmt.Errorf("CSV 행 출력 실패: %w", err)
		}
	}

	f.writer.Flush()
	return f.writer.Error()
}
