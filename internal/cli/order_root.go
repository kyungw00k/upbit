package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/i18n"
)

var orderCmd = &cobra.Command{
	Use:     "order",
	Short:   i18n.T(i18n.MsgOrderShort),
	GroupID: "trading",
	// 인자 없이 실행 시 도움말 표시
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(orderCmd)
}
