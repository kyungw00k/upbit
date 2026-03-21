package cli

import (
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/i18n"
)

// emptyMessage 빈 결과일 때 stderr에 안내 메시지 출력 후 true 반환
func emptyMessage(data interface{}, msg string) bool {
	rv := reflect.ValueOf(data)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			fmt.Fprintln(os.Stderr, msg)
			return true
		}
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Slice && rv.Len() == 0 {
		fmt.Fprintln(os.Stderr, msg)
		return true
	}
	return false
}

// RequireArgs 정확히 n개의 인자를 요구하는 커스텀 Args 함수
func RequireArgs(n int, usage string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < n {
			return fmt.Errorf("%s\n\n%s: %s", usage, i18n.T(i18n.MsgUsagePrefix), cmd.CommandPath())
		}
		return nil
	}
}

// RequireMinArgs 최소 n개의 인자를 요구하는 커스텀 Args 함수
func RequireMinArgs(n int, usage string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < n {
			return fmt.Errorf("%s\n\n%s: %s", usage, i18n.T(i18n.MsgUsagePrefix), cmd.CommandPath())
		}
		return nil
	}
}
