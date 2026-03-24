package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/api/exchange"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var orderCancelCmd = &cobra.Command{
	Use:   "cancel [uuid]",
	Short: i18n.T(i18n.MsgOrderCancelShort),
	Args:  cobra.MaximumNArgs(1),
	Example: `  upbit order cancel 12345678-abcd-efgh-ijkl-1234567890ab   # 단건 취소
  upbit order cancel --all                                  # 전체 대기 주문 취소
  upbit order cancel --all -m KRW-BTC                       # 특정 마켓 전체 취소
  upbit order cancel --all --force                          # 확인 없이 전체 취소`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		ec := exchange.NewExchangeClient(client)

		all, _ := cmd.Flags().GetBool("all")
		market, _ := cmd.Flags().GetString("market")

		if all {
			// 전체 대기 주문 취소
			msg := i18n.T(i18n.MsgCancelAllOrders)
			if market != "" {
				msg = i18n.Tf(i18n.MsgCancelAllMarket, market)
			}
			confirmed, err := output.Confirm(msg, GetForce())
			if err != nil {
				return err
			}
			if !confirmed {
				fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgCancelAborted))
				return nil
			}

			result, err := ec.BatchCancelOrders(cmd.Context(), "all", market, "")
			if err != nil {
				return err
			}

			// JSON/CSV 등은 전체 데이터 출력
			formatter := GetFormatter()
			if _, ok := formatter.(*output.TableFormatter); !ok {
				return formatter.Format(result)
			}

			// 테이블: 결과 요약
			fmt.Fprint(os.Stdout, i18n.Tf(i18n.MsgCancelSuccess, result.Success.Count))
			for _, o := range result.Success.Orders {
				fmt.Fprintf(os.Stdout, "  %s  %s\n", o.UUID, o.Market)
			}
			if result.Failed.Count > 0 {
				fmt.Fprint(os.Stdout, i18n.Tf(i18n.MsgCancelFailed, result.Failed.Count))
				for _, o := range result.Failed.Orders {
					fmt.Fprintf(os.Stdout, "  %s  %s\n", o.UUID, o.Market)
				}
			}
			return nil
		}

		// 단건 취소
		if len(args) == 0 {
			return fmt.Errorf("%s", i18n.T(i18n.ErrCancelNoUUID))
		}

		uuid := args[0]
		msg := i18n.Tf(i18n.MsgCancelSingleOrder, uuid)
		confirmed, err := output.Confirm(msg, GetForce())
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgCancelAborted))
			return nil
		}

		formatter := GetFormatterWithColumns(orderShowColumns)
		order, err := ec.CancelOrder(cmd.Context(), uuid)
		if err != nil {
			return err
		}
		return formatter.Format(order)
	},
}

func init() {
	f := orderCancelCmd.Flags()
	f.Bool("all", false, i18n.T(i18n.FlagCancelAllUsage))
	f.StringP("market", "m", "", i18n.T(i18n.FlagCancelMarketUsage))
	AddForceFlag(orderCancelCmd)
	orderCmd.AddCommand(orderCancelCmd)
}
