package config

import (
	"fmt"
	"strings"

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
	Github  *GithubConfig   `yaml:"github,omitempty" mapstructure:"github"`
	Jenkins *JenkinsConfig `yaml:"jenkins,omitempty" mapstructure:"jenkins"`
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
		return nil, fmt.Errorf("<?> Error: No profile named %s registered, please re-check the username and try again", name)
	}

	return profile, nil
}

// GetPlatform retrieves platform using platform name: name
//
// Parameters:
//   - platformName: platform name
//
// Error possible causes:
//   - platform isn't configured for this profile
//
// Examples:
// platform, err := p.GetPlatform("github")
func (p *Profile) GetPlatform(platformName string) (interface{}, error) {
	switch strings.ToLower(platformName) {
	case strings.ToLower(constants.DEFAULT_CONFIG_PLATFORM):
		if p.Github == nil {
			return nil, fmt.Errorf("<?> Error: Github configuration not found for this profile")
		}
	default:
		return nil, fmt.Errorf("<?> Error: Unsupported platform: %s", platformName)
	}

	return nil, nil
}
