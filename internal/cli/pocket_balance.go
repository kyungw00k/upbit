package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/api/pocket"
	"github.com/kyungw00k/upbit/internal/i18n"
)

var pocketBalanceCmd = &cobra.Command{
	Use:   "balance <uuid>",
	Short: i18n.T(i18n.MsgPocketBalanceShort),
	Args:  RequireArgs(1, i18n.T(i18n.ErrPocketUUIDRequired)),
	Example: `  upbit pocket balance 9ca023a5-851b-4fec-9f0a-48cd83c2eaae   # 서브포켓 잔고
  upbit pocket list                                        # UUID 확인`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		pc := pocket.NewPocketClient(client)

		accounts, err := pc.SubPocketBalance(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return GetFormatterWithColumns(balanceColumns).Format(accounts)
	},
}

func init() {
	pocketCmd.AddCommand(pocketBalanceCmd)
}
