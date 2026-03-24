package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/api/exchange"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var orderShowColumns = []output.TableColumn{
	{Header: "UUID", Key: "uuid"},
	{Header: i18n.T(i18n.HdrMarket), Key: "market"},
	{Header: i18n.T(i18n.HdrSide), Key: "side"},
	{Header: i18n.T(i18n.HdrOrdType), Key: "ord_type"},
	{Header: i18n.T(i18n.HdrOrderPrice), Key: "price", Format: "number"},
	{Header: i18n.T(i18n.HdrOrderVolume), Key: "volume", Format: "number"},
	{Header: i18n.T(i18n.HdrRemainingVol), Key: "remaining_volume", Format: "number"},
	{Header: i18n.T(i18n.HdrExecutedVol), Key: "executed_volume", Format: "number"},
	{Header: i18n.T(i18n.HdrState), Key: "state"},
	{Header: i18n.T(i18n.HdrCreatedAt), Key: "created_at", Format: "datetime"},
}

var orderShowCmd = &cobra.Command{
	Use:   "show <uuid...>",
	Short: i18n.T(i18n.MsgOrderShowShort),
	Args:  RequireMinArgs(1, i18n.T(i18n.ErrOrderShowArgs)),
	Example: `  upbit order show 12345678-abcd-efgh-ijkl-1234567890ab            # 단일 조회
  upbit order show uuid1 uuid2 uuid3                                # 복수 UUID 조회
  upbit order show my-order-001 --id                                # Identifier로 조회`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		ec := exchange.NewExchangeClient(client)

		byID, _ := cmd.Flags().GetBool("id")

		// Identifier 조회는 단일만 지원
		if byID {
			formatter := GetFormatterWithColumns(orderShowColumns)
			order, err := ec.GetOrderByIdentifier(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return formatter.Format(order)
		}

		// UUID 1개: 기존 단건 상세 조회
		if len(args) == 1 {
			formatter := GetFormatterWithColumns(orderShowColumns)
			order, err := ec.GetOrder(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return formatter.Format(order)
		}

		// UUID 2개 이상: 복수 주문 조회
		formatter := GetFormatterWithColumns(orderShowColumns)
		orders, err := ec.GetOrdersByUUIDs(cmd.Context(), args)
		if err != nil {
			return err
		}
		return formatter.Format(orders)
	},
}

func init() {
	orderShowCmd.Flags().Bool("id", false, i18n.T(i18n.FlagOrderShowIDUsage))
	orderCmd.AddCommand(orderShowCmd)
}
