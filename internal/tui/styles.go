package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// StyleRise 상승 (초록)
	StyleRise = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	// StyleFall 하락 (빨강)
	StyleFall = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	// StyleFlat 보합 (흰색)
	StyleFlat = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	// StyleHeader 헤더
	StyleHeader = lipgloss.NewStyle().Bold(true).Underline(true)
	// StyleBar 바 차트 (파랑)
	StyleBar = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	// StyleHint 하단 힌트
	StyleHint = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	// StyleTitle 제목
	StyleTitle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	// StyleBox 확인 프롬프트 박스
	StyleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("4")).
			Padding(1, 2)
	// StyleHighlight 강조 텍스트
	StyleHighlight = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("3"))
)

// PriceStyle 변동 방향에 따른 스타일 반환
func PriceStyle(change string) lipgloss.Style {
	switch change {
	case "RISE":
		return StyleRise
	case "FALL":
		return StyleFall
	default:
		return StyleFlat
	}
}

// SideStyle BID/ASK에 따른 스타일 반환
func SideStyle(askBid string) lipgloss.Style {
	if askBid == "ASK" {
		return StyleFall
	}
	return StyleRise
}

// TruncateToHeight 터미널 높이 초과 시 위에서부터 유지하고 아래를 자름
func TruncateToHeight(s string, height int) string {
	if height <= 0 {
		return s
	}
	lines := strings.Split(s, "\n")
	if len(lines) > height {
		lines = lines[:height]
	}
	return strings.Join(lines, "\n")
}
