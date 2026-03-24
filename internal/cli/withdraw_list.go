package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/api/wallet"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var withdrawListColumns = []output.TableColumn{
	{Header: "UUID", Key: "uuid"},
	{Header: i18n.T(i18n.HdrCurrency), Key: "currency"},
	{Header: i18n.T(i18n.HdrAmount), Key: "amount", Format: "number"},
	{Header: i18n.T(i18n.HdrFee), Key: "fee", Format: "number"},
	{Header: i18n.T(i18n.HdrState), Key: "state"},
	{Header: i18n.T(i18n.HdrTransType), Key: "transaction_type"},
	{Header: i18n.T(i18n.HdrCreatedAt), Key: "created_at", Format: "datetime"},
}

var withdrawListCmd = &cobra.Command{
	Use:   "list",
	Short: i18n.T(i18n.MsgWithdrawListShort),
	Example: `  upbit withdraw list                        # 최근 출금 목록
  upbit withdraw list --currency BTC         # BTC 출금만
  upbit withdraw list --state DONE           # 완료된 출금만
  upbit withdraw list -c 50                  # 50개 조회
  upbit withdraw list --page 2               # 2페이지`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		wc := wallet.NewWalletClient(client)
		formatter := GetFormatterWithColumns(withdrawListColumns)

		currency, _ := cmd.Flags().GetString("currency")
		state, _ := cmd.Flags().GetString("state")
		count, _ := cmd.Flags().GetInt("count")
		page, _ := cmd.Flags().GetInt("page")

		currency = strings.ToUpper(currency)
		state = strings.ToUpper(state)

		withdrawals, err := wc.ListWithdrawals(cmd.Context(), currency, state, count, page)
		if err != nil {
			return err
		}
		if emptyMessage(withdrawals, i18n.T(i18n.MsgWithdrawListEmpty)) {
			return nil
		}
		return formatter.Format(withdrawals)
	},
}

func init() {
	f := withdrawListCmd.Flags()
	f.String("currency", "", i18n.T(i18n.FlagCurrencyUsage))
	f.String("state", "", i18n.T(i18n.FlagWithdrawStateUsage))
	f.IntP("count", "c", 20, i18n.T(i18n.FlagCountUsage))
	f.Int("page", 1, i18n.T(i18n.FlagPageUsage))
	withdrawCmd.AddCommand(withdrawListCmd)
}
