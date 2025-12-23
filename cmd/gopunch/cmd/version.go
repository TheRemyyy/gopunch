package cmd

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan, color.Bold)
		white := color.New(color.FgWhite)

		cyan.Println("⚡ GoPunch")
		white.Printf("  Version:    %s\n", Version)
		white.Printf("  Go:         %s\n", runtime.Version())
		white.Printf("  OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
		if BuildDate != "dev" {
			white.Printf("  Built:      %s\n", BuildDate)
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
