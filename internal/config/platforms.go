package config

// NOTE: The base configuration of each platform is defined here

// GitHub base configuration
type GithubConfig struct {
	Token             string `yaml:"token" mapstructure:"token"`
	DefaultRepository string `yaml:"default_repository,omitempty" mapstructure:"default_repository"`
	BaseURL           string `yaml:"base_url,omitempty" mapstructure:"base_url"`
}

// Jenkins base configuration
type JenkinsConfig struct {
	BaseURL  string `yaml:"base_url" mapstructure:"base_url"`
	Username string `yaml:"username,omitempty" mapstructure:"username"`
	APIToken string `yaml:"api_token,omitempty" mapstructure:"api_token"`
	Password string `yaml:"password,omitempty" mapstructure:"password"`

	// Jenkins-specific fields
	JobName  string `yaml:"job_name,omitempty" mapstructure:"job_name"`
	ViewName string `yaml:"view_name,omitempty" mapstructure:"view_name"`

	// Optional settings
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify,omitempty" mapstructure:"insecure_skip_verify"`
	TimeoutSeconds     int    `yaml:"timeout_seconds,omitempty" mapstructure:"timeout_seconds"`
	CACertPath         string `yaml:"ca_cert_path,omitempty" mapstructure:"ca_cert_path"`
}
