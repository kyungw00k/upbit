package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/quotation"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var tickerColumns = []output.TableColumn{
	{Header: i18n.T(i18n.HdrMarket), Key: "market"},
	{Header: i18n.T(i18n.HdrPrice), Key: "trade_price", Format: "number"},
	{Header: i18n.T(i18n.HdrChange), Key: "signed_change_price", Format: "number"},
	{Header: i18n.T(i18n.HdrChangeRate), Key: "signed_change_rate", Format: "percent"},
	{Header: i18n.T(i18n.HdrVolume24h), Key: "acc_trade_volume_24h", Format: "number"},
	{Header: i18n.T(i18n.HdrTradePrice24h), Key: "acc_trade_price_24h", Format: "number"},
	{Header: i18n.T(i18n.HdrHigh), Key: "high_price", Format: "number"},
	{Header: i18n.T(i18n.HdrLow), Key: "low_price", Format: "number"},
}

var tickerCmd = &cobra.Command{
	Use:        "ticker [market...]",
	Short:      i18n.T(i18n.MsgTickerShort),
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
			return fmt.Errorf("%s", i18n.Tf(i18n.ErrTickerNoMarket, cmd.CommandPath()))
		}

		tickers, err := qc.GetTickers(cmd.Context(), args)
		if err != nil {
			return err
		}

		return formatter.Format(tickers)
	},
}

func init() {
	tickerCmd.Flags().StringP("quote", "q", "", i18n.T(i18n.FlagTickerQuoteUsage))
	rootCmd.AddCommand(tickerCmd)
}
