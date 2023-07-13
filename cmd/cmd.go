package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "opentelemetry-test-server",
	Long:  "",
}

func init() {
	runCmd.AddCommand(httpServerCmd())
}

func Execute() {
	if err := runCmd.Execute(); err != nil {
		fmt.Printf("cmd err: %s\n", err)
		os.Exit(1)
	}
}
