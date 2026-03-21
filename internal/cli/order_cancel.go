package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/exchange"
	"github.com/kyungw00k/upbit/internal/output"
)

var orderCancelCmd = &cobra.Command{
	Use:   "cancel [uuid]",
	Short: "주문 취소",
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
			msg := "모든 대기 주문을 취소합니다"
			if market != "" {
				msg = fmt.Sprintf("%s 마켓의 모든 대기 주문을 취소합니다", market)
			}
			confirmed, err := output.Confirm(msg, GetForce())
			if err != nil {
				return err
			}
			if !confirmed {
				fmt.Fprintln(os.Stderr, "취소가 중단되었습니다")
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
			fmt.Fprintf(os.Stdout, "취소 성공  %d건\n", result.Success.Count)
			for _, o := range result.Success.Orders {
				fmt.Fprintf(os.Stdout, "  %s  %s\n", o.UUID, o.Market)
			}
			if result.Failed.Count > 0 {
				fmt.Fprintf(os.Stdout, "취소 실패  %d건\n", result.Failed.Count)
				for _, o := range result.Failed.Orders {
					fmt.Fprintf(os.Stdout, "  %s  %s\n", o.UUID, o.Market)
				}
			}
			return nil
		}

		// 단건 취소
		if len(args) == 0 {
			return fmt.Errorf("취소할 주문의 UUID를 지정하거나 --all 플래그를 사용하세요")
		}

		uuid := args[0]
		msg := fmt.Sprintf("주문 %s을(를) 취소합니다", uuid)
		confirmed, err := output.Confirm(msg, GetForce())
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Fprintln(os.Stderr, "취소가 중단되었습니다")
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
	f.Bool("all", false, "모든 대기 주문 일괄 취소")
	f.StringP("market", "m", "", "마켓 필터 (--all과 함께 사용)")
	AddForceFlag(orderCancelCmd)
	orderCmd.AddCommand(orderCancelCmd)
}
