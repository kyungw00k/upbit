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

var pocketTransferColumns = []output.TableColumn{
	{Header: "UUID", Key: "uuid"},
	{Header: "From", Key: "from"},
	{Header: "To", Key: "to"},
	{Header: i18n.T(i18n.HdrCurrency), Key: "currency"},
	{Header: i18n.T(i18n.HdrAmount), Key: "amount", Format: "number"},
	{Header: i18n.T(i18n.HdrState), Key: "state"},
	{Header: i18n.T(i18n.HdrCreatedAt), Key: "created_at", Format: "datetime"},
}

var pocketTransferCmd = &cobra.Command{
	Use:   "transfer",
	Short: i18n.T(i18n.MsgPocketTransferShort),
}

var pocketTransferSendCmd = &cobra.Command{
	Use:   "send <to> <currency> <amount>",
	Short: i18n.T(i18n.MsgPocketTransferSend),
	Args:  RequireArgs(3, i18n.T(i18n.ErrPocketTransferArgs)),
	Example: `  upbit pocket transfer send 9ca023a5-851b-4fec-9f0a-48cd83c2eaae XRP 10   # XRP 10개 이전
  upbit pocket transfer send 9ca023a5-... XRP 10 --id my-transfer-01    # 식별자 지정`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		pc := pocket.NewPocketClient(client)

		identifier, _ := cmd.Flags().GetString("id")
		req := &pocket.TransferRequest{
			To:         args[0],
			Currency:   strings.ToUpper(args[1]),
			Amount:     args[2],
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

		transfer, err := pc.Transfer(cmd.Context(), req)
		if err != nil {
			return err
		}
		return GetFormatterWithColumns(pocketTransferColumns).Format(transfer)
	},
}

var pocketTransferListCmd = &cobra.Command{
	Use:   "list",
	Short: i18n.T(i18n.MsgPocketTransferList),
	Example: `  upbit pocket transfer list                       # 최근 이전 내역
  upbit pocket transfer list --currency XRP        # 통화 필터
  upbit pocket transfer list --direction out       # 나간 이전만
  upbit pocket transfer list --states WAIT,OK      # 상태 필터`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		pc := pocket.NewPocketClient(client)

		currency, _ := cmd.Flags().GetString("currency")
		direction, _ := cmd.Flags().GetString("direction")
		statesFlag, _ := cmd.Flags().GetStringSlice("states")
		limit, _ := cmd.Flags().GetInt("limit")

		var states []string
		for _, s := range statesFlag {
			s = strings.TrimSpace(s)
			if s != "" {
				states = append(states, strings.ToUpper(s))
			}
		}

		transfers, err := pc.ListTransfers(cmd.Context(), pocket.ListTransfersOptions{
			Direction: direction,
			States:    states,
			Currency:  strings.ToUpper(currency),
			Limit:     limit,
		})
		if err != nil {
			return err
		}
		return GetFormatterWithColumns(pocketTransferColumns).Format(transfers)
	},
}

func init() {
	pocketTransferSendCmd.Flags().String("id", "", "transfer identifier")
	AddForceFlag(pocketTransferSendCmd)
	pocketTransferListCmd.Flags().String("currency", "", i18n.T(i18n.FlagCurrencyUsage))
	pocketTransferListCmd.Flags().String("direction", "", i18n.T(i18n.FlagPocketDirection))
	pocketTransferListCmd.Flags().StringSlice("states", nil, i18n.T(i18n.FlagPocketStates))
	pocketTransferListCmd.Flags().Int("limit", 0, "result limit")
	pocketTransferCmd.AddCommand(pocketTransferSendCmd)
	pocketTransferCmd.AddCommand(pocketTransferListCmd)
	pocketCmd.AddCommand(pocketTransferCmd)
}
