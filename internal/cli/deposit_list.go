package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/api/wallet"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var depositListColumns = []output.TableColumn{
	{Header: "UUID", Key: "uuid"},
	{Header: i18n.T(i18n.HdrCurrency), Key: "currency"},
	{Header: i18n.T(i18n.HdrAmount), Key: "amount", Format: "number"},
	{Header: i18n.T(i18n.HdrState), Key: "state"},
	{Header: i18n.T(i18n.HdrTransType), Key: "transaction_type"},
	{Header: i18n.T(i18n.HdrCreatedAt), Key: "created_at", Format: "datetime"},
}

var depositListCmd = &cobra.Command{
	Use:   "list",
	Short: i18n.T(i18n.MsgDepositListShort),
	Example: `  upbit deposit list                        # 최근 입금 목록
  upbit deposit list --currency BTC         # BTC 입금만
  upbit deposit list --state ACCEPTED       # 완료된 입금만
  upbit deposit list -c 50                  # 50개 조회
  upbit deposit list --page 2               # 2페이지`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		wc := wallet.NewWalletClient(client)
		formatter := GetFormatterWithColumns(depositListColumns)

		currency, _ := cmd.Flags().GetString("currency")
		state, _ := cmd.Flags().GetString("state")
		count, _ := cmd.Flags().GetInt("count")
		page, _ := cmd.Flags().GetInt("page")

		currency = strings.ToUpper(currency)
		state = strings.ToUpper(state)

		deposits, err := wc.ListDeposits(cmd.Context(), currency, state, count, page)
		if err != nil {
			return err
		}
		if emptyMessage(deposits, i18n.T(i18n.MsgDepositListEmpty)) {
			return nil
		}
		return formatter.Format(deposits)
	},
}

func init() {
	f := depositListCmd.Flags()
	f.String("currency", "", i18n.T(i18n.FlagCurrencyUsage))
	f.String("state", "", i18n.T(i18n.FlagStateUsage))
	f.IntP("count", "c", 100, i18n.T(i18n.FlagCountUsage))
	f.Int("page", 1, i18n.T(i18n.FlagPageUsage))
	depositCmd.AddCommand(depositListCmd)
}
