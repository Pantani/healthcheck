package cmd

import (
	"fmt"
	"github.com/Pantani/healthcheck/internal/config"
	"github.com/spf13/cobra"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"os"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "Metric Collector",
		Short: "Collect metrics and format for Prometheus pull",
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	loadConf := func() { config.InitConfig() }
	loadLogger := func() { logger.InitLogger() }
	cobra.OnInitialize(loadConf)
	cobra.OnInitialize(loadLogger)
}
