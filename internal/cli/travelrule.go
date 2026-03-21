package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/wallet"
	"github.com/kyungw00k/upbit/internal/output"
)

// travelruleCmd 트래블룰 부모 명령
var travelruleCmd = &cobra.Command{
	Use:     "travelrule",
	Short:   "트래블룰 검증 관리",
	GroupID: "wallet",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var vaspColumns = []output.TableColumn{
	{Header: "UUID", Key: "vasp_uuid"},
	{Header: "이름", Key: "vasp_name"},
	{Header: "이름(영문)", Key: "en_name"},
	{Header: "입금가능", Key: "depositable"},
	{Header: "출금가능", Key: "withdrawable"},
}

var travelruleVerifyColumns = []output.TableColumn{
	{Header: "입금UUID", Key: "deposit_uuid"},
	{Header: "검증결과", Key: "verification_result"},
	{Header: "입금상태", Key: "deposit_state"},
}

var travelruleVaspsCmd = &cobra.Command{
	Use:   "vasps",
	Short: "트래블룰 지원 거래소 목록 조회",
	Example: `  upbit travelrule vasps              # VASP 목록 조회
  upbit travelrule vasps -o json      # JSON 출력`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		wc := wallet.NewWalletClient(client)
		formatter := GetFormatterWithColumns(vaspColumns)

		vasps, err := wc.GetTravelRuleVASPs(cmd.Context())
		if err != nil {
			return err
		}

		return formatter.Format(vasps)
	},
}

var travelruleVerifyTxIDCmd = &cobra.Command{
	Use:   "verify-txid <txid>",
	Short: "TxID 기반 트래블룰 검증 요청",
	Args:  RequireArgs(1, "검증할 입금의 TxID를 지정하세요"),
	Example: `  upbit travelrule verify-txid abc123 --vasp 8d4fe968-... --currency ETH --net-type ETH`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		wc := wallet.NewWalletClient(client)
		formatter := GetFormatterWithColumns(travelruleVerifyColumns)

		vaspUUID, _ := cmd.Flags().GetString("vasp")
		currency, _ := cmd.Flags().GetString("currency")
		netType, _ := cmd.Flags().GetString("net-type")

		req := &wallet.TravelRuleVerifyByTxIDRequest{
			TxID:     args[0],
			VaspUUID: vaspUUID,
			Currency: currency,
			NetType:  netType,
		}

		result, err := wc.VerifyTravelRuleByTxID(cmd.Context(), req)
		if err != nil {
			return err
		}

		return formatter.Format(result)
	},
}

var travelruleVerifyUUIDCmd = &cobra.Command{
	Use:   "verify-uuid <deposit-uuid>",
	Short: "UUID 기반 트래블룰 검증 요청",
	Args:  RequireArgs(1, "검증할 입금의 UUID를 지정하세요"),
	Example: `  upbit travelrule verify-uuid 5b871d34-... --vasp 8d4fe968-...`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		wc := wallet.NewWalletClient(client)
		formatter := GetFormatterWithColumns(travelruleVerifyColumns)

		vaspUUID, _ := cmd.Flags().GetString("vasp")

		req := &wallet.TravelRuleVerifyByUUIDRequest{
			DepositUUID: args[0],
			VaspUUID:    vaspUUID,
		}

		result, err := wc.VerifyTravelRuleByUUID(cmd.Context(), req)
		if err != nil {
			return err
		}

		return formatter.Format(result)
	},
}

func init() {
	// verify-txid 플래그
	ftxid := travelruleVerifyTxIDCmd.Flags()
	ftxid.String("vasp", "", "상대 거래소 UUID (필수)")
	ftxid.String("currency", "", "통화 코드 (예: ETH) (필수)")
	ftxid.String("net-type", "", "네트워크 식별자 (예: ETH) (필수)")
	_ = travelruleVerifyTxIDCmd.MarkFlagRequired("vasp")
	_ = travelruleVerifyTxIDCmd.MarkFlagRequired("currency")
	_ = travelruleVerifyTxIDCmd.MarkFlagRequired("net-type")

	// verify-uuid 플래그
	fuuid := travelruleVerifyUUIDCmd.Flags()
	fuuid.String("vasp", "", "상대 거래소 UUID (필수)")
	_ = travelruleVerifyUUIDCmd.MarkFlagRequired("vasp")

	// 서브커맨드 등록
	travelruleCmd.AddCommand(travelruleVaspsCmd)
	travelruleCmd.AddCommand(travelruleVerifyTxIDCmd)
	travelruleCmd.AddCommand(travelruleVerifyUUIDCmd)

	rootCmd.AddCommand(travelruleCmd)
}
