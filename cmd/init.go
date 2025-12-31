package cmd

import (
	"fmt"
	"os"

	"github.com/ignorant05/Uniflow/internal/config"
	errorhandling "github.com/ignorant05/Uniflow/internal/errorHandling"
	"github.com/ignorant05/Uniflow/internal/helpers"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"
)

// Init command flags representatives
var (
	// --force (-f) flag
	// UTILITY: Forces Initialization
	forceInit bool

	// --verbose (-v)
	// UTILITY: verbose output
	initVerbose bool
)

// command: init
//
// Example:
//   - uniflow init
//   - uniflow init --force
var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "Initialize uniflow configuration",
	Long: `Initialize creates the necessary configuration files and directories
for uniflow to function properly. 

Note: This should be run once before using other commands.

Example:

	# Inilialize configuration (default config and it's crusual)
	uniflow init

	# Activating verbose output 
	uniflow i --verbose`,
	Run: runInit,
}

// Commands and subcommands configuration
func init() {
	initCmd.Flags().BoolVarP(&forceInit, "force", "f", false, "Overwrite an existing configuration")
	initCmd.Flags().BoolVarP(&initVerbose, "verbose", "v", false, "verbose output")
	rootCmd.AddCommand(initCmd)
}

// runInit is the init command main function
//
// Errors possible causes:
//
//   - Dir already exists with that name
//   - Config file already exists with that name
//   - Permission error
func runInit(cmd *cobra.Command, args []string) {
	// if verbose mode active
	if initVerbose {
		fmt.Println("</> Info: Running in verbose mode...")
	}
	fmt.Println("❯ Initialize Uniflow configuration...")

	// Getting configuration dir
	configDir, err := helpers.GetConfigDir()
	if err != nil {
		errMsg := fmt.Errorf("<?> Error: Failed to resolve configuration directory.\n<?> Error: %w", err)
		errorhandling.HandleError(errMsg)
	}

	// Creating new directory
	if err := os.MkdirAll(configDir, 0700); err != nil {
		errMsg := fmt.Errorf("<?> Error: Failed to create configuration directory.\n<?> Error: %w", err)
		errorhandling.HandleError(errMsg)
	}
	fmt.Printf("✓ Created configuration directory: %s\n", configDir)

	// Get configuration path
	configPath, err := helpers.GetConfigPath()
	if err != nil {
		errMsg := fmt.Errorf("<?> Error: Failed to resolve configuration path.\n<?> Error: %w", err)
		errorhandling.HandleError(errMsg)
	}

	// Returns error
	// To avoid overwriting existing files without user permissions
	if _, err := os.Stat(configPath); err == nil && !forceInit {
		errMsg := fmt.Errorf("<?> Error: Configuration path already exists: %s.\n<?> Error: %w.\n<.> Solution: use --force (or -f) to overwrite.\n", configPath, err)
		errorhandling.HandleError(errMsg)
	}

	// Default configuration
	cfg := config.NewDefaultConfig()

	// Generate configuration
	data, err := yaml.Marshal(cfg)
	if err != nil {
		errMsg := fmt.Errorf("<?> Error: Failed to generate configuration file.\n<?> Error: %w.\n", err)
		errorhandling.HandleError(errMsg)
	}

	// Write to file: saving data
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		errMsg := fmt.Errorf("<?> Error: Failed to write configuration file at %s.\n<?> Error: %w.\n", configPath, err)
		errorhandling.HandleError(errMsg)
	}
	fmt.Printf("✓ Created configuration file at: %s.", configPath)

	// creating logs directory
	logsDir := configDir + "/logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		errMsg := fmt.Errorf("<?> Error: Failed to create logs directory.\n<?> Error: %w", err)
		errorhandling.HandleError(errMsg)
	}
	fmt.Printf("✓ Created logs directory: %s\n", logsDir)

	fmt.Println("\n✓ Initialization complete!")
	fmt.Println("\n❯ Next steps:")
	fmt.Println("  1. Set your API tokens as environment variables:")
	fmt.Println("     export GITHUB_TOKEN=your_token_here")
	fmt.Println("\n  2. Or update the config file directly:")
	fmt.Printf("     %s\n", configPath)
	fmt.Println("\n  3. Verify your configuration:")
	fmt.Println("     uniflow config list")
	fmt.Println("\n  4. Start triggering workflows:")
	fmt.Println("     uniflow trigger <workflow>")
}
