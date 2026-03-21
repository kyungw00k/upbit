package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/quotation"
	"github.com/kyungw00k/upbit/internal/output"
)

var tickSizeColumns = []output.TableColumn{
	{Header: "마켓", Key: "market"},
	{Header: "호가통화", Key: "quote_currency"},
	{Header: "호가단위", Key: "tick_size"},
}

var tickSizeCmd = &cobra.Command{
	Use:     "tick-size <market...>",
	Short:   "호가 단위 조회",
	GroupID: "quotation",
	Args:    RequireMinArgs(1, "마켓 코드를 지정하세요 (예: KRW-BTC)"),
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
