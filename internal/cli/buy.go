package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/exchange"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var buyCmd = &cobra.Command{
	Use:     "buy <market>",
	Short:   i18n.T(i18n.MsgBuyShort),
	GroupID: "trading",
	Args:    RequireArgs(1, i18n.T(i18n.ErrBuyArgsRequired)),
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
			return fmt.Errorf("%s", i18n.T(i18n.ErrBuyParamsRequired))
		}

		// 지정가 주문 시 호가 단위 자동 보정
		var wasAdjusted bool
		var originalPrice string
		if ordType == "limit" {
			originalPrice = orderPrice
			adjustedPrice, adjusted, adjErr := adjustPrice(cmd.Context(), client, market, orderPrice, "bid")
			if adjErr != nil {
				fmt.Fprint(os.Stderr, i18n.Tf(i18n.MsgTickCheckFailed, adjErr))
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
			msg = i18n.Tf(i18n.MsgBuyOrderAdjusted, market, orderPrice, originalPrice, orderPrice, orderVolume, ordType)
		} else {
			msg = i18n.Tf(i18n.MsgBuyOrderNormal, market, describeBuyOrder(ordType, orderPrice, orderVolume), ordType)
		}

		// --force여도 보정 사실은 stderr에 출력
		if wasAdjusted && GetForce() {
			fmt.Fprint(os.Stderr, i18n.Tf(i18n.MsgTickAdjusted, originalPrice, orderPrice))
		}

		confirmed, err := output.Confirm(msg, GetForce())
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgOrderCancelled))
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
		return i18n.Tf(i18n.MsgDescLimitOrder, price, volume)
	case "price":
		return i18n.Tf(i18n.MsgDescPriceOrder, price)
	default:
		return ""
	}
}

func init() {
	f := buyCmd.Flags()
	f.StringP("price", "p", "", i18n.T(i18n.FlagPriceUsage))
	f.StringP("volume", "V", "", i18n.T(i18n.FlagVolumeUsage))
	f.StringP("total", "t", "", i18n.T(i18n.FlagTotalUsage))
	f.String("tif", "", "Time in Force (ioc, fok, post_only)")
	f.String("smp", "", i18n.T(i18n.FlagSMPUsage))
	f.String("id", "", i18n.T(i18n.FlagIdentifierUsage))
	f.Bool("test", false, i18n.T(i18n.FlagTestUsage))
	AddForceFlag(buyCmd)
	rootCmd.AddCommand(buyCmd)
}
