package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/api/pocket"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var pocketUniversalCmd = &cobra.Command{
	Use:   "universal",
	Short: i18n.T(i18n.MsgPocketUniversalShort),
}

var pocketUniversalSendCmd = &cobra.Command{
	Use:   "send <from> <to> <currency> <amount>",
	Short: i18n.T(i18n.MsgPocketUniversalSend),
	Args:  RequireArgs(4, i18n.T(i18n.ErrPocketUnivTransferArg)),
	Example: `  upbit pocket universal send 11111111-... 9ca023a5-851b-4fec-9f0a-48cd83c2eaae XRP 10   # 포켓 간 이전
  upbit pocket list                                                                     # UUID 확인
  upbit pocket universal send 9ca023a5-... 11111111-... KRW 100000 --id my-xfer        # 식별자 지정`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		pc := pocket.NewPocketClient(client)

		identifier, _ := cmd.Flags().GetString("id")
		req := &pocket.UniversalTransferRequest{
			From:       args[0],
			To:         args[1],
			Currency:   strings.ToUpper(args[2]),
			Amount:     args[3],
			Identifier: identifier,
		}

		msg := i18n.Tf(i18n.MsgPocketConfirmSend, req.Amount, req.Currency, req.To)
		confirmed, err := output.Confirm(msg, GetForce())
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgCancelAborted))
			return nil
		}

		transfer, err := pc.UniversalTransfer(cmd.Context(), req)
		if err != nil {
			return err
		}
		return GetFormatterWithColumns(pocketTransferColumns).Format(transfer)
	},
}

var pocketUniversalListCmd = &cobra.Command{
	Use:   "list",
	Short: i18n.T(i18n.MsgPocketUniversalList),
	Example: `  upbit pocket universal list                                        # 최근 메인포켓 이전 내역
  upbit pocket universal list --from 9ca023a5-... --to 11111111-...  # 방향 필터`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		pc := pocket.NewPocketClient(client)

		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		currency, _ := cmd.Flags().GetString("currency")
		statesFlag, _ := cmd.Flags().GetStringSlice("states")
		limit, _ := cmd.Flags().GetInt("limit")

		var states []string
		for _, s := range statesFlag {
			s = strings.TrimSpace(s)
			if s != "" {
				states = append(states, strings.ToUpper(s))
			}
		}

		transfers, err := pc.ListUniversalTransfers(cmd.Context(), pocket.ListUniversalTransfersOptions{
			ListTransfersOptions: pocket.ListTransfersOptions{
				States:   states,
				Currency: strings.ToUpper(currency),
				Limit:    limit,
			},
			From: from,
			To:   to,
		})
		if err != nil {
			return err
		}
		return GetFormatterWithColumns(pocketTransferColumns).Format(transfers)
	},
}

func init() {
	pocketUniversalSendCmd.Flags().String("id", "", "transfer identifier")
	AddForceFlag(pocketUniversalSendCmd)
	pocketUniversalListCmd.Flags().String("from", "", "sender pocket UUID or 'main'")
	pocketUniversalListCmd.Flags().String("to", "", "receiver pocket UUID or 'main'")
	pocketUniversalListCmd.Flags().String("currency", "", i18n.T(i18n.FlagCurrencyUsage))
	pocketUniversalListCmd.Flags().StringSlice("states", nil, i18n.T(i18n.FlagPocketStates))
	pocketUniversalListCmd.Flags().Int("limit", 0, "result limit")
	pocketUniversalCmd.AddCommand(pocketUniversalSendCmd)
	pocketUniversalCmd.AddCommand(pocketUniversalListCmd)
	pocketCmd.AddCommand(pocketUniversalCmd)
}
