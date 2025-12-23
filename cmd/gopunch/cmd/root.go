package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version   = "2.0.0"
	BuildDate = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "gopunch",
	Short: "⚡ Dead Simple Uptime Monitoring Tool",
	Long: `GoPunch is a lightweight CLI utility for checking uptime and response times.

Features:
  • Multi-URL support with customizable HTTP methods
  • Interval-based continuous monitoring
  • Real-time TUI dashboard
  • Alerting via webhooks (Discord, Slack)
  • Export to CSV/JSON formats
  • Concurrent requests with configurable parallelism`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var cfgFile string

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is gopunch.json)")
}
