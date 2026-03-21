package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/quotation"
	"github.com/kyungw00k/upbit/internal/output"
)

var marketColumns = []output.TableColumn{
	{Header: "마켓", Key: "market"},
	{Header: "한글명", Key: "korean_name"},
	{Header: "영문명", Key: "english_name"},
}

var marketCmd = &cobra.Command{
	Use:     "market",
	Short:   "마켓 목록 조회",
	GroupID: "quotation",
	Example: `  upbit market              # 전체 마켓 목록
  upbit market -q KRW       # KRW 마켓만
  upbit market -q BTC       # BTC 마켓만
  upbit market -q USDT      # USDT 마켓만`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		qc := quotation.NewQuotationClient(client)
		formatter := GetFormatterWithColumns(marketColumns)

		markets, err := qc.GetMarkets(cmd.Context())
		if err != nil {
			return err
		}

		// --quote 필터 적용
		quote, _ := cmd.Flags().GetString("quote")
		if quote != "" {
			prefix := strings.ToUpper(quote) + "-"
			var filtered []interface{}
			for _, m := range markets {
				if strings.HasPrefix(m.Market, prefix) {
					filtered = append(filtered, m)
				}
			}
			if emptyMessage(filtered, "--quote 필터 결과가 없습니다") {
				return nil
			}
			return formatter.Format(filtered)
		}

		return formatter.Format(markets)
	},
}

func init() {
	marketCmd.Flags().StringP("quote", "q", "", "호가 통화 필터 (KRW, BTC, USDT)")
	rootCmd.AddCommand(marketCmd)
}
