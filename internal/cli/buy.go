package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/exchange"
	"github.com/kyungw00k/upbit/internal/output"
)

var buyCmd = &cobra.Command{
	Use:     "buy <market>",
	Short:   "매수 주문",
	GroupID: "trading",
	Args:    RequireArgs(1, "마켓 코드를 지정하세요 (예: KRW-BTC)"),
	Example: `  upbit buy KRW-BTC -p 50000000 -V 0.001   # 지정가 매수
  upbit buy KRW-BTC -t 100000              # 시장가 매수 (총액 지정)
  upbit buy KRW-BTC -p 50000000 -V 0.001 --test  # 테스트 주문
  upbit buy KRW-BTC -t 100000 --force      # 확인 프롬프트 스킵`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		ec := exchange.NewExchangeClient(client)
		formatter := GetFormatterWithColumns(orderShowColumns)

		market := args[0]
		price, _ := cmd.Flags().GetString("price")
		volume, _ := cmd.Flags().GetString("volume")
		total, _ := cmd.Flags().GetString("total")
		tif, _ := cmd.Flags().GetString("tif")
		smp, _ := cmd.Flags().GetString("smp")
		identifier, _ := cmd.Flags().GetString("id")
		test, _ := cmd.Flags().GetBool("test")

		// 주문 유형 자동 판별
		var ordType, orderPrice, orderVolume string
		switch {
		case price != "" && volume != "":
			// 지정가 매수: --price + --volume
			ordType = "limit"
			orderPrice = price
			orderVolume = volume
		case total != "":
			// 시장가 매수: --total (총액)
			ordType = "price"
			orderPrice = total
		default:
			return fmt.Errorf("매수 주문에는 --price와 --volume (지정가) 또는 --total (시장가)이 필요합니다")
		}

		// 지정가 주문 시 호가 단위 자동 보정
		var wasAdjusted bool
		var originalPrice string
		if ordType == "limit" {
			originalPrice = orderPrice
			adjustedPrice, adjusted, adjErr := adjustPrice(cmd.Context(), client, market, orderPrice, "bid")
			if adjErr != nil {
				fmt.Fprintf(os.Stderr, "호가 단위 확인 실패: %v\n", adjErr)
			} else {
				wasAdjusted = adjusted
				orderPrice = adjustedPrice
			}
		}

		req := &exchange.OrderRequest{
			Market:      market,
			Side:        "bid",
			OrdType:     ordType,
			Price:       orderPrice,
			Volume:      orderVolume,
			TimeInForce: tif,
			SMPType:     smp,
			Identifier:  identifier,
		}

		// 확인 프롬프트
		var msg string
		if wasAdjusted {
			msg = fmt.Sprintf("매수 주문: %s 단가=%s (호가 보정: %s→%s), 수량=%s (유형: %s)", market, orderPrice, originalPrice, orderPrice, orderVolume, ordType)
		} else {
			msg = fmt.Sprintf("매수 주문: %s %s (유형: %s)", market, describeBuyOrder(ordType, orderPrice, orderVolume), ordType)
		}

		// --force여도 보정 사실은 stderr에 출력
		if wasAdjusted && GetForce() {
			fmt.Fprintf(os.Stderr, "호가 보정: %s → %s\n", originalPrice, orderPrice)
		}

		confirmed, err := output.Confirm(msg, GetForce())
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Fprintln(os.Stderr, "주문이 취소되었습니다")
			return nil
		}

		if test {
			order, err := ec.TestOrder(cmd.Context(), req)
			if err != nil {
				return err
			}
			return formatter.Format(order)
		}

		order, err := ec.CreateOrder(cmd.Context(), req)
		if err != nil {
			return err
		}
		return formatter.Format(order)
	},
}

func describeBuyOrder(ordType, price, volume string) string {
	switch ordType {
	case "limit":
		return fmt.Sprintf("단가=%s, 수량=%s", price, volume)
	case "price":
		return fmt.Sprintf("총액=%s", price)
	default:
		return ""
	}
}

func init() {
	f := buyCmd.Flags()
	f.StringP("price", "p", "", "주문 단가")
	f.StringP("volume", "V", "", "주문 수량")
	f.StringP("total", "t", "", "주문 총액 (시장가 매수)")
	f.String("tif", "", "Time in Force (ioc, fok, post_only)")
	f.String("smp", "", "자기 거래 방지 (cancel_maker, cancel_taker, reduce)")
	f.String("id", "", "클라이언트 지정 주문 식별자")
	f.Bool("test", false, "테스트 주문 (실제 체결 안됨)")
	AddForceFlag(buyCmd)
	rootCmd.AddCommand(buyCmd)
}
