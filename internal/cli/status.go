package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/api/wallet"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var statusColumns = []output.TableColumn{
	{Header: i18n.T(i18n.HdrCurrency), Key: "currency"},
	{Header: i18n.T(i18n.HdrWalletState), Key: "wallet_state"},
	{Header: i18n.T(i18n.HdrBlockState), Key: "block_state"},
	{Header: i18n.T(i18n.HdrNetwork), Key: "network_name"},
	{Header: i18n.T(i18n.HdrNetworkType), Key: "net_type"},
}

var walletCmd = &cobra.Command{
	Use:     "wallet",
	Short:   i18n.T(i18n.MsgWalletShort),
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
