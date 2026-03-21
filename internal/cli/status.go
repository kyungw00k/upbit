package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/wallet"
	"github.com/kyungw00k/upbit/internal/output"
)

var statusColumns = []output.TableColumn{
	{Header: "통화", Key: "currency"},
	{Header: "지갑 상태", Key: "wallet_state"},
	{Header: "블록 상태", Key: "block_state"},
	{Header: "네트워크", Key: "network_name"},
	{Header: "네트워크 타입", Key: "net_type"},
}

var walletCmd = &cobra.Command{
	Use:     "wallet",
	Short:   "입출금 서비스 상태 조회",
	GroupID: "wallet",
	Example: `  upbit wallet                # 전체 입출금 서비스 상태
  upbit wallet -o json        # JSON 출력`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		wc := wallet.NewWalletClient(client)
		formatter := GetFormatterWithColumns(statusColumns)

		statuses, err := wc.GetServiceStatus(cmd.Context())
		if err != nil {
			return err
		}
		return formatter.Format(statuses)
	},
}

func init() {
	rootCmd.AddCommand(walletCmd)
}
