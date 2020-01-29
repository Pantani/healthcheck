package cmd

import (
	"github.com/Pantani/healthcheck/internal/collector"
	"github.com/spf13/cobra"
)

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Collect metrics from fixtures",
	Run:   provideMetricsCollector,
}

func provideMetricsCollector(_ *cobra.Command, args []string) {
	collector.MetricsCollector()
}

func init() {
	rootCmd.AddCommand(metricsCmd)
}
