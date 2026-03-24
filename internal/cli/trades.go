package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/api/quotation"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var tradesColumns = []output.TableColumn{
	{Header: i18n.T(i18n.HdrTradeTime), Key: "trade_time_utc", Format: "time"},
	{Header: i18n.T(i18n.HdrTradePrice), Key: "trade_price", Format: "number"},
	{Header: i18n.T(i18n.HdrTradeVolume), Key: "trade_volume", Format: "number"},
	{Header: i18n.T(i18n.HdrAskBid), Key: "ask_bid"},
	{Header: i18n.T(i18n.HdrChangePrice), Key: "change_price", Format: "number"},
}

var tradesCmd = &cobra.Command{
	Use:     "trades [market]",
	Short:   i18n.T(i18n.MsgTradesShort),
	GroupID: "quotation",
	Args:    RequireArgs(1, i18n.T(i18n.ErrCandleMarketRequired)),
	Example: `  upbit trades KRW-BTC                  # 최근 20건
  upbit trades KRW-BTC -c 100          # 최근 100건
  upbit trades KRW-BTC -o json         # JSON 출력`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		qc := quotation.NewQuotationClient(client)
		formatter := GetFormatterWithColumns(tradesColumns)

		count, _ := cmd.Flags().GetInt("count")

		trades, err := qc.GetRecentTrades(cmd.Context(), args[0], count)
		if err != nil {
			return err
		}

		return formatter.Format(trades)
	},
}

func init() {
	tradesCmd.Flags().IntP("count", "c", 20, i18n.T(i18n.FlagCountUsage))
	rootCmd.AddCommand(tradesCmd)
}
