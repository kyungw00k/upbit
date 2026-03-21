package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/wallet"
	"github.com/kyungw00k/upbit/internal/output"
)

var withdrawResultColumns = []output.TableColumn{
	{Header: "UUID", Key: "uuid"},
	{Header: "통화", Key: "currency"},
	{Header: "금액", Key: "amount", Format: "number"},
	{Header: "수수료", Key: "fee", Format: "number"},
	{Header: "상태", Key: "state"},
	{Header: "생성시각", Key: "created_at", Format: "datetime"},
}

var withdrawRequestCmd = &cobra.Command{
	Use:   "request <currency>",
	Short: "출금 요청",
	Args:  RequireArgs(1, "출금할 통화 코드를 지정하세요 (예: BTC, KRW)"),
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
			return fmt.Errorf("--amount 플래그는 필수입니다")
		}

		// KRW 자동 감지: 원화 출금
		if currency == "KRW" {
			if twoFactorType == "" {
				return fmt.Errorf("KRW 출금 시 --two-factor 플래그가 필요합니다 (kakao, naver, hana)")
			}

			msg := fmt.Sprintf("출금: KRW %s", amount)
			confirmed, err := output.Confirm(msg, GetForce())
			if err != nil {
				return err
			}
			if !confirmed {
				fmt.Fprintln(os.Stderr, "출금이 취소되었습니다")
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
			return fmt.Errorf("디지털 자산 출금에는 --to (수신 주소) 플래그가 필요합니다")
		}
		if netType == "" {
			netType = currency // 기본값: currency와 동일
		}
		netType = strings.ToUpper(netType)

		// 확인 프롬프트
		msg := fmt.Sprintf("출금: %s %s -> %s", currency, amount, address)
		if secondaryAddr != "" {
			msg += fmt.Sprintf(" (tag: %s)", secondaryAddr)
		}
		confirmed, err := output.Confirm(msg, GetForce())
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Fprintln(os.Stderr, "출금이 취소되었습니다")
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
	f.String("amount", "", "출금 수량/금액 (필수)")
	f.String("to", "", "수신 주소 (디지털 자산 필수)")
	f.String("secondary-address", "", "2차 주소 (태그/메모)")
	f.String("net-type", "", "네트워크 유형 (기본: currency와 동일)")
	f.String("tx-type", "", "트랜잭션 유형 (default, internal)")
	f.String("two-factor", "", "2차 인증 수단 (kakao, naver, hana) — KRW 전용")
	AddForceFlag(withdrawRequestCmd)
	withdrawCmd.AddCommand(withdrawRequestCmd)
}
