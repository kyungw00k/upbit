package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/wallet"
	"github.com/kyungw00k/upbit/internal/output"
)

var withdrawListColumns = []output.TableColumn{
	{Header: "UUID", Key: "uuid"},
	{Header: "통화", Key: "currency"},
	{Header: "금액", Key: "amount", Format: "number"},
	{Header: "수수료", Key: "fee", Format: "number"},
	{Header: "상태", Key: "state"},
	{Header: "유형", Key: "transaction_type"},
	{Header: "생성일", Key: "created_at", Format: "datetime"},
}

var withdrawListCmd = &cobra.Command{
	Use:   "list",
	Short: "출금 목록 조회",
	Example: `  upbit withdraw list                        # 최근 출금 목록
  upbit withdraw list --currency BTC         # BTC 출금만
  upbit withdraw list --state DONE           # 완료된 출금만
  upbit withdraw list -c 50                  # 50개 조회
  upbit withdraw list --page 2               # 2페이지`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		wc := wallet.NewWalletClient(client)
		formatter := GetFormatterWithColumns(withdrawListColumns)

		currency, _ := cmd.Flags().GetString("currency")
		state, _ := cmd.Flags().GetString("state")
		count, _ := cmd.Flags().GetInt("count")
		page, _ := cmd.Flags().GetInt("page")

		currency = strings.ToUpper(currency)
		state = strings.ToUpper(state)

		withdrawals, err := wc.ListWithdrawals(cmd.Context(), currency, state, count, page)
		if err != nil {
			return err
		}
		if emptyMessage(withdrawals, "출금 내역이 없습니다") {
			return nil
		}
		return formatter.Format(withdrawals)
	},
}

func init() {
	f := withdrawListCmd.Flags()
	f.String("currency", "", "통화 코드 필터 (예: BTC, KRW)")
	f.String("state", "", "출금 상태 필터 (WAITING, PROCESSING, DONE, FAILED, CANCELLED, REJECTED)")
	f.IntP("count", "c", 20, "조회 개수")
	f.Int("page", 1, "페이지 번호")
	withdrawCmd.AddCommand(withdrawListCmd)
}
