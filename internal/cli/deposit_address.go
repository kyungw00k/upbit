package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/wallet"
	"github.com/kyungw00k/upbit/internal/output"
)

var depositAddressColumns = []output.TableColumn{
	{Header: "통화", Key: "currency"},
	{Header: "네트워크", Key: "net_type"},
	{Header: "주소", Key: "deposit_address"},
	{Header: "태그", Key: "secondary_address"},
}

var depositAddressCmd = &cobra.Command{
	Use:   "address [currency]",
	Short: "입금 주소 조회",
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
			if emptyMessage(addresses, "등록된 입금 주소가 없습니다 (upbit deposit address create <currency>로 생성)") {
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
	Short: "입금 주소 생성 요청",
	Args:  RequireArgs(1, "통화 코드를 지정하세요 (예: BTC)"),
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

		fmt.Fprintf(os.Stderr, "주소 생성이 요청되었습니다.\n")
		fmt.Fprintf(os.Stderr, "'upbit deposit address %s --net-type %s' 명령으로 생성 결과를 확인하세요.\n", currency, netType)
		return nil
	},
}

func init() {
	depositAddressCmd.Flags().String("net-type", "", "네트워크 유형 (기본: currency와 동일)")
	depositAddressCreateCmd.Flags().String("net-type", "", "네트워크 유형 (기본: currency와 동일)")

	depositAddressCmd.AddCommand(depositAddressCreateCmd)
	depositCmd.AddCommand(depositAddressCmd)
}
