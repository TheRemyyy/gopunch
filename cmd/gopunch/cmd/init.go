package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [filename]",
	Short: "Generate a sample configuration file",
	Long: `Create a sample gopunch.json configuration file in the current directory.

Examples:
  gopunch init                    # Creates gopunch.json
  gopunch init myconfig.json      # Creates myconfig.json`,
	Run: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	filename := "gopunch.json"
	if len(args) > 0 {
		filename = args[0]
	}

	// Check if file exists
	if _, err := os.Stat(filename); err == nil {
		color.Red("✗ File %s already exists", filename)
		os.Exit(1)
	}

	config := Config{
		URLs: []string{
			"https://example.com",
			"https://api.example.com/health",
		},
		Interval: 30,
		Timeout:  10,
		Method:   "GET",
		Headers: map[string]string{
			"User-Agent": "GoPunch/2.0",
		},
		Insecure:    false,
		Follow:      true,
		Concurrency: 10,
		Retries:     2,
		Alerting: &AlertConfig{
			Enabled:  false,
			Cooldown: 300,
			Webhook: &WebhookConfig{
				URL:    "https://discord.com/api/webhooks/YOUR_WEBHOOK_ID",
				Method: "POST",
			},
		},
	}

	file, err := os.Create(filename)
	if err != nil {
		color.Red("✗ Failed to create file: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		color.Red("✗ Failed to write config: %v", err)
		os.Exit(1)
	}

	absPath, _ := filepath.Abs(filename)
	color.Green("✓ Created %s", absPath)
	fmt.Println()
	fmt.Println("Edit the file and run:")
	color.Cyan("  gopunch watch --config %s", filename)
	fmt.Println()
}
