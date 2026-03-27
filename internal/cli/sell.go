package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/api/exchange"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var sellCmd = &cobra.Command{
	Use:     "sell <market>",
	Short:   i18n.T(i18n.MsgSellShort),
	GroupID: "trading",
	Args:    RequireArgs(1, i18n.T(i18n.ErrSellArgsRequired)),
	Example: `  upbit sell KRW-BTC -p 50000000 -V 0.001   # 지정가 매도
  upbit sell KRW-BTC -V 0.001               # 시장가 매도 (수량 지정)
  upbit sell KRW-BTC -V 100%                # 전량 시장가 매도
  upbit sell KRW-BTC -p 55000000 -V 50%    # 보유량의 50% 지정가 매도
  upbit sell KRW-BTC -V 0.001 --best       # 최유리 지정가 매도
  upbit sell KRW-BTC -V 0.001 --best --tif fok  # 최유리 지정가 FOK 매도
  upbit sell KRW-BTC --watch 55000000 -p 54500000 -V 0.001  # 예약-지정가 매도
  upbit sell KRW-BTC -p 50000000 -V 0.001 --test  # 테스트 주문
  upbit sell KRW-BTC -V 0.001 --force       # 확인 프롬프트 스킵`,
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
		tif, _ := cmd.Flags().GetString("tif")
		smp, _ := cmd.Flags().GetString("smp")
		identifier, _ := cmd.Flags().GetString("id")
		test, _ := cmd.Flags().GetBool("test")
		best, _ := cmd.Flags().GetBool("best")
		watchPrice, _ := cmd.Flags().GetString("watch")

		// 퍼센트 해석: -V 50% → 보유 코인 기준 실제 수량으로 변환
		if isPercent(volume) {
			origVolume := volume
			_, volume, _, err = resolvePercentOrder(cmd.Context(), client, market, "ask", price, volume, "")
			if err != nil {
				return err
			}
			pct := strings.TrimSuffix(origVolume, "%")
			fmt.Fprint(os.Stderr, i18n.Tf(i18n.MsgPercentResolved, i18n.T(i18n.FlagVolumeUsage), pct, volume))
		}

		// 주문 유형 자동 판별
		var ordType, orderPrice, orderVolume string
		switch {
		case watchPrice != "" && price != "" && volume != "":
			// 예약-지정가 매도: --watch + --price + --volume
			ordType = "limit"
			orderPrice = price
			orderVolume = volume
		case best && volume != "":
			// 최유리 지정가 매도: --best + --volume
			ordType = "best"
			orderVolume = volume
		case price != "" && volume != "":
			// 지정가 매도: --price + --volume
			ordType = "limit"
			orderPrice = price
			orderVolume = volume
		case volume != "" && price == "":
			// 시장가 매도: --volume만
			ordType = "market"
			orderVolume = volume
		default:
			return fmt.Errorf("%s", i18n.T(i18n.ErrSellParamsRequired))
		}

		// 지정가 주문 시 호가 단위 자동 보정
		var wasAdjusted bool
		var originalPrice string
		if ordType == "limit" {
			originalPrice = orderPrice
			adjustedPrice, adjusted, adjErr := adjustPrice(cmd.Context(), client, market, orderPrice, "ask")
			if adjErr != nil {
				fmt.Fprint(os.Stderr, i18n.Tf(i18n.MsgTickCheckFailed, adjErr))
			} else {
				wasAdjusted = adjusted
				orderPrice = adjustedPrice
			}
			// 예약 주문의 감시가도 호가 보정
			if watchPrice != "" {
				adjWatch, _, adjErr := adjustPrice(cmd.Context(), client, market, watchPrice, "ask")
				if adjErr == nil {
					watchPrice = adjWatch
				}
			}
		}

		req := &exchange.OrderRequest{
			Market:      market,
			Side:        "ask",
			OrdType:     ordType,
			Price:       orderPrice,
			Volume:      orderVolume,
			WatchPrice:  watchPrice,
			TimeInForce: tif,
			SMPType:     smp,
			Identifier:  identifier,
		}

		// 확인 프롬프트
		var msg string
		if wasAdjusted {
			msg = i18n.Tf(i18n.MsgSellOrderAdjusted, market, orderPrice, originalPrice, orderPrice, orderVolume, ordType)
		} else {
			msg = i18n.Tf(i18n.MsgSellOrderNormal, market, describeSellOrder(ordType, orderPrice, orderVolume, watchPrice), ordType)
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

func describeSellOrder(ordType, price, volume, watchPrice string) string {
	if watchPrice != "" {
		return i18n.Tf(i18n.MsgDescReservedOrder, watchPrice, price, volume)
	}
	switch ordType {
	case "limit":
		return i18n.Tf(i18n.MsgDescLimitOrder, price, volume)
	case "market":
		return i18n.Tf(i18n.MsgDescMarketSell, volume)
	case "best":
		return i18n.Tf(i18n.MsgDescBestOrder, volume)
	default:
		return ""
	}
}

func init() {
	f := sellCmd.Flags()
	f.StringP("price", "p", "", i18n.T(i18n.FlagPriceUsage))
	f.StringP("volume", "V", "", i18n.T(i18n.FlagVolumeUsage))
	f.Bool("best", false, i18n.T(i18n.FlagBestUsage))
	f.String("watch", "", i18n.T(i18n.FlagWatchPriceUsage))
	f.String("tif", "", "Time in Force (ioc, fok, post_only)")
	f.String("smp", "", i18n.T(i18n.FlagSMPUsage))
	f.String("id", "", i18n.T(i18n.FlagIdentifierUsage))
	f.Bool("test", false, i18n.T(i18n.FlagTestUsage))
	AddForceFlag(sellCmd)
	rootCmd.AddCommand(sellCmd)
}
