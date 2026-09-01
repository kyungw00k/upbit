package cli

import (
	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/i18n"
)

var pocketCmd = &cobra.Command{
	Use:   "pocket",
	Short: i18n.T(i18n.MsgPocketRootShort),
}

func init() {
	rootCmd.AddCommand(pocketCmd)
}
