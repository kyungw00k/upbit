package output

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kyungw00k/upbit/internal/i18n"
)

// Confirm 사용자에게 확인 프롬프트를 표시하고 결과를 반환
//
// force=true: 즉시 true 반환
// non-tty:   즉시 false 반환 + stderr에 "--force 필요" 메시지
// tty:       stdin에서 Y/n 입력 대기, 기본값 Y (Enter 입력 시 true)
func Confirm(message string, force bool) (bool, error) {
	// --force 플래그가 설정된 경우 즉시 승인
	if force {
		return true, nil
	}

	// non-tty 환경에서는 자동으로 거부 (안전 장치)
	if !IsTTY() {
		fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgConfirmNonTTY))
		return false, nil
	}

	// tty 환경에서 사용자 입력 대기
	fmt.Fprintf(os.Stderr, "%s [Y/n]: ", message)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("%s: %w", i18n.T(i18n.ErrInputRead), err)
	}

	input = strings.TrimSpace(strings.ToLower(input))

	// 기본값 Y: 빈 입력 또는 'y' 입력 시 true
	switch input {
	case "", "y", "yes":
		return true, nil
	default:
		return false, nil
	}
}
