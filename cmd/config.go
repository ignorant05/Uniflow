package cmd

import (
	"fmt"
	"strings"

	"github.com/ignorant05/Uniflow/cmd/helpers"
	"github.com/ignorant05/Uniflow/internal/config"
	"github.com/ignorant05/Uniflow/internal/constants"
	"github.com/spf13/cobra"
)

var (
	profileFlag string
	showSecrets bool
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Uniflow configuration",
	Long: `Manage configuration settings for Uniflow,

Available subcommands: 
	list	 - Show current configuration
	set		 - Update configuration values
	get		 - Get a specific configuration values
	validate - Validate configuration file`,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show current configuration",
	Long: `Display the current configuration settings.

Example:
	uniflow config list
	uniflow config list --profile prod
	uniflow config list --show-secrets`,
	RunE: runConfigList,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value> \nNote:\n  key format: profiles.<profile>.<platform>.<field>",
	Short: "Update a configuration value",
	Long: `Update a specific configuration value.

Key format: profiles.<profile>.<platform>.<field>

Examples:
	uniflow config set profiles.default.github.default_repository your/repository
	uniflow config set default_platform github`,
	Args: cobra.ExactArgs(2),
	RunE: runConfigSet,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Show configuration value",
	Long: `Show a specific configuration value.

Key format: profiles.<profile>.<platform>.<field>

Examples:
	uniflow config get profiles.default.github.default_repository 
	uniflow config get default_platform`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigGet,
}

var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration file",
	Long: `Check if the configuration file is valid and all required fields are set.

Example:
	uniflow config validate`,
	RunE: runConfigValidate,
}

func init() {
	configListCmd.Flags().StringVarP(&profileFlag, "profile", "p", "default", "Profile to display")
	configListCmd.Flags().BoolVar(&showSecrets, "show-secrets", false, "Show sensitive values (tokens)")

	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configValidateCmd)

	rootCmd.AddCommand(configCmd)
}

func runConfigList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Printf("<.> Info: Configuration (Profile: %s)\n", profileFlag)
	fmt.Println(strings.Repeat("─", 60))

	fmt.Printf("<.> Info: Show default platform: %s\n", cfg.DefaultPlatform)
	fmt.Printf("<.> Info: Show version: %s\n", cfg.Version)

	profile, err := cfg.GetProfile(profileFlag)
	if err != nil {
		return err
	}

	if profile.Github != nil {
		fmt.Println("\n❯❯❯ GitHub:")
		fmt.Printf("  Base URL:     %s\n", profile.Github.BaseURL)
		fmt.Printf("  Default Repository: %s\n", helpers.ValueOrEmpty(profile.Github.DefaultRepository))
		fmt.Printf("  Token:        %s\n", helpers.MaskSecret(profile.Github.Token, showSecrets))
	}

	fmt.Println()
	fmt.Printf("\nAvailable Profiles: %s\n", strings.Join(getProfileNames(cfg), ", "))

	return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key, val := args[0], args[1]

	if verbose {
		fmt.Printf("<.> Info: Setting %s = %s\n", key, val)
	}

	if err := config.Update(key, val); err != nil {
		return err
	}

	fmt.Printf("<✓> Updated %s\n", key)

	return nil
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	key := args[0]

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	parts := strings.Split(key, ".")

	switch parts[0] {
	case constants.DEFAULT_PLATFORM:
		fmt.Println(cfg.DefaultPlatform)
	case constants.VERSION:
		fmt.Println(cfg.Version)
	case constants.PROFILES:
		if len(parts) < 4 {
			return fmt.Errorf("<?> Error: Invalid key format.\n<.> Use: profiles.<profile>.<platform>.<field>\n")
		}

		profileName, platform, filed := parts[1], parts[2], parts[3]

		profile, err := cfg.GetProfile(profileName)
		if err != nil {
			return err
		}

		val, err := getPlatformFiledValue(profile, platform, filed)
		if err != nil {
			return err
		}

		fmt.Println(val)

	default:
		return fmt.Errorf("<?> Error: Unknown config key: %s", key)

	}

	return nil
}

func runConfigValidate(cmd *cobra.Command, args []string) error {
	fmt.Println("❯❯❯ Validating...")
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if err := config.ValidateAndReport(cfg); err != nil {
		return err
	}

	fmt.Println("<✓> Configuration is valid!")
	return nil
}

func getProfileNames(cfg *config.Config) []string {
	profiles := cfg.Profiles
	names := make([]string, 0, len(profiles))
	for name := range profiles {
		names = append(names, name)
	}

	return names
}

func getPlatformFiledValue(profile *config.Profile, platform, field string) (string, error) {
	switch platform {
	case constants.GITHUB:
		{
			switch field {
			case constants.TOKEN_FIELD:
				return profile.Github.Token, nil
			case constants.DEFAULT_REPOSITORY_FIELD:
				return profile.Github.DefaultRepository, nil
			case constants.BASE_URL_FIELD:
				return profile.Github.BaseURL, nil
			default:
				return "", fmt.Errorf("<?> Error: Invalid field: %s.\n", field)
			}
		}
	}

	return "", fmt.Errorf("<?> Error: Unknown platform: %s.\n", platform)
}
