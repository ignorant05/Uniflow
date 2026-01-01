package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	constants "github.com/ignorant05/Uniflow/internal/constants/config"
	"github.com/ignorant05/Uniflow/internal/helpers"
	"github.com/spf13/viper"
)

// Regex pattern
var envVarRegex = regexp.MustCompile(`\$\{([^}]+)\}`)

// Load loads reads and parses configuration file
//
// Parameters:
//   - None
//
// Error possible causes:
//   - configuration file not found
//   - env vars not found
//   - internal error (failed to read/parse config file)
//
// Examples:
// config, err := Load()
func Load() (*Config, error) {
	configPath, err := helpers.GetConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("<?> Error: Configuration file was not found at %s.\nPlease run: 'uniflow init' to create it", configPath)
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("<?> Error: Failed to read configuration file\nError: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("<?> Error: Failed to parse configuration file\nError: %w", err)
	}

	if err := resolveEnvVars(&cfg); err != nil {
		return nil, fmt.Errorf("<?> Error: Failed to resolve environment variables\nError: %w", err)
	}

	return &cfg, nil
}

// resolveEnvVars resolves env vars
//
// Parameters:
//   - cfg: configuration
//
// Examples:
// err := resolveEnvVars(cfg)
func resolveEnvVars(cfg *Config) error {
	for profileName, profile := range cfg.Profiles {
		if profile.Github != nil {
			profile.Github.Token = resolveEnvVar(profile.Github.Token)
			profile.Github.DefaultRepository = resolveEnvVar(profile.Github.DefaultRepository)
			profile.Github.BaseURL = resolveEnvVar(profile.Github.BaseURL)
		}

		cfg.Profiles[profileName] = profile
	}

	return nil
}

// resolveEnvVars resolves env vars
//
// Parameters:
//   - value: replaces value
//
// Examples:
// err := resolveEnvVars("something")
func resolveEnvVar(value string) string {
	return envVarRegex.ReplaceAllStringFunc(value, func(match string) string {
		varName := match[2 : len(value)-1]

		if envVal := os.Getenv(varName); envVal != "" {
			return envVal
		}

		return match
	})
}

// Save configuration
//
// Parameters:
//   - cfg: configuration struct
//
// Error possible causes:
//   - internal error (failed to save config file)
//
// Examples:
// err := Save(cfg)
func Save(cfg *Config) error {
	configPath, err := helpers.GetConfigPath()
	if err != nil {
		return err
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	v.Set(constants.DEFAULT_PLATFORM, cfg.DefaultPlatform)
	v.Set(constants.VERSION, cfg.Version)
	v.Set(constants.PROFILES, cfg.Profiles)

	if err := v.WriteConfig(); err != nil {
		return fmt.Errorf("<?> Error : Failed to save configuration file.\nError: %w", err)
	}

	return nil
}

// Update updates key with val
//
// Parameters:
//   - key: field name
//   - val: field value
//
// Error possible causes:
//   - configuration file not found
//   - env vars not found
//   - internal error (failed to read/parse config file)
//
// Examples:
// err := Update("token", "gibbrich string")
func Update(key, val string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	parts := strings.Split(key, ".")
	if len(parts) < 2 {
		return fmt.Errorf("<?> Error: Invalid config key format. Use: profiles.<profile>.<platform>.<field>")
	}

	fmt.Println(parts)

	switch parts[0] {
	case constants.DEFAULT_PLATFORM:
		cfg.DefaultPlatform = val
	case constants.PROFILES:
		if len(parts) < 4 {
			return fmt.Errorf("invalid profile key format. Use: profiles.<profile>.<platform>.<field>")
		}

		profileName, platform, field := parts[1], parts[2], parts[3]

		profile, err := cfg.GetProfile(profileName)
		if err != nil {
			return err
		}

		if err := updatePlatformField(profile, platform, field, val); err != nil {
			return err
		}

	default:
		return fmt.Errorf("<?> Error: Unknown config section: %s", parts[0])
	}

	return Save(cfg)
}

// updatePlatformField updates platform field
//
// Parameters:
//   - profile: profile object
//   - platform: platform name
//   - field: field name
//   - field value: val name
//
// Error possible causes:
//   - configuration file not found
//   - env vars not found
//   - internal error (failed to read/parse config file)
//
// Examples:
// err := updatePlatformField(profile, "github", "token", "gibbrich string")
func updatePlatformField(profile *Profile, platform, field, val string) error {
	switch platform {
	case constants.GITHUB:
		if profile.Github == nil {
			profile.Github = &GithubConfig{}
		}

		switch field {
		case constants.TOKEN_FIELD:
			profile.Github.Token = val
		case constants.DEFAULT_REPOSITORY_FIELD:
			profile.Github.DefaultRepository = val
		case constants.BASE_URL_FIELD:
			profile.Github.BaseURL = val
		default:
			return fmt.Errorf("<?> Error: Unknown github field: %s", field)
		}
	default:
		return fmt.Errorf("<?> Error: Unsupported platform: %s", platform)
	}

	return nil
}
