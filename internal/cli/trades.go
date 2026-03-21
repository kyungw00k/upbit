package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/quotation"
	"github.com/kyungw00k/upbit/internal/output"
)

var tradesColumns = []output.TableColumn{
	{Header: "체결 시각", Key: "trade_time_utc", Format: "time"},
	{Header: "체결가", Key: "trade_price", Format: "number"},
	{Header: "체결량", Key: "trade_volume", Format: "number"},
	{Header: "매수/매도", Key: "ask_bid"},
	{Header: "전일 대비", Key: "change_price", Format: "number"},
}

var tradesCmd = &cobra.Command{
	Use:     "trades [market]",
	Short:   "체결 내역 조회",
	GroupID: "quotation",
	Args:    RequireArgs(1, "마켓 코드를 지정하세요 (예: KRW-BTC)"),
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
	tradesCmd.Flags().IntP("count", "c", 20, "조회 개수")
	rootCmd.AddCommand(tradesCmd)
}
