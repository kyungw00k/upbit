package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/quotation"
	"github.com/kyungw00k/upbit/internal/output"
)

var orderbookLevelsColumns = []output.TableColumn{
	{Header: "마켓", Key: "market"},
	{Header: "지원 단위", Key: "supported_levels"},
}

var orderbookLevelsCmd = &cobra.Command{
	Use:     "orderbook-levels <market...>",
	Short:   "호가 모아보기 단위 조회",
	GroupID: "quotation",
	Args:    RequireMinArgs(1, "마켓 코드를 지정하세요 (예: KRW-BTC)"),
	Example: `  upbit orderbook-levels KRW-BTC KRW-ETH    # 호가 모아보기 단위 조회
  upbit orderbook-levels KRW-BTC            # 단일 마켓`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		qc := quotation.NewQuotationClient(client)
		formatter := GetFormatterWithColumns(orderbookLevelsColumns)

		levels, err := qc.GetOrderbookLevels(cmd.Context(), args)
		if err != nil {
			return err
		}

		return formatter.Format(levels)
	},
}

func init() {
	rootCmd.AddCommand(orderbookLevelsCmd)
}
