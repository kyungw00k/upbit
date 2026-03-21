package output

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
)

// TableColumn 테이블 컬럼 정의
type TableColumn struct {
	Header string // 표시할 헤더 (예: "마켓", "현재가")
	Key    string // JSON 키 (예: "market", "trade_price")
	Format string // "number", "percent", "date", "string" (선택, 기본: 자동)
}

// TableFormatter 테이블 출력 포맷터
type TableFormatter struct {
	writer  *os.File
	columns []TableColumn
}

// NewTableFormatter TableFormatter 생성
func NewTableFormatter() *TableFormatter {
	return &TableFormatter{
		writer: os.Stdout,
	}
}

// NewTableFormatterWithColumns 컬럼이 지정된 TableFormatter 생성
func NewTableFormatterWithColumns(columns []TableColumn) *TableFormatter {
	return &TableFormatter{
		writer:  os.Stdout,
		columns: columns,
	}
}

// padRight 문자열을 표시 폭 기준으로 오른쪽 패딩
func padRight(s string, width int) string {
	sw := runewidth.StringWidth(s)
	if sw >= width {
		return s
	}
	return s + strings.Repeat(" ", width-sw)
}

// Format 테이블로 출력
func (f *TableFormatter) Format(data interface{}) error {
	// 컬럼이 지정된 경우 컬럼 기반 출력
	if len(f.columns) > 0 {
		return f.FormatWithColumns(data, f.columns)
	}

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
	case reflect.Map, reflect.Struct:
		return f.formatSingle(data)
	default:
		// 기본값: 그냥 출력
		fmt.Fprintln(f.writer, fmt.Sprintf("%v", data))
		return nil
	}
}

// FormatWithColumns 지정된 컬럼만 테이블로 출력
func (f *TableFormatter) FormatWithColumns(data interface{}, columns []TableColumn) error {
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
		return f.formatSliceWithColumns(rv, columns)
	case reflect.Map, reflect.Struct:
		return f.formatSingleWithColumns(data, columns)
	default:
		fmt.Fprintln(f.writer, fmt.Sprintf("%v", data))
		return nil
	}
}

// formatSliceWithColumns 슬라이스를 지정된 컬럼으로 테이블 출력
func (f *TableFormatter) formatSliceWithColumns(rv reflect.Value, columns []TableColumn) error {
	if rv.Len() == 0 {
		return nil
	}

	// 1. 모든 행의 데이터를 미리 계산
	allRows := make([][]string, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		row := rv.Index(i).Interface()
		m := toMap(row)
		values := make([]string, len(columns))
		if m != nil {
			for j, col := range columns {
				values[j] = formatColumnValue(m[col.Key], col.Format)
			}
		}
		allRows[i] = values
	}

	// 2. 헤더 포함 각 컬럼의 최대 표시 폭 계산
	colWidths := make([]int, len(columns))
	for j, col := range columns {
		colWidths[j] = runewidth.StringWidth(col.Header)
	}
	for _, row := range allRows {
		for j, val := range row {
			w := runewidth.StringWidth(val)
			if w > colWidths[j] {
				colWidths[j] = w
			}
		}
	}

	// 3. 헤더 출력 (패딩 적용)
	headerParts := make([]string, len(columns))
	for j, col := range columns {
		headerParts[j] = padRight(col.Header, colWidths[j])
	}
	fmt.Fprintln(f.writer, strings.Join(headerParts, "  "))

	// 4. 데이터 행 출력 (패딩 적용)
	for _, row := range allRows {
		parts := make([]string, len(row))
		for j, val := range row {
			parts[j] = padRight(val, colWidths[j])
		}
		fmt.Fprintln(f.writer, strings.Join(parts, "  "))
	}

	return nil
}

// formatSingleWithColumns 단일 항목을 지정된 컬럼으로 키-값 테이블 출력
func (f *TableFormatter) formatSingleWithColumns(data interface{}, columns []TableColumn) error {
	m := toMap(data)
	if m == nil {
		return nil
	}

	// Label 최대 폭 계산
	maxLabelWidth := 0
	for _, col := range columns {
		w := runewidth.StringWidth(col.Header)
		if w > maxLabelWidth {
			maxLabelWidth = w
		}
	}

	// Label: 값 출력
	for _, col := range columns {
		label := padRight(col.Header, maxLabelWidth)
		value := formatColumnValue(m[col.Key], col.Format)
		fmt.Fprintf(f.writer, "%s  %s\n", label, value)
	}

	return nil
}

// toMap 임의의 값을 map[string]interface{}로 변환
func toMap(item interface{}) map[string]interface{} {
	b, err := json.Marshal(item)
	if err != nil {
		return nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}
	return m
}

