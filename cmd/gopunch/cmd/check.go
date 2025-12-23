package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/TheRemyyy/gopunch/internal/checker"
)

var (
	checkTimeout     int
	checkMethod      string
	checkHeaders     []string
	checkData        string
	checkInsecure    bool
	checkFollowRedir bool
	checkExpect      []int
	checkRetries     int
	checkFormat      string
	checkQuiet       bool
	checkConcurrency int
)

var checkCmd = &cobra.Command{
	Use:   "check <url> [url...]",
	Short: "Perform one-time health check",
	Long: `Check one or more URLs, TCP ports, DNS records, or SSL certificates.

Supports schemes:
  http://, https://  - Standard HTTP check
  tcp://host:port    - TCP port check
  dns://host         - DNS resolution check
  ssl://host:port    - SSL certificate expiry check

If config file exists, defaults are loaded from it.`,
	Args: cobra.ArbitraryArgs,
	Run:  runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().IntVarP(&checkTimeout, "timeout", "t", 10, "Request timeout in seconds")
	checkCmd.Flags().StringVarP(&checkMethod, "method", "m", "GET", "HTTP method")
	checkCmd.Flags().StringArrayVarP(&checkHeaders, "header", "H", nil, "Custom headers")
	checkCmd.Flags().StringVarP(&checkData, "data", "d", "", "Request body")
	checkCmd.Flags().BoolVarP(&checkInsecure, "insecure", "k", false, "Skip TLS verification")
	checkCmd.Flags().BoolVarP(&checkFollowRedir, "follow", "L", true, "Follow redirects")
	checkCmd.Flags().IntSliceVarP(&checkExpect, "expect", "e", nil, "Expected status codes")
	checkCmd.Flags().IntVarP(&checkRetries, "retries", "r", 0, "Number of retries on failure")
	checkCmd.Flags().StringVarP(&checkFormat, "format", "f", "table", "Output format")
	checkCmd.Flags().BoolVarP(&checkQuiet, "quiet", "q", false, "Minimal output")
	checkCmd.Flags().IntVarP(&checkConcurrency, "concurrency", "c", 10, "Max concurrent requests")
}

func runCheck(cmd *cobra.Command, args []string) {
	// Load config
	var cfg *Config
	if cfgFile != "" {
		var err error
		cfg, err = LoadConfig(cfgFile)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}
	} else {
		if _, err := os.Stat("gopunch.json"); err == nil {
			loaded, err := LoadConfig("gopunch.json")
			if err == nil {
				cfg = loaded
			}
		}
	}

	// Apply config defaults
	if cfg != nil {
		if len(args) == 0 {
			args = cfg.URLs
		}
		if !cmd.Flags().Changed("timeout") && cfg.Timeout > 0 {
			checkTimeout = cfg.Timeout
		}
		if !cmd.Flags().Changed("method") && cfg.Method != "" {
			checkMethod = cfg.Method
		}
		if !cmd.Flags().Changed("header") && cfg.Headers != nil {
			for k, v := range cfg.Headers {
				checkHeaders = append(checkHeaders, fmt.Sprintf("%s:%s", k, v))
			}
		}
		if !cmd.Flags().Changed("insecure") {
			checkInsecure = cfg.Insecure
		}
		if !cmd.Flags().Changed("follow") {
			checkFollowRedir = cfg.Follow
		}
		if !cmd.Flags().Changed("concurrency") && cfg.Concurrency > 0 {
			checkConcurrency = cfg.Concurrency
		}
		if !cmd.Flags().Changed("expect") && len(cfg.ExpectedCodes) > 0 {
			checkExpect = cfg.ExpectedCodes
		}
		if !cmd.Flags().Changed("retries") && cfg.Retries > 0 {
			checkRetries = cfg.Retries
		}
	}

	opts := checker.Options{
		Timeout:         time.Duration(checkTimeout) * time.Second,
		Method:          strings.ToUpper(checkMethod),
		Headers:         parseHeaders(checkHeaders),
		Body:            checkData,
		Insecure:        checkInsecure,
		FollowRedirects: checkFollowRedir,
		ExpectedCodes:   checkExpect,
		Retries:         checkRetries,
	}

	results := checker.CheckURLs(args, opts, checkConcurrency)

	if checkQuiet {
		for _, r := range results {
			if r.Error != nil || !r.Success {
				os.Exit(1)
			}
		}
		os.Exit(0)
	}

	switch checkFormat {
	case "json":
		printJSON(results)
	case "csv":
		printCSV(results)
	case "minimal":
		printMinimal(results)
	default:
		printTable(results)
	}

	for _, r := range results {
		if r.Error != nil || !r.Success {
			os.Exit(1)
		}
	}
}

