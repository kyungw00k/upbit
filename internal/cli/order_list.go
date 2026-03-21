package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/exchange"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var orderListColumns = []output.TableColumn{
	{Header: "UUID", Key: "uuid"},
	{Header: i18n.T(i18n.HdrMarket), Key: "market"},
	{Header: i18n.T(i18n.HdrSide), Key: "side"},
	{Header: i18n.T(i18n.HdrOrdType), Key: "ord_type"},
	{Header: i18n.T(i18n.HdrOrderPrice), Key: "price", Format: "number"},
	{Header: i18n.T(i18n.HdrOrderVolume), Key: "volume", Format: "number"},
	{Header: i18n.T(i18n.HdrExecutedVol), Key: "executed_volume", Format: "number"},
	{Header: i18n.T(i18n.HdrState), Key: "state"},
	{Header: i18n.T(i18n.HdrCreatedAt), Key: "created_at", Format: "datetime"},
}

var orderListCmd = &cobra.Command{
	Use:   "list",
	Short: i18n.T(i18n.MsgOrderListShort),
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
			if emptyMessage(orders, i18n.T(i18n.MsgOrderListClosed)) {
				return nil
			}
			return formatter.Format(orders)
		}

		orders, err := ec.ListOpenOrders(cmd.Context(), market, count, page)
		if err != nil {
			return err
		}
		if emptyMessage(orders, i18n.T(i18n.MsgOrderListOpen)) {
			return nil
		}
		return formatter.Format(orders)
	},
}

func init() {
	f := orderListCmd.Flags()
	f.Bool("closed", false, i18n.T(i18n.FlagClosedUsage))
	f.StringP("market", "m", "", i18n.T(i18n.FlagMarketFilterUsage))
	f.IntP("count", "c", 20, i18n.T(i18n.FlagCountUsage))
	f.Int("page", 1, i18n.T(i18n.FlagPageUsage))
	orderCmd.AddCommand(orderListCmd)
}
