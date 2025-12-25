package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// root flags
var (
	// current tool version
	version = "1.0.0"

	// verbose output (global)
	verbose bool
)

// Uniflow command initialization
var rootCmd = &cobra.Command{
	Use:   "Uniflow",
	Short: "A powerful workflow orchestration tool",
	Long: `uniflow is a CLI tool for managing and triggering automated workflows.
It provides commands to initialize configurations, trigger workflows, check status, and view logs.`,
	Version: version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// verbose flag
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output for debugging")

	// version
	rootCmd.SetVersionTemplate(`{{.Version}}`)
}