func printTable(results []checker.Result) {
	green := color.New(color.FgGreen, color.Bold)
	red := color.New(color.FgRed, color.Bold)
	yellow := color.New(color.FgYellow)
	cyan := color.New(color.FgCyan)

	fmt.Println()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Status", "Target", "Code/Info", "Time", "Note"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	// SetHeaderLine not available in v0.0.5, using default
	table.SetTablePadding("  ")
	table.SetNoWhiteSpace(true)

	for _, r := range results {
		var status, statusColor string
		if r.Error != nil {
			status = "✗"
			statusColor = red.Sprint(status)
		} else if r.Success {
			status = "✓"
			statusColor = green.Sprint(status)
		} else {
			status = "!"
			statusColor = yellow.Sprint(status)
		}

		code := "-"
		note := ""

		if r.StatusCode > 0 {
			code = fmt.Sprintf("%d", r.StatusCode)
		} else if r.Info != "" {
			code = r.Info
		}

		if r.Size > 0 {
			note = formatBytes(r.Size)
		} else if r.Retries > 0 {
			note = fmt.Sprintf("%d retries", r.Retries)
		}

		timeStr := fmt.Sprintf("%dms", r.Duration.Milliseconds())

		table.Append([]string{
			statusColor,
			cyan.Sprint(r.URL),
			code,
			timeStr,
			note,
		})
	}
	table.Render()
	fmt.Println()
}

func printJSON(results []checker.Result) {
	fmt.Println("[")
	for i, r := range results {
		errStr := "null"
		if r.Error != nil {
			errStr = fmt.Sprintf("\"%s\"", r.Error.Error())
		}
		info := r.Info
		if info == "" && r.StatusCode > 0 {
			info = fmt.Sprintf("%d", r.StatusCode)
		}

		fmt.Printf(`  {"url":"%s","info":"%s","duration_ms":%d,"size":%d,"success":%t,"error":%s}`,
			r.URL, info, r.Duration.Milliseconds(), r.Size, r.Success, errStr)
		if i < len(results)-1 {
			fmt.Println(",")
		} else {
			fmt.Println()
		}
	}
	fmt.Println("]")
}

func printCSV(results []checker.Result) {
	fmt.Println("url,info,duration_ms,size,success,error")
	for _, r := range results {
		errStr := ""
		if r.Error != nil {
			errStr = r.Error.Error()
		}
		info := r.Info
		if info == "" && r.StatusCode > 0 {
			info = fmt.Sprintf("%d", r.StatusCode)
		}
		fmt.Printf("%s,%s,%d,%d,%t,%s\n",
			r.URL, info, r.Duration.Milliseconds(), r.Size, r.Success, errStr)
	}
}

func printMinimal(results []checker.Result) {
	for _, r := range results {
		if r.Error != nil {
			fmt.Printf("ERR %s\n", r.URL)
		} else {
			info := r.Info
			if info == "" {
				info = fmt.Sprintf("%d", r.StatusCode)
			}
			fmt.Printf("%s %s\n", info, r.URL)
		}
	}
}
