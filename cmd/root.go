package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
	verbose bool
)

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
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output for debugging")
	rootCmd.SetVersionTemplate(`{{.Version}}`)
}