// KST 시간대
var kstLoc = time.FixedZone("KST", 9*60*60)

// formatColumnValue 컬럼 포맷에 따라 값을 문자열로 변환
func formatColumnValue(v interface{}, format string) string {
	if v == nil {
		return ""
	}

	switch format {
	case "percent":
		if f, ok := toFloat64(v); ok {
			return FormatPercent(f)
		}
		return fmt.Sprintf("%v", v)
	case "number":
		return formatValue(v)
	case "datetime":
		// ISO 8601 시간 문자열을 KST "2006-01-02 15:04:05" 포맷으로 변환
		return formatDateTimeKST(v)
	case "time":
		// HH:MM:SS 시간 문자열을 KST로 변환 (UTC 입력 가정)
		return formatTimeKST(v)
	case "date":
		// 날짜만 표시 (시간대 변환 없음)
		return fmt.Sprintf("%v", v)
	case "string":
		return fmt.Sprintf("%v", v)
	default:
		return formatValue(v)
	}
}

// formatDateTimeKST ISO 8601 시간을 KST로 변환하여 포맷
func formatDateTimeKST(v interface{}) string {
	s, ok := v.(string)
	if !ok {
		return fmt.Sprintf("%v", v)
	}
	// 여러 형식 시도
	formats := []string{
		time.RFC3339,                // 2006-01-02T15:04:05+09:00
		"2006-01-02T15:04:05",      // UTC (오프셋 없음)
		"2006-01-02T15:04:05Z",     // UTC 명시
		"2006-01-02 15:04:05",      // 공백 구분
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t.In(kstLoc).Format("2006-01-02 15:04:05")
		}
	}
	return s
}

// formatTimeKST UTC HH:MM:SS를 KST로 변환
func formatTimeKST(v interface{}) string {
	s, ok := v.(string)
	if !ok {
		return fmt.Sprintf("%v", v)
	}
	// HH:MM:SS 형식 파싱
	t, err := time.Parse("15:04:05", s)
	if err != nil {
		return s
	}
	// UTC로 해석하여 KST로 변환
	utc := time.Date(0, 1, 1, t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
	return utc.In(kstLoc).Format("15:04:05")
}

// toFloat64 interface{}를 float64로 변환
func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case string:
		f, err := strconv.ParseFloat(val, 64)
		return f, err == nil
	default:
		return 0, false
	}
}

// ToMapSlice 슬라이스를 []map[string]interface{}로 변환 (외부 패키지용)
func ToMapSlice(data interface{}) []map[string]interface{} {
	rv := reflect.ValueOf(data)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Slice {
		return nil
	}
	result := make([]map[string]interface{}, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		m := toMap(rv.Index(i).Interface())
		if m != nil {
			result = append(result, m)
		}
	}
	return result
}

// FormatNumberPublic 숫자를 사람이 읽기 좋은 문자열로 변환 (외부 패키지용)
func FormatNumberPublic(f float64) string {
	return formatNumber(f)
}

// formatSlice 슬라이스를 테이블로 출력
func (f *TableFormatter) formatSlice(rv reflect.Value) error {
	if rv.Len() == 0 {
		return nil
	}

	// 첫 번째 요소로부터 컬럼 정보 추출
	first := rv.Index(0).Interface()
	columns, jsonKeys := extractColumns(first)

	if len(columns) == 0 {
		return nil
	}

	// 모든 행의 데이터를 미리 계산
	allRows := make([][]string, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		row := rv.Index(i).Interface()
		values := extractValues(row, jsonKeys)
		formatted := make([]string, len(values))
		for j, v := range values {
			formatted[j] = formatValue(v)
		}
		allRows[i] = formatted
	}

	// 헤더 포함 각 컬럼의 최대 표시 폭 계산
	colWidths := make([]int, len(columns))
	for j, col := range columns {
		colWidths[j] = runewidth.StringWidth(col)
	}
	for _, row := range allRows {
		for j, val := range row {
			w := runewidth.StringWidth(val)
			if w > colWidths[j] {
				colWidths[j] = w
			}
		}
	}

	// 헤더 출력 (패딩 적용)
	headerParts := make([]string, len(columns))
	for j, col := range columns {
		headerParts[j] = padRight(col, colWidths[j])
	}
	fmt.Fprintln(f.writer, strings.Join(headerParts, "  "))

	// 데이터 행 출력 (패딩 적용)
	for _, row := range allRows {
		parts := make([]string, len(row))
		for j, val := range row {
			parts[j] = padRight(val, colWidths[j])
		}
		fmt.Fprintln(f.writer, strings.Join(parts, "  "))
	}

	return nil
}

