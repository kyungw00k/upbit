package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/api/pocket"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
)

var pocketAPIKeyColumns = []output.TableColumn{
	{Header: i18n.T(i18n.HdrAccessKey), Key: "access_key"},
	{Header: "Permissions", Key: "permissions"},
	{Header: "Expired At", Key: "expired_at"},
}

var pocketAPIKeysCmd = &cobra.Command{
	Use:   "apikeys [uuid...]",
	Short: i18n.T(i18n.MsgPocketAPIKeysShort),
	Example: `  upbit pocket apikeys                    # 전체 포켓의 API Key
  upbit pocket apikeys --include-expired  # 만료된 Key 포함`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		pc := pocket.NewPocketClient(client)

		includeExpired, _ := cmd.Flags().GetBool("include-expired")
		groups, err := pc.ListPocketAPIKeys(cmd.Context(), args, includeExpired)
		if err != nil {
			return err
		}

		// 포켓별 그룹을 (pocket uuid, key) 평면 리스트로 펼친다
		rows := make([]pocketAPIKeyRow, 0, len(groups)*2)
		for _, g := range groups {
			for _, k := range g.Keys {
				rows = append(rows, pocketAPIKeyRow{
					PocketUUID:  g.UUID,
					AccessKey:   k.AccessKey,
					Permissions: k.Permissions,
					ExpiredAt:   k.ExpiredAt,
				})
			}
		}
		formatter := GetFormatterWithColumns(append([]output.TableColumn{
			{Header: "Pocket UUID", Key: "pocket_uuid"},
		}, pocketAPIKeyColumns...))
		return formatter.Format(rows)
	},
}

type pocketAPIKeyRow struct {
	PocketUUID  string   `json:"pocket_uuid"`
	AccessKey   string   `json:"access_key"`
	Permissions []string `json:"permissions"`
	ExpiredAt   string   `json:"expired_at"`
}

func init() {
	pocketAPIKeysCmd.Flags().Bool("include-expired", false, i18n.T(i18n.FlagPocketIncExpired))
	pocketCmd.AddCommand(pocketAPIKeysCmd)
}
