package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/api/quotation"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var tickSizeColumns = []output.TableColumn{
	{Header: i18n.T(i18n.HdrMarket), Key: "market"},
	{Header: i18n.T(i18n.HdrQuoteCurrency), Key: "quote_currency"},
	{Header: i18n.T(i18n.HdrTickSize), Key: "tick_size"},
}

var tickSizeCmd = &cobra.Command{
	Use:     "tick-size <market...>",
	Short:   i18n.T(i18n.MsgTickSizeShort),
	GroupID: "quotation",
	Args:    RequireMinArgs(1, i18n.T(i18n.ErrCandleMarketRequired)),
	Example: `  upbit tick-size KRW-BTC KRW-XRP    # 호가 단위 조회
  upbit tick-size KRW-BTC            # 단일 마켓`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		qc := quotation.NewQuotationClient(client)
		formatter := GetFormatterWithColumns(tickSizeColumns)

		tickSizes, err := qc.GetTickSizes(cmd.Context(), args)
		if err != nil {
			return err
		}

		return formatter.Format(tickSizes)
	},
}

func init() {
	rootCmd.AddCommand(tickSizeCmd)
}
