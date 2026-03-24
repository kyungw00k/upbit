package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/api/quotation"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var orderbookLevelsColumns = []output.TableColumn{
	{Header: i18n.T(i18n.HdrMarket), Key: "market"},
	{Header: i18n.T(i18n.HdrSupportedLevels), Key: "supported_levels"},
}

var orderbookLevelsCmd = &cobra.Command{
	Use:     "orderbook-levels <market...>",
	Short:   i18n.T(i18n.MsgOrderbookLevelsShort),
	GroupID: "quotation",
	Args:    RequireMinArgs(1, i18n.T(i18n.ErrOrderbookMarket)),
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
