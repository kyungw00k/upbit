package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/exchange"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var orderReplaceCmd = &cobra.Command{
	Use:   "replace <uuid>",
	Short: i18n.T(i18n.MsgOrderReplaceShort),
	Args:  RequireArgs(1, i18n.T(i18n.ErrReplaceArgs)),
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
			return fmt.Errorf("%s", i18n.T(i18n.ErrReplaceNoParam))
		}

		// limit 타입이면 price 필수
		if ordType == "limit" && price == "" {
			return fmt.Errorf("%s", i18n.T(i18n.ErrReplaceLimitNoPrice))
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
		msg := i18n.Tf(i18n.MsgReplaceConfirm, uuid)
		if price != "" {
			msg += i18n.Tf(i18n.MsgReplacePrice, price)
		}
		if volume != "" {
			msg += i18n.Tf(i18n.MsgReplaceVolume, volume)
		} else {
			msg += i18n.T(i18n.MsgReplaceRemain)
		}

		confirmed, err := output.Confirm(msg, GetForce())
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgReplaceCancelled))
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
	f.StringP("price", "p", "", i18n.T(i18n.FlagNewPriceUsage))
	f.StringP("volume", "V", "", i18n.T(i18n.FlagNewVolumeUsage))
	f.String("ord-type", "limit", i18n.T(i18n.FlagOrdTypeUsage))
	AddForceFlag(orderReplaceCmd)
	orderCmd.AddCommand(orderReplaceCmd)
}