// formatSingle 단일 항목을 키-값 테이블로 출력
func (f *TableFormatter) formatSingle(data interface{}) error {
	// JSON으로 직렬화 후 맵으로 변환
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("직렬화 실패: %w", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return fmt.Errorf("역직렬화 실패: %w", err)
	}

	// 키의 최대 표시 폭 계산
	maxKeyWidth := 0
	for k := range m {
		w := runewidth.StringWidth(k)
		if w > maxKeyWidth {
			maxKeyWidth = w
		}
	}

	for k, v := range m {
		fmt.Fprintf(f.writer, "%s  %s\n", padRight(k, maxKeyWidth), formatValue(v))
	}

	return nil
}

// extractColumns 항목에서 컬럼 헤더와 JSON 키를 추출
func extractColumns(item interface{}) (headers []string, jsonKeys []string) {
	// JSON으로 변환 후 키 목록 추출
	b, err := json.Marshal(item)
	if err != nil {
		return nil, nil
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, nil
	}

	// struct 필드 순서 유지를 위해 reflect 사용
	rv := reflect.ValueOf(item)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() == reflect.Struct {
		rt := rv.Type()
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			jsonTag := field.Tag.Get("json")
			if jsonTag == "" || jsonTag == "-" {
				continue
			}
			// omitempty 제거
			jsonKey := strings.Split(jsonTag, ",")[0]
			if jsonKey == "-" {
				continue
			}
			// 해당 키가 실제 데이터에 있는 경우만 포함
			if _, exists := m[jsonKey]; exists {
				headers = append(headers, strings.ToUpper(jsonKey))
				jsonKeys = append(jsonKeys, jsonKey)
			}
		}
		return headers, jsonKeys
	}

	// map인 경우 키 목록 반환
	for k := range m {
		headers = append(headers, strings.ToUpper(k))
		jsonKeys = append(jsonKeys, k)
	}
	return headers, jsonKeys
}

// extractValues 항목에서 지정된 키의 값을 추출
func extractValues(item interface{}, keys []string) []interface{} {
	b, err := json.Marshal(item)
	if err != nil {
		result := make([]interface{}, len(keys))
		return result
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		result := make([]interface{}, len(keys))
		return result
	}

	result := make([]interface{}, len(keys))
	for i, k := range keys {
		result[i] = m[k]
	}
	return result
}

// formatValue 값을 사람이 읽기 좋은 문자열로 변환
// 숫자: 천 단위 쉼표
// float 변화율: 부호 포함 퍼센트
func formatValue(v interface{}) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case float64:
		return formatNumber(val)
	case int:
		return formatInt(int64(val))
	case int64:
		return formatInt(val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	case string:
		return val
	case []interface{}:
		// 배열: 각 요소를 formatValue로 재귀 처리
		parts := make([]string, len(val))
		for i, elem := range val {
			parts[i] = formatValue(elem)
		}
		return strings.Join(parts, ", ")
	default:
		return fmt.Sprintf("%v", v)
	}
}

// formatNumber 숫자를 천 단위 쉼표로 포맷
// 소수점이 있는 경우 소수점 유지, 정수인 경우 쉼표만 적용
func formatNumber(f float64) string {
	// 정수 여부 확인
	if f == math.Trunc(f) && math.Abs(f) < 1e15 {
		return formatInt(int64(f))
	}

	// 소수점 포함: 적절한 소수점 자리수로 표현
	formatted := strconv.FormatFloat(f, 'f', -1, 64)
	parts := strings.Split(formatted, ".")
	if len(parts) == 2 {
		return formatIntStr(parts[0]) + "." + parts[1]
	}
	return formatted
}

// formatInt 정수를 천 단위 쉼표로 포맷
func formatInt(n int64) string {
	return formatIntStr(strconv.FormatInt(n, 10))
}

// formatIntStr 정수 문자열에 천 단위 쉼표 추가
func formatIntStr(s string) string {
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = s[1:]
	}

	n := len(s)
	if n <= 3 {
		if neg {
			return "-" + s
		}
		return s
	}

	var b strings.Builder
	mod := n % 3
	if mod == 0 {
		mod = 3
	}
	b.WriteString(s[:mod])
	for i := mod; i < n; i += 3 {
		b.WriteByte(',')
		b.WriteString(s[i : i+3])
	}

	if neg {
		return "-" + b.String()
	}
	return b.String()
}

// FormatPercent 변화율을 +/-부호 포함 퍼센트로 포맷
func FormatPercent(rate float64) string {
	if rate > 0 {
		return fmt.Sprintf("+%.2f%%", rate*100)
	} else if rate < 0 {
		return fmt.Sprintf("%.2f%%", rate*100)
	}
	return "0.00%"
}
