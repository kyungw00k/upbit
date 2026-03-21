package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/wallet"
	"github.com/kyungw00k/upbit/internal/output"
)

var depositShowColumns = []output.TableColumn{
	{Header: "UUID", Key: "uuid"},
	{Header: "통화", Key: "currency"},
	{Header: "금액", Key: "amount", Format: "number"},
	{Header: "수수료", Key: "fee", Format: "number"},
	{Header: "상태", Key: "state"},
	{Header: "유형", Key: "transaction_type"},
	{Header: "TXID", Key: "txid"},
	{Header: "생성일", Key: "created_at", Format: "datetime"},
	{Header: "완료일", Key: "done_at", Format: "datetime"},
}

var depositShowCmd = &cobra.Command{
	Use:   "show <uuid>",
	Short: "개별 입금 조회",
	Args:  RequireArgs(1, "입금 UUID를 지정하세요"),
	Example: `  upbit deposit show 94332e99-3a87-4a35-ad98-28b0c969f830`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		wc := wallet.NewWalletClient(client)
		formatter := GetFormatterWithColumns(depositShowColumns)

		deposit, err := wc.GetDeposit(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return formatter.Format(deposit)
	},
}

func init() {
	depositCmd.AddCommand(depositShowCmd)
}
