package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/exchange"
	"github.com/kyungw00k/upbit/internal/output"
)

var orderListColumns = []output.TableColumn{
	{Header: "UUID", Key: "uuid"},
	{Header: "마켓", Key: "market"},
	{Header: "방향", Key: "side"},
	{Header: "유형", Key: "ord_type"},
	{Header: "가격", Key: "price", Format: "number"},
	{Header: "수량", Key: "volume", Format: "number"},
	{Header: "체결량", Key: "executed_volume", Format: "number"},
	{Header: "상태", Key: "state"},
	{Header: "생성시각", Key: "created_at", Format: "datetime"},
}

var orderListCmd = &cobra.Command{
	Use:   "list",
	Short: "주문 목록 조회",
	Example: `  upbit order list                    # 체결 대기 주문
  upbit order list --closed           # 종료 주문
  upbit order list -m KRW-BTC         # 특정 마켓만
  upbit order list -c 50              # 50개 조회
  upbit order list --closed --page 2  # 종료 주문 2페이지`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		ec := exchange.NewExchangeClient(client)
		formatter := GetFormatterWithColumns(orderListColumns)

		closed, _ := cmd.Flags().GetBool("closed")
		market, _ := cmd.Flags().GetString("market")
		count, _ := cmd.Flags().GetInt("count")
		page, _ := cmd.Flags().GetInt("page")

		if closed {
			orders, err := ec.ListClosedOrders(cmd.Context(), market, count, page)
			if err != nil {
				return err
			}
			if emptyMessage(orders, "종료된 주문이 없습니다") {
				return nil
			}
			return formatter.Format(orders)
		}

		orders, err := ec.ListOpenOrders(cmd.Context(), market, count, page)
		if err != nil {
			return err
		}
		if emptyMessage(orders, "대기 중인 주문이 없습니다") {
			return nil
		}
		return formatter.Format(orders)
	},
}

func init() {
	f := orderListCmd.Flags()
	f.Bool("closed", false, "종료 주문 조회")
	f.StringP("market", "m", "", "마켓 필터 (예: KRW-BTC)")
	f.IntP("count", "c", 20, "조회 개수")
	f.Int("page", 1, "페이지 번호")
	orderCmd.AddCommand(orderListCmd)
}
