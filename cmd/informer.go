package cmd

import (
	"github.com/practice/opentelemetry-practice/pkg/common"
	"github.com/practice/opentelemetry-practice/pkg/k8s_resource_otel"
	"github.com/spf13/cobra"
)

func informerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "k8sInformer",
		Short: "run k8s resource informer server",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := &common.ServerConfig{
				Debug:          debug,
				Port:           serverPort,
				JaegerEndpoint: jaegerEndpoint,
			}
			k8s_resource_otel.K8sResourceInformer(cfg)
		},
	}
	return cmd
}
