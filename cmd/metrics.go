package cmd

import (
	"github.com/Pantani/healthcheck/internal/collector"
	"github.com/spf13/cobra"
)

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Run scheduled checks from the fixture file",
	Args:  cobra.NoArgs,
	RunE:  provideMetricsCollector,
}

func provideMetricsCollector(cmd *cobra.Command, _ []string) error {
	cfg, err := buildConfig(cmd)
	if err != nil {
		return err
	}
	ctx, stop := commandContext(cmd)
	defer stop()
	return collector.MetricsCollector(ctx, cfg)
}

func init() {
	rootCmd.AddCommand(metricsCmd)
}
