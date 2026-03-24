package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/api/wallet"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var depositAddressColumns = []output.TableColumn{
	{Header: i18n.T(i18n.HdrCurrency), Key: "currency"},
	{Header: i18n.T(i18n.HdrNetwork), Key: "net_type"},
	{Header: i18n.T(i18n.HdrDepositAddr), Key: "deposit_address"},
	{Header: i18n.T(i18n.HdrTag), Key: "secondary_address"},
}

var depositAddressCmd = &cobra.Command{
	Use:   "address [currency]",
	Short: i18n.T(i18n.MsgDepositAddressShort),
	Args:  cobra.MaximumNArgs(1),
	Example: `  upbit deposit address             # 전체 입금 주소 목록
  upbit deposit address BTC         # BTC 입금 주소 조회 (net_type 자동: BTC)
  upbit deposit address BTC --net-type BTC   # net_type 명시`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		wc := wallet.NewWalletClient(client)
		formatter := GetFormatterWithColumns(depositAddressColumns)

		if len(args) == 0 {
			// 전체 입금 주소 목록
			addresses, err := wc.ListDepositAddresses(cmd.Context())
			if err != nil {
				return err
			}
			if emptyMessage(addresses, i18n.T(i18n.MsgDepositAddressEmpty)) {
				return nil
			}
			return formatter.Format(addresses)
		}

		// 개별 입금 주소 조회
		currency := strings.ToUpper(args[0])
		netType, _ := cmd.Flags().GetString("net-type")
		if netType == "" {
			netType = currency // 기본값: currency와 동일
		}
		netType = strings.ToUpper(netType)

		address, err := wc.GetDepositAddress(cmd.Context(), currency, netType)
		if err != nil {
			return err
		}
		return formatter.Format(address)
	},
}

var depositAddressCreateCmd = &cobra.Command{
	Use:   "create <currency>",
	Short: i18n.T(i18n.MsgDepositCreateShort),
	Args:  RequireArgs(1, i18n.T(i18n.ErrDepositCurrRequired)),
	Example: `  upbit deposit address create BTC
  upbit deposit address create ETH --net-type ETH`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		wc := wallet.NewWalletClient(client)

		currency := strings.ToUpper(args[0])
		netType, _ := cmd.Flags().GetString("net-type")
		if netType == "" {
			netType = currency
		}
		netType = strings.ToUpper(netType)

		_, err = wc.CreateDepositAddress(cmd.Context(), currency, netType)
		if err != nil {
			return err
		}

		fmt.Fprint(os.Stderr, i18n.T(i18n.MsgDepositAddrCreated))
		fmt.Fprint(os.Stderr, i18n.Tf(i18n.MsgDepositAddrCheck, currency, netType))
		return nil
	},
}

func init() {
	depositAddressCmd.Flags().String("net-type", "", i18n.T(i18n.FlagNetTypeUsage))
	depositAddressCreateCmd.Flags().String("net-type", "", i18n.T(i18n.FlagNetTypeUsage))

	depositAddressCmd.AddCommand(depositAddressCreateCmd)
	depositCmd.AddCommand(depositAddressCmd)
}
