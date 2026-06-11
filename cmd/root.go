package cmd

import (
	"context"
	"fmt"
	"github.com/Pantani/healthcheck/internal/config"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	fixturesPath string
	redisURL     string
	httpTimeout  time.Duration
)

var (
	rootCmd = &cobra.Command{
		Use:          "healthcheck",
		Short:        "Run scheduled HTTP health checks and create PagerDuty incidents on failures",
		SilenceUsage: true,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	defaults := config.Default()
	rootCmd.PersistentFlags().StringVar(&fixturesPath, "fixtures", defaults.Fixtures.Path, "path to the JSON fixture file")
	rootCmd.PersistentFlags().StringVar(&redisURL, "redis-url", defaults.Redis.URL, "Redis URL used to store the previous check value")
	rootCmd.PersistentFlags().DurationVar(&httpTimeout, "http-timeout", defaults.HTTP.Timeout, "HTTP request timeout")
}

func buildConfig(cmd *cobra.Command) (config.Configuration, error) {
	cfg, err := config.Load()
	if err != nil {
		return config.Configuration{}, err
	}
	if flagChanged(cmd, "fixtures") {
		cfg.Fixtures.Path = fixturesPath
	}
	if flagChanged(cmd, "redis-url") {
		cfg.Redis.URL = redisURL
	}
	if flagChanged(cmd, "http-timeout") {
		cfg.HTTP.Timeout = httpTimeout
	}
	config.Apply(cfg)
	return cfg, nil
}

func commandContext(cmd *cobra.Command) (context.Context, func()) {
	ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
	return ctx, cancel
}

func flagChanged(cmd *cobra.Command, name string) bool {
	if flag := cmd.Flags().Lookup(name); flag != nil {
		return flag.Changed
	}
	if flag := cmd.InheritedFlags().Lookup(name); flag != nil {
		return flag.Changed
	}
	return false
}
