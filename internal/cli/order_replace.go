package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/exchange"
	"github.com/kyungw00k/upbit/internal/output"
)

var orderReplaceCmd = &cobra.Command{
	Use:   "replace <uuid>",
	Short: "주문 정정 (취소 후 재주문)",
	Args:  RequireArgs(1, "정정할 주문 UUID를 지정하세요"),
	Example: `  upbit order replace 12345678-abcd -p 51000000 -V 0.001   # 가격·수량 변경
  upbit order replace 12345678-abcd -p 51000000            # 가격만 변경
  upbit order replace 12345678-abcd -V 0.002               # 수량만 변경
  upbit order replace 12345678-abcd -p 51000000 --force    # 확인 없이 정정`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		ec := exchange.NewExchangeClient(client)
		formatter := GetFormatterWithColumns(orderShowColumns)

		uuid := args[0]
		price, _ := cmd.Flags().GetString("price")
		volume, _ := cmd.Flags().GetString("volume")
		ordType, _ := cmd.Flags().GetString("ord-type")

		if price == "" && volume == "" {
			return fmt.Errorf("정정할 --price 또는 --volume을 지정하세요")
		}

		// limit 타입이면 price 필수
		if ordType == "limit" && price == "" {
			return fmt.Errorf("limit 주문 정정에는 --price가 필요합니다")
		}

		req := &exchange.CancelAndNewOrderRequest{
			PrevOrderUUID: uuid,
			NewOrdType:    ordType,
		}

		if price != "" {
			req.NewPrice = price
		}
		if volume != "" {
			req.NewVolume = volume
		} else {
			req.NewVolume = "remain_only"
		}

		// 확인 프롬프트
		msg := fmt.Sprintf("주문 %s 정정: ", uuid)
		if price != "" {
			msg += fmt.Sprintf("단가=%s ", price)
		}
		if volume != "" {
			msg += fmt.Sprintf("수량=%s", volume)
		} else {
			msg += "(잔량 유지)"
		}

		confirmed, err := output.Confirm(msg, GetForce())
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Fprintln(os.Stderr, "정정이 취소되었습니다")
			return nil
		}

		result, err := ec.CancelAndNewOrder(cmd.Context(), req)
		if err != nil {
			return err
		}
		return formatter.Format(result)
	},
}

func init() {
	f := orderReplaceCmd.Flags()
	f.StringP("price", "p", "", "신규 주문 단가")
	f.StringP("volume", "V", "", "신규 주문 수량")
	f.String("ord-type", "limit", "주문 유형 (limit, best)")
	AddForceFlag(orderReplaceCmd)
	orderCmd.AddCommand(orderReplaceCmd)
}
