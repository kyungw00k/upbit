package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/exchange"
	"github.com/kyungw00k/upbit/internal/output"
)

var orderChanceCmd = &cobra.Command{
	Use:   "chance <market>",
	Short: "마켓별 주문 가능 정보 조회",
	Args:  RequireArgs(1, "마켓 코드를 지정하세요 (예: KRW-BTC)"),
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
			{"마켓", fmt.Sprintf("%s (%s)", chance.Market.ID, chance.Market.Name)},
			{"상태", chance.Market.State},
			{"매수 수수료", fmt.Sprintf("%s (Maker: %s)", chance.BidFee, chance.MakerBidFee)},
			{"매도 수수료", fmt.Sprintf("%s (Maker: %s)", chance.AskFee, chance.MakerAskFee)},
			{"매수 가능", fmt.Sprintf("%s %s", chance.BidAccount.Balance, chance.BidAccount.Currency)},
			{"매도 가능", fmt.Sprintf("%s %s", chance.AskAccount.Balance, chance.AskAccount.Currency)},
			{"최소 주문", fmt.Sprintf("%s %s", chance.Market.Bid.MinTotal, chance.BidAccount.Currency)},
			{"매수 유형", strings.Join(chance.Market.BidTypes, ", ")},
			{"매도 유형", strings.Join(chance.Market.AskTypes, ", ")},
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
