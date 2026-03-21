package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/exchange"
	"github.com/kyungw00k/upbit/internal/output"
)

var orderShowColumns = []output.TableColumn{
	{Header: "UUID", Key: "uuid"},
	{Header: "마켓", Key: "market"},
	{Header: "방향", Key: "side"},
	{Header: "유형", Key: "ord_type"},
	{Header: "가격", Key: "price", Format: "number"},
	{Header: "수량", Key: "volume", Format: "number"},
	{Header: "잔여수량", Key: "remaining_volume", Format: "number"},
	{Header: "체결량", Key: "executed_volume", Format: "number"},
	{Header: "상태", Key: "state"},
	{Header: "생성시각", Key: "created_at", Format: "datetime"},
}

var orderShowCmd = &cobra.Command{
	Use:   "show <uuid...>",
	Short: "주문 상세 조회",
	Args:  RequireMinArgs(1, "주문 UUID를 지정하세요"),
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
	orderShowCmd.Flags().Bool("id", false, "Identifier로 조회 (UUID 대신)")
	orderCmd.AddCommand(orderShowCmd)
}
