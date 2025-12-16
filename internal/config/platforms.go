package config

// NOTE: The base configuration of each platform is defined here

// GitHub base configuration
type GithubConfig struct {
	Token             string `yaml:"token" mapstructure:"token"`
	DefaultRepository string `yaml:"default_repository,omitempty" mapstructure:"default_repository"`
	BaseURL           string `yaml:"base_url,omitempty" mapstructure:"base_url"`
}
