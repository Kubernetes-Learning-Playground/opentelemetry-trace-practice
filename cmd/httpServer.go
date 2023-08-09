package cmd

import (
	"github.com/practice/opentelemetry-practice/pkg/common"
	"github.com/practice/opentelemetry-practice/pkg/server"
	"github.com/spf13/cobra"
)

func httpServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "httpServer",
		Short: "run http server",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := &common.ServerConfig{
				Debug:          debug,
				Port:           serverPort,
				JaegerEndpoint: jaegerEndpoint,
			}
			// 启动http server
			server.HttpServer(cfg)
		},
	}
	return cmd
}
