package cli

import (
	"github.com/spf13/cobra"
)

var withdrawCmd = &cobra.Command{
	Use:     "withdraw",
	Short:   "출금 관리 (조회, 요청, 취소)",
	GroupID: "wallet",
	// 인자 없이 실행 시 도움말 표시
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(withdrawCmd)
}
