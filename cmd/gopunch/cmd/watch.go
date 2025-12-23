package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/TheRemyyy/gopunch/internal/alerter"
	"github.com/TheRemyyy/gopunch/internal/checker"
)

var (
	watchInterval    int
	watchTimeout     int
	watchMethod      string
	watchHeaders     []string
	watchData        string
	watchInsecure    bool
	watchFollowRedir bool
	watchExpect      []int
	watchConcurrency int
	watchQuiet       bool
)

var watchCmd = &cobra.Command{
	Use:   "watch [url...]",
	Short: "Continuously monitor URLs at regular intervals",
	Long: `Watch one or more URLs and display their status in real-time.
Press Ctrl+C to stop and see summary statistics.

If a config file is present (or specified via --config), URLs and settings 
from the config will be used. Command line arguments override config values.`,
	Run: runWatch,
}

func init() {
	rootCmd.AddCommand(watchCmd)

	watchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 5, "Check interval in seconds")
	watchCmd.Flags().IntVarP(&watchTimeout, "timeout", "t", 10, "Request timeout in seconds")
	watchCmd.Flags().StringVarP(&watchMethod, "method", "m", "GET", "HTTP method")
	watchCmd.Flags().StringArrayVarP(&watchHeaders, "header", "H", nil, "Custom headers")
	watchCmd.Flags().StringVarP(&watchData, "data", "d", "", "Request body")
	watchCmd.Flags().BoolVarP(&watchInsecure, "insecure", "k", false, "Skip TLS verification")
	watchCmd.Flags().BoolVarP(&watchFollowRedir, "follow", "L", true, "Follow redirects")
	watchCmd.Flags().IntSliceVarP(&watchExpect, "expect", "e", nil, "Expected status codes")
	watchCmd.Flags().IntVarP(&watchConcurrency, "concurrency", "c", 10, "Max concurrent requests")
	watchCmd.Flags().BoolVarP(&watchQuiet, "quiet", "q", false, "Minimal output")
}

type WatchStats struct {
	URL           string
	Checks        int
	Successes     int
	Failures      int
	TotalTime     time.Duration
	MinTime       time.Duration
	MaxTime       time.Duration
	ResponseTimes []time.Duration
	LastSuccess   bool
}

