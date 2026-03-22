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
// tty:       stdin에서 Y/n 입력 대기
func Confirm(message string, force bool) (bool, error) {
	if force {
		return true, nil
	}

	if !IsTTY() {
		fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgConfirmNonTTY))
		return false, nil
	}

	fmt.Fprintf(os.Stderr, "%s [Y/n]: ", message)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("%s: %w", i18n.T(i18n.ErrInputRead), err)
	}

	input = strings.TrimSpace(strings.ToLower(input))

	switch input {
	case "", "y", "yes":
		return true, nil
	default:
		return false, nil
	}
}
