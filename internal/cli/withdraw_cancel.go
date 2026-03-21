package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/wallet"
	"github.com/kyungw00k/upbit/internal/output"
)

var withdrawCancelCmd = &cobra.Command{
	Use:   "cancel <uuid>",
	Short: "출금 취소",
	Args:  RequireArgs(1, "출금 UUID를 지정하세요"),
	Example: `  upbit withdraw cancel 9f432943-54e0-40b7-825f-b6fec8b42b79
  upbit withdraw cancel 9f432943-54e0-40b7-825f-b6fec8b42b79 --force`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		wc := wallet.NewWalletClient(client)
		formatter := GetFormatterWithColumns(withdrawResultColumns)

		uuid := args[0]

		msg := fmt.Sprintf("출금 %s을(를) 취소합니다", uuid)
		confirmed, err := output.Confirm(msg, GetForce())
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Fprintln(os.Stderr, "취소가 중단되었습니다")
			return nil
		}

		withdrawal, err := wc.CancelWithdrawal(cmd.Context(), uuid)
		if err != nil {
			return err
		}
		return formatter.Format(withdrawal)
	},
}

func init() {
	AddForceFlag(withdrawCancelCmd)
	withdrawCmd.AddCommand(withdrawCancelCmd)
}
