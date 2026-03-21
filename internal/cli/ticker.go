package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/quotation"
	"github.com/kyungw00k/upbit/internal/output"
)

var tickerColumns = []output.TableColumn{
	{Header: "마켓", Key: "market"},
	{Header: "현재가", Key: "trade_price", Format: "number"},
	{Header: "전일 대비", Key: "signed_change_price", Format: "number"},
	{Header: "변동률", Key: "signed_change_rate", Format: "percent"},
	{Header: "거래량(24h)", Key: "acc_trade_volume_24h", Format: "number"},
	{Header: "거래대금(24h)", Key: "acc_trade_price_24h", Format: "number"},
	{Header: "고가", Key: "high_price", Format: "number"},
	{Header: "저가", Key: "low_price", Format: "number"},
}

var tickerCmd = &cobra.Command{
	Use:        "ticker [market...]",
	Short:      "현재가 조회",
	SuggestFor: []string{"tick", "price"},
	GroupID:    "quotation",
	Args:       cobra.ArbitraryArgs,
	Example: `  upbit ticker KRW-BTC                  # 비트코인 현재가
  upbit ticker KRW-BTC KRW-ETH         # 복수 마켓
  upbit ticker --quote KRW              # KRW 마켓 전체 시세
  upbit ticker -q KRW,BTC              # KRW+BTC 마켓 전체 시세
  upbit ticker KRW-BTC -o json         # JSON 출력
  upbit ticker KRW-BTC --json trade_price,signed_change_rate`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		qc := quotation.NewQuotationClient(client)
		formatter := GetFormatterWithColumns(tickerColumns)

		quote, _ := cmd.Flags().GetString("quote")

		// --quote 플래그가 있으면 마켓 단위 전체 현재가 조회
		if quote != "" {
			currencies := strings.Split(quote, ",")
			tickers, err := qc.GetAllTickers(cmd.Context(), currencies)
			if err != nil {
				return err
			}
			return formatter.Format(tickers)
		}

		// 인자가 없으면 에러
		if len(args) == 0 {
			return fmt.Errorf("마켓 코드를 지정하세요 (예: KRW-BTC) 또는 --quote 플래그를 사용하세요\n\n사용법: %s", cmd.CommandPath())
		}

		tickers, err := qc.GetTickers(cmd.Context(), args)
		if err != nil {
			return err
		}

		return formatter.Format(tickers)
	},
}

func init() {
	tickerCmd.Flags().StringP("quote", "q", "", "마켓 통화 코드로 전체 시세 조회 (예: KRW, BTC, USDT)")
	rootCmd.AddCommand(tickerCmd)
}