func runWatch(cmd *cobra.Command, args []string) {
	// Load config if specified or exists
	var cfg *Config
	if cfgFile != "" {
		var err error
		cfg, err = LoadConfig(cfgFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Try default gopunch.json
		if _, err := os.Stat("gopunch.json"); err == nil {
			loaded, err := LoadConfig("gopunch.json")
			if err == nil {
				cfg = loaded
			}
		}
	}

	// Merge config with flags
	urls := args
	if cfg != nil {
		if len(urls) == 0 {
			urls = cfg.URLs
		}

		// Apply defaults from config if flags not changed
		if !cmd.Flags().Changed("interval") && cfg.Interval > 0 {
			watchInterval = cfg.Interval
		}
		if !cmd.Flags().Changed("timeout") && cfg.Timeout > 0 {
			watchTimeout = cfg.Timeout
		}
		if !cmd.Flags().Changed("method") && cfg.Method != "" {
			watchMethod = cfg.Method
		}
		if !cmd.Flags().Changed("header") && cfg.Headers != nil {
			for k, v := range cfg.Headers {
				watchHeaders = append(watchHeaders, fmt.Sprintf("%s:%s", k, v))
			}
		}
		if !cmd.Flags().Changed("insecure") {
			watchInsecure = cfg.Insecure
		}
		if !cmd.Flags().Changed("follow") {
			watchFollowRedir = cfg.Follow
		}
		if !cmd.Flags().Changed("concurrency") && cfg.Concurrency > 0 {
			watchConcurrency = cfg.Concurrency
		}
		if !cmd.Flags().Changed("expect") && len(cfg.ExpectedCodes) > 0 {
			watchExpect = cfg.ExpectedCodes
		}
	}

	if len(urls) == 0 {
		cmd.Help()
		os.Exit(1)
	}

	opts := checker.Options{
		Timeout:         time.Duration(watchTimeout) * time.Second,
		Method:          strings.ToUpper(watchMethod),
		Headers:         parseHeaders(watchHeaders),
		Body:            watchData,
		Insecure:        watchInsecure,
		FollowRedirects: watchFollowRedir,
		ExpectedCodes:   watchExpect,
		Retries:         1,
	}

	if cfg != nil && cfg.Retries > 0 && !cmd.Flags().Changed("retries") {
		opts.Retries = cfg.Retries
	}

	// Setup Alerter
	var alertSystem *alerter.Alerter
	if cfg != nil && cfg.Alerting != nil && cfg.Alerting.Enabled {
		wc := cfg.Alerting.Webhook
		webhookConfig := &alerter.WebhookConfig{URL: wc.URL, Method: wc.Method}

		alertSystem = alerter.New(alerter.Config{
			Enabled:  cfg.Alerting.Enabled,
			Cooldown: time.Duration(cfg.Alerting.Cooldown) * time.Second,
			Webhook:  webhookConfig,
		})
		fmt.Println("🔔 Alerting enabled")
	}

	green := color.New(color.FgGreen, color.Bold)
	red := color.New(color.FgRed, color.Bold)
	cyan := color.New(color.FgCyan)
	yellow := color.New(color.FgYellow)

	stats := make(map[string]*WatchStats)
	for _, url := range urls {
		stats[url] = &WatchStats{
			URL:         url,
			MinTime:     time.Hour,
			LastSuccess: true, // Assume healthy start to avoid immediate alert? Or false?
			// Better true to avoid alerting on first run unless it fails
		}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(time.Duration(watchInterval) * time.Second)
	defer ticker.Stop()

	cyan.Printf("\n⚡ Watching %d URL(s) every %ds (Ctrl+C to stop)\n\n", len(urls), watchInterval)

	// Initial check
	runWatchCycle(urls, opts, stats, green, red, yellow, watchQuiet, watchConcurrency, alertSystem)

	for {
		select {
		case <-sigChan:
			fmt.Println()
			printWatchSummary(stats)
			return
		case <-ticker.C:
			runWatchCycle(urls, opts, stats, green, red, yellow, watchQuiet, watchConcurrency, alertSystem)
		}
	}
}

func runWatchCycle(urls []string, opts checker.Options, stats map[string]*WatchStats,
	green, red, yellow *color.Color, quiet bool, concurrency int, alert *alerter.Alerter) {

	results := checker.CheckURLs(urls, opts, concurrency)
	timestamp := time.Now().Format("15:04:05")

	for _, r := range results {
		s := stats[r.URL]
		s.Checks++
		s.TotalTime += r.Duration
		s.ResponseTimes = append(s.ResponseTimes, r.Duration)

		if r.Duration < s.MinTime {
			s.MinTime = r.Duration
		}
		if r.Duration > s.MaxTime {
			s.MaxTime = r.Duration
		}

		if r.Success && r.Error == nil {
			s.Successes++

			// Recovery alert
			if !s.LastSuccess && alert != nil {
				go alert.SendRecoveryAlert(r.URL)
			}
			s.LastSuccess = true

			if !quiet {
				green.Printf("[%s] ✓ %s - %d %dms\n", timestamp, r.URL, r.StatusCode, r.Duration.Milliseconds())
			}
		} else {
			s.Failures++
			errMsg := ""
			if r.Error != nil {
				errMsg = r.Error.Error()
			} else {
				errMsg = fmt.Sprintf("status %d", r.StatusCode)
			}

			// Failure alert
			if s.LastSuccess && alert != nil {
				go alert.SendAlert(alerter.Alert{
					URL:       r.URL,
					Status:    fmt.Sprintf("%d %s", r.StatusCode, r.Status),
					Error:     errMsg,
					Timestamp: time.Now(),
				})
			}
			// Repeat alerts handled by cooldown in alerter
			if !s.LastSuccess && alert != nil {
				go alert.SendAlert(alerter.Alert{
					URL:       r.URL,
					Status:    fmt.Sprintf("%d %s", r.StatusCode, r.Status),
					Error:     errMsg,
					Timestamp: time.Now(),
				})
			}

			s.LastSuccess = false

			if !quiet {
				red.Printf("[%s] ✗ %s - %s\n", timestamp, r.URL, errMsg)
			}
		}
	}
}

func printWatchSummary(stats map[string]*WatchStats) {
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	cyan.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	cyan.Println("              WATCH SUMMARY              ")
	cyan.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"URL", "Checks", "Success", "Failed", "Uptime", "Avg", "Min", "Max"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, s := range stats {
		uptime := 0.0
		if s.Checks > 0 {
			uptime = float64(s.Successes) / float64(s.Checks) * 100
		}

		avgTime := time.Duration(0)
		if s.Checks > 0 {
			avgTime = s.TotalTime / time.Duration(s.Checks)
		}

		uptimeStr := fmt.Sprintf("%.1f%%", uptime)
		if uptime >= 99 {
			uptimeStr = green.Sprint(uptimeStr)
		} else if uptime < 90 {
			uptimeStr = red.Sprint(uptimeStr)
		}

		table.Append([]string{
			s.URL,
			fmt.Sprintf("%d", s.Checks),
			green.Sprintf("%d", s.Successes),
			red.Sprintf("%d", s.Failures),
			uptimeStr,
			fmt.Sprintf("%dms", avgTime.Milliseconds()),
			fmt.Sprintf("%dms", s.MinTime.Milliseconds()),
			fmt.Sprintf("%dms", s.MaxTime.Milliseconds()),
		})
	}

	fmt.Println()
	table.Render()
	fmt.Println()
}
