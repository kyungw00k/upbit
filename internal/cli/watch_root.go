package cli

import (
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:     "watch",
	Short:   "실시간 WebSocket 스트림 구독",
	GroupID: "realtime",
	Example: `  upbit watch ticker KRW-BTC KRW-ETH
  upbit watch orderbook KRW-BTC
  upbit watch trade KRW-BTC
  upbit watch candle KRW-BTC -i 1m
  upbit watch my-order
  upbit watch my-asset`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
