package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/wallet"
	"github.com/kyungw00k/upbit/internal/output"
)

var depositListColumns = []output.TableColumn{
	{Header: "UUID", Key: "uuid"},
	{Header: "통화", Key: "currency"},
	{Header: "금액", Key: "amount", Format: "number"},
	{Header: "상태", Key: "state"},
	{Header: "유형", Key: "transaction_type"},
	{Header: "생성일", Key: "created_at", Format: "datetime"},
}

var depositListCmd = &cobra.Command{
	Use:   "list",
	Short: "입금 목록 조회",
	Example: `  upbit deposit list                        # 최근 입금 목록
  upbit deposit list --currency BTC         # BTC 입금만
  upbit deposit list --state ACCEPTED       # 완료된 입금만
  upbit deposit list -c 50                  # 50개 조회
  upbit deposit list --page 2               # 2페이지`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		wc := wallet.NewWalletClient(client)
		formatter := GetFormatterWithColumns(depositListColumns)

		currency, _ := cmd.Flags().GetString("currency")
		state, _ := cmd.Flags().GetString("state")
		count, _ := cmd.Flags().GetInt("count")
		page, _ := cmd.Flags().GetInt("page")

		currency = strings.ToUpper(currency)
		state = strings.ToUpper(state)

		deposits, err := wc.ListDeposits(cmd.Context(), currency, state, count, page)
		if err != nil {
			return err
		}
		if emptyMessage(deposits, "입금 내역이 없습니다") {
			return nil
		}
		return formatter.Format(deposits)
	},
}

func init() {
	f := depositListCmd.Flags()
	f.String("currency", "", "통화 코드 필터 (예: BTC, KRW)")
	f.String("state", "", "입금 상태 필터 (PROCESSING, ACCEPTED, CANCELLED, REJECTED)")
	f.IntP("count", "c", 100, "조회 개수")
	f.Int("page", 1, "페이지 번호")
	depositCmd.AddCommand(depositListCmd)
}
