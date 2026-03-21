package output

import (
	"os"

	"golang.org/x/term"
)

// Formatter 출력 포맷터 인터페이스
type Formatter interface {
	Format(data interface{}) error
}

// IsTTY 현재 stdout이 tty인지 확인
func IsTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// NewFormatter 포맷터 생성
// output: "auto" | "table" | "json" | "jsonl" | "csv"
// jsonFields: 쉼표로 구분된 필드 목록 (비어있으면 전체 출력)
func NewFormatter(output string, jsonFields string) Formatter {
	return NewFormatterWithColumns(output, jsonFields, nil)
}

// NewFormatterWithColumns 컬럼 정의가 포함된 포맷터 생성
// columns가 지정되면 테이블 출력 시 해당 컬럼만 표시 (JSON/CSV는 전체 데이터 유지)
func NewFormatterWithColumns(output string, jsonFields string, columns []TableColumn) Formatter {
	isTTY := IsTTY()

	// --json fields 지정 시 JSONFieldsFormatter 우선
	if jsonFields != "" {
		return NewJSONFieldsFormatter(jsonFields, isTTY)
	}

	switch output {
	case "table":
		if len(columns) > 0 {
			return NewTableFormatterWithColumns(columns)
		}
		return NewTableFormatter()
	case "json":
		return NewJSONFormatter(isTTY)
	case "jsonl":
		return NewJSONLFormatter()
	case "csv":
		return NewCSVFormatter()
	case "auto", "":
		if isTTY {
			if len(columns) > 0 {
				return NewTableFormatterWithColumns(columns)
			}
			return NewTableFormatter()
		}
		return NewJSONFormatter(false)
	default:
		// 알 수 없는 포맷은 auto로 처리
		if isTTY {
			if len(columns) > 0 {
				return NewTableFormatterWithColumns(columns)
			}
			return NewTableFormatter()
		}
		return NewJSONFormatter(false)
	}
}
