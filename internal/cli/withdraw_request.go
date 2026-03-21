package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/wallet"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var withdrawResultColumns = []output.TableColumn{
	{Header: "UUID", Key: "uuid"},
	{Header: i18n.T(i18n.HdrCurrency), Key: "currency"},
	{Header: i18n.T(i18n.HdrAmount), Key: "amount", Format: "number"},
	{Header: i18n.T(i18n.HdrFee), Key: "fee", Format: "number"},
	{Header: i18n.T(i18n.HdrState), Key: "state"},
	{Header: i18n.T(i18n.HdrCreatedAt), Key: "created_at", Format: "datetime"},
}

var withdrawRequestCmd = &cobra.Command{
	Use:   "request <currency>",
	Short: i18n.T(i18n.MsgWithdrawRequestShort),
	Args:  RequireArgs(1, i18n.T(i18n.ErrWithdrawCurrRequired)),
	Example: `  upbit withdraw request BTC --amount 0.5 --to 1A1zP1...       # BTC 출금
  upbit withdraw request XRP --amount 100 --to rN9qN... --secondary-address 3057887915
  upbit withdraw request KRW --amount 50000                    # KRW 원화 출금
  upbit withdraw request BTC --amount 0.01 --to addr --net-type BTC --tx-type internal  # 바로출금
  upbit withdraw request KRW --amount 50000 --force            # 확인 없이 출금`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		wc := wallet.NewWalletClient(client)
		formatter := GetFormatterWithColumns(withdrawResultColumns)

		currency := strings.ToUpper(args[0])
		amount, _ := cmd.Flags().GetString("amount")
		address, _ := cmd.Flags().GetString("to")
		secondaryAddr, _ := cmd.Flags().GetString("secondary-address")
		netType, _ := cmd.Flags().GetString("net-type")
		txType, _ := cmd.Flags().GetString("tx-type")
		twoFactorType, _ := cmd.Flags().GetString("two-factor")

		if amount == "" {
			return fmt.Errorf("%s", i18n.T(i18n.ErrAmountRequired))
		}

		// KRW 자동 감지: 원화 출금
		if currency == "KRW" {
			if twoFactorType == "" {
				return fmt.Errorf("%s", i18n.T(i18n.ErrTwoFactorRequired))
			}

			msg := i18n.Tf(i18n.MsgWithdrawConfirmKRW, amount)
			confirmed, err := output.Confirm(msg, GetForce())
			if err != nil {
				return err
			}
			if !confirmed {
				fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgWithdrawCancelled))
				return nil
			}

			withdrawal, err := wc.WithdrawKRW(cmd.Context(), amount, twoFactorType)
			if err != nil {
				return err
			}
			return formatter.Format(withdrawal)
		}

		// 디지털 자산 출금
		if address == "" {
			return fmt.Errorf("%s", i18n.T(i18n.ErrWithdrawAddrRequired))
		}
		if netType == "" {
			netType = currency // 기본값: currency와 동일
		}
		netType = strings.ToUpper(netType)

		// 확인 프롬프트
		msg := i18n.Tf(i18n.MsgWithdrawConfirmCoin, currency, amount, address)
		if secondaryAddr != "" {
			msg += fmt.Sprintf(" (tag: %s)", secondaryAddr)
		}
		confirmed, err := output.Confirm(msg, GetForce())
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgWithdrawCancelled))
			return nil
		}

		req := &wallet.WithdrawCoinRequest{
			Currency:         currency,
			NetType:          netType,
			Amount:           amount,
			Address:          address,
			SecondaryAddress: secondaryAddr,
			TransactionType:  txType,
		}

		withdrawal, err := wc.WithdrawCoin(cmd.Context(), req)
		if err != nil {
			return err
		}
		return formatter.Format(withdrawal)
	},
}

func init() {
	f := withdrawRequestCmd.Flags()
	f.String("amount", "", i18n.T(i18n.FlagAmountUsage))
	f.String("to", "", i18n.T(i18n.FlagToUsage))
	f.String("secondary-address", "", i18n.T(i18n.FlagSecondaryAddrUsage))
	f.String("net-type", "", i18n.T(i18n.FlagNetTypeUsage))
	f.String("tx-type", "", i18n.T(i18n.FlagTxTypeUsage))
	f.String("two-factor", "", i18n.T(i18n.FlagTwoFactorUsage))
	AddForceFlag(withdrawRequestCmd)
	withdrawCmd.AddCommand(withdrawRequestCmd)
}
