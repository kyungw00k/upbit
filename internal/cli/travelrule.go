package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/api/wallet"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

// travelruleCmd 트래블룰 부모 명령
var travelruleCmd = &cobra.Command{
	Use:     "travelrule",
	Short:   i18n.T(i18n.MsgTravelruleShort),
	GroupID: "wallet",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var vaspColumns = []output.TableColumn{
	{Header: "UUID", Key: "vasp_uuid"},
	{Header: i18n.T(i18n.HdrVaspName), Key: "vasp_name"},
	{Header: i18n.T(i18n.HdrVaspNameEn), Key: "en_name"},
	{Header: i18n.T(i18n.HdrDepositable), Key: "depositable"},
	{Header: i18n.T(i18n.HdrWithdrawable), Key: "withdrawable"},
}

var travelruleVerifyColumns = []output.TableColumn{
	{Header: i18n.T(i18n.HdrDepositUUID), Key: "deposit_uuid"},
	{Header: i18n.T(i18n.HdrVerifyResult), Key: "verification_result"},
	{Header: i18n.T(i18n.HdrDepositState), Key: "deposit_state"},
}

var travelruleVaspsCmd = &cobra.Command{
	Use:   "vasps",
	Short: i18n.T(i18n.MsgTravelruleVaspsShort),
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
	Short: i18n.T(i18n.MsgTravelruleTxIDShort),
	Args:  RequireArgs(1, i18n.T(i18n.ErrTravelruleTxIDArgs)),
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
	Short: i18n.T(i18n.MsgTravelruleUUIDShort),
	Args:  RequireArgs(1, i18n.T(i18n.ErrTravelruleUUIDArgs)),
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
	ftxid.String("vasp", "", i18n.T(i18n.FlagVaspUsage))
	ftxid.String("currency", "", i18n.T(i18n.FlagTRCurrencyUsage))
	ftxid.String("net-type", "", i18n.T(i18n.FlagTRNetTypeUsage))
	_ = travelruleVerifyTxIDCmd.MarkFlagRequired("vasp")
	_ = travelruleVerifyTxIDCmd.MarkFlagRequired("currency")
	_ = travelruleVerifyTxIDCmd.MarkFlagRequired("net-type")

	// verify-uuid 플래그
	fuuid := travelruleVerifyUUIDCmd.Flags()
	fuuid.String("vasp", "", i18n.T(i18n.FlagVaspUsage))
	_ = travelruleVerifyUUIDCmd.MarkFlagRequired("vasp")

	// 서브커맨드 등록
	travelruleCmd.AddCommand(travelruleVaspsCmd)
	travelruleCmd.AddCommand(travelruleVerifyTxIDCmd)
	travelruleCmd.AddCommand(travelruleVerifyUUIDCmd)

	rootCmd.AddCommand(travelruleCmd)
}
