package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/api/wallet"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var withdrawShowColumns = []output.TableColumn{
	{Header: "UUID", Key: "uuid"},
	{Header: i18n.T(i18n.HdrCurrency), Key: "currency"},
	{Header: i18n.T(i18n.HdrAmount), Key: "amount", Format: "number"},
	{Header: i18n.T(i18n.HdrFee), Key: "fee", Format: "number"},
	{Header: i18n.T(i18n.HdrState), Key: "state"},
	{Header: i18n.T(i18n.HdrTransType), Key: "transaction_type"},
	{Header: i18n.T(i18n.HdrTXID), Key: "txid"},
	{Header: i18n.T(i18n.HdrCreatedAt), Key: "created_at", Format: "datetime"},
	{Header: i18n.T(i18n.HdrDoneAt), Key: "done_at", Format: "datetime"},
}

var withdrawShowCmd = &cobra.Command{
	Use:   "show <uuid>",
	Short: i18n.T(i18n.MsgWithdrawShowShort),
	Args:  RequireArgs(1, i18n.T(i18n.ErrWithdrawUUIDRequired)),
	Example: `  upbit withdraw show 94332e99-3a87-4a35-ad98-28b0c969f830`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		wc := wallet.NewWalletClient(client)
		formatter := GetFormatterWithColumns(withdrawShowColumns)

		withdrawal, err := wc.GetWithdrawal(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return formatter.Format(withdrawal)
	},
}

func init() {
	withdrawCmd.AddCommand(withdrawShowCmd)
}
