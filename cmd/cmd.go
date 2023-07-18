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

var (
	debug          bool
	serverPort     string
	jaegerEndpoint string
)

func init() {
	runCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug mode")
	runCmd.PersistentFlags().StringVarP(&serverPort, "port", "p", "8080", "server port")
	runCmd.PersistentFlags().StringVarP(&jaegerEndpoint, "jaegerEndpoint", "j", "http://localhost:14268/api/traces", "jaeger endpoint for trace")
	runCmd.AddCommand(httpServerCmd())
}

func Execute() {
	if err := runCmd.Execute(); err != nil {
		fmt.Printf("cmd err: %s\n", err)
		os.Exit(1)
	}
}
