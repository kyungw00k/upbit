package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/wallet"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var withdrawCancelCmd = &cobra.Command{
	Use:   "cancel <uuid>",
	Short: i18n.T(i18n.MsgWithdrawCancelShort),
	Args:  RequireArgs(1, i18n.T(i18n.ErrWithdrawUUIDRequired)),
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

		msg := i18n.Tf(i18n.MsgWithdrawCancelConfirm, uuid)
		confirmed, err := output.Confirm(msg, GetForce())
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgCancelAborted))
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
