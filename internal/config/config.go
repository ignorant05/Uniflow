package config

import (
	"fmt"

	"github.com/ignorant05/Uniflow/cmd/constants"
)

// The main configuration structure
type Config struct {
	DefaultPlatform string              `yaml:"default_platform" mapstructure:"default_platform"`
	Profiles        map[string]*Profile `yaml:"profiles" mapstructure:"profiles"`
	Version         string              `yaml:"version" mapstructure:"version"`
}

// The configuration profile (dev, prod, staging, etc...)
type Profile struct {
	Github *GithubConfig `yaml:"github,omitempty" mapstructure:"github"`
}

// NewDefaultConfig creates configuration with default values
func NewDefaultConfig() *Config {
	return &Config{
		DefaultPlatform: constants.DEFAULT_CONFIG_PLATFORM,
		Version:         constants.DEFAULT_CONFIG_VERSION,
		Profiles: map[string]*Profile{
			constants.DEFAULT_CONFIG_PROFILE: {
				Github: &GithubConfig{
					Token:             constants.DEFAULT_GITHUB_TOKEN_PLACEHOLDER,
					DefaultRepository: constants.DEFAULT_GITHUB_REPOSITORY,
					BaseURL:           constants.DEFAULT_GITHUB_BASE_URL,
				},
			},
		},
	}
}

// GetProfile retrieves profile using profile_name: name
//
// Parameters:
//   - name: profile name
//
// Error possible causes:
//   - profile not found (not configured with name: name)
//   - platform isn't configured for this profile
//
// Examples:
// profile, err := cfg.GetProfile("production")
func (cfg *Config) GetProfile(name string) (*Profile, error) {
	if name == "" {
		name = constants.DEFAULT_CONFIG_PROFILE
	}

	profile, exists := cfg.Profiles[name]
	if !exists {
		return nil, fmt.Errorf("<?> Error: No profile named %s registered, please re-check the username and try again...\n\n", name)
	}

	return profile, nil
}
