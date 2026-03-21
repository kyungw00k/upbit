package cli

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api/exchange"
	"github.com/kyungw00k/upbit/internal/i18n"
	"github.com/kyungw00k/upbit/internal/output"
	"github.com/kyungw00k/upbit/internal/types"
)

var apiKeysColumns = []output.TableColumn{
	{Header: i18n.T(i18n.HdrAccessKey), Key: "access_key"},
	{Header: i18n.T(i18n.HdrExpireAt), Key: "expire_at", Format: "datetime"},
}

var apiKeysCmd = &cobra.Command{
	Use:        "api-keys",
	Short:      i18n.T(i18n.MsgApiKeysShort),
	SuggestFor: []string{"apikey", "keys"},
	GroupID:    "util",
	Example: `  upbit api-keys              # 만료되지 않은 API 키만
  upbit api-keys --all        # 만료된 키 포함 전체
  upbit api-keys -o json      # JSON 출력`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := GetClientE(true)
		if err != nil {
			return err
		}
		ec := exchange.NewExchangeClient(client)
		formatter := GetFormatterWithColumns(apiKeysColumns)

		keys, err := ec.GetAPIKeys(cmd.Context())
		if err != nil {
			return err
		}

		showAll, _ := cmd.Flags().GetBool("all")
		now := time.Now()

		var filtered []types.APIKey
		for _, k := range keys {
			// 만료 여부 확인
			if !showAll && k.ExpireAt != "" {
				expireTime, err := time.Parse(time.RFC3339, k.ExpireAt)
				if err == nil && expireTime.Before(now) {
					continue // 만료된 키 스킵
				}
			}
			filtered = append(filtered, types.APIKey{
				AccessKey: maskAccessKey(k.AccessKey),
				ExpireAt:  k.ExpireAt,
			})
		}

		if emptyMessage(filtered, i18n.T(i18n.MsgApiKeysEmpty)) {
			return nil
		}
		return formatter.Format(filtered)
	},
}

// maskAccessKey access key 마스킹 (앞 4자리 + ****)
func maskAccessKey(key string) string {
	if len(key) <= 4 {
		return "****"
	}
	return key[:4] + "****"
}

func init() {
	apiKeysCmd.Flags().Bool("all", false, i18n.T(i18n.FlagApiKeysAllUsage))
	rootCmd.AddCommand(apiKeysCmd)
}
