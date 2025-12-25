package constants

// Default configuration values.
const (
	// DEFAULT_CONFIG_PLATFORM is the default platform to be configured
	DEFAULT_CONFIG_PLATFORM = "github"

	// DEFAULT_CONFIG_VERSION is the default version to be configured
	// This is often updated incrementally
	DEFAULT_CONFIG_VERSION = "1.0"

	// DEFAULT_CONFIG_PROFIL is the default profile to be configured
	DEFAULT_CONFIG_PROFILE = "default"
)

// Default placeholders for github default configuration.
const (
	// DEFAULT_GITHUB_TOKEN_PLACEHOLDER is a placeholder for github token.
	// Must be configured in ~/.zshrc (or ~/.bashrc) file
	// Or use export GITHUB_TOKEN="your token here" in your terminal
	DEFAULT_GITHUB_TOKEN_PLACEHOLDER = "${GITHUB_TOKEN}"

	// DEFAULT_GITHUB_REPOSITORY is the default repo name to be configured.
	DEFAULT_GITHUB_REPOSITORY = ""

	// DEFAULT_GITHUB_BASE_URL is the default baseURL value.
	DEFAULT_GITHUB_BASE_URL = "https://api.github.com"
)
