package main

import (
	"fmt"
	"os"

	"github.com/kyungw00k/upbit/internal/api"
	"github.com/kyungw00k/upbit/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		// APIError인 경우 종료 코드 사용
		if apiErr, ok := api.AsAPIError(err); ok {
			fmt.Fprintln(os.Stderr, apiErr.Error())
			os.Exit(apiErr.ExitCode())
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
