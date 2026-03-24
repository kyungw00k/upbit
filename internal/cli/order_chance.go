package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/api/exchange"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var orderChanceCmd = &cobra.Command{
	Use:   "chance <market>",
	Short: i18n.T(i18n.MsgOrderChanceShort),
	Args:  RequireArgs(1, i18n.T(i18n.ErrCandleMarketRequired)),
	Example: `  upbit order chance KRW-BTC   # KRW-BTC 주문 가능 정보
  upbit order chance KRW-ETH   # KRW-ETH 주문 가능 정보`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		ec := exchange.NewExchangeClient(client)

		market := args[0]
		chance, err := ec.GetOrderChance(cmd.Context(), market)
		if err != nil {
			return err
		}

		// JSON/CSV는 전체 데이터 출력
		formatter := GetFormatter()
		if _, ok := formatter.(*output.TableFormatter); !ok {
			return formatter.Format(chance)
		}

		// 테이블: 핵심 정보만 커스텀 렌더링
		type kv struct{ label, value string }
		rows := []kv{
			{i18n.T(i18n.LblMarket), fmt.Sprintf("%s (%s)", chance.Market.ID, chance.Market.Name)},
			{i18n.T(i18n.LblState), chance.Market.State},
			{i18n.T(i18n.LblBidFee), fmt.Sprintf("%s (Maker: %s)", chance.BidFee, chance.MakerBidFee)},
			{i18n.T(i18n.LblAskFee), fmt.Sprintf("%s (Maker: %s)", chance.AskFee, chance.MakerAskFee)},
			{i18n.T(i18n.LblBidAvailable), fmt.Sprintf("%s %s", chance.BidAccount.Balance, chance.BidAccount.Currency)},
			{i18n.T(i18n.LblAskAvailable), fmt.Sprintf("%s %s", chance.AskAccount.Balance, chance.AskAccount.Currency)},
			{i18n.T(i18n.LblMinOrder), fmt.Sprintf("%s %s", chance.Market.Bid.MinTotal, chance.BidAccount.Currency)},
			{i18n.T(i18n.LblBidTypes), strings.Join(chance.Market.BidTypes, ", ")},
			{i18n.T(i18n.LblAskTypes), strings.Join(chance.Market.AskTypes, ", ")},
		}

		maxLabelWidth := 0
		for _, r := range rows {
			if w := runewidth.StringWidth(r.label); w > maxLabelWidth {
				maxLabelWidth = w
			}
		}

		padRight := func(s string, width int) string {
			sw := runewidth.StringWidth(s)
			if sw >= width {
				return s
			}
			return s + strings.Repeat(" ", width-sw)
		}

		w := os.Stdout
		for _, r := range rows {
			fmt.Fprintf(w, "%s  %s\n", padRight(r.label, maxLabelWidth), r.value)
		}
		return nil
	},
}

func init() {
	orderCmd.AddCommand(orderChanceCmd)
}
