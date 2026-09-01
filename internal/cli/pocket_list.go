package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/api/pocket"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var pocketListColumns = []output.TableColumn{
	{Header: "UUID", Key: "uuid"},
	{Header: "Name", Key: "name"},
}

var pocketListCmd = &cobra.Command{
	Use:     "list",
	Short:   i18n.T(i18n.MsgPocketListShort),
	Example: `  upbit pocket list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		pc := pocket.NewPocketClient(client)

		pockets, err := pc.ListPockets(cmd.Context())
		if err != nil {
			return err
		}
		return GetFormatterWithColumns(pocketListColumns).Format(pockets)
	},
}

func init() {
	pocketCmd.AddCommand(pocketListCmd)
}
