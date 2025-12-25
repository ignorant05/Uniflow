package config

import (
	"fmt"
	"slices"
	"strings"

	constants "github.com/ignorant05/Uniflow/internal/constants/config"
	"github.com/ignorant05/Uniflow/internal/helpers"
)

// ValidationError struct
type ValidationError struct {
	// at what level ?
	Field string

	// the error itself
	Message string
}

// Formatting the error
func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Validate validates overall config
//
// Parameters:
//   - None
//
// Error possible causes:
//   - at least one platform must be configured (github)
//   - platform must be valid (github, jenkins, gitlab-ci, circleci)
//   - at least one profile must be configured
//
// Examples:
// errs := cfg.Validate()
func (cfg *Config) Validate() []error {
	var errors []error
	if cfg.DefaultPlatform == "" {
		errors = append(errors, &ValidationError{
			Field:   constants.VALIDATOR_PLATFORM,
			Message: "<?> Error: Default platform cannot be empty",
		})

	}

	if !slices.Contains(constants.ValidPlarforms, cfg.DefaultPlatform) {
		errors = append(errors, &ValidationError{
			Field:   constants.VALIDATOR_PLATFORM,
			Message: fmt.Sprintf("<?> Error: Must be one of: %s", strings.Join(constants.ValidPlarforms, ", ")),
		})
	}

	if len(cfg.Profiles) == 0 {
		errors = append(errors, &ValidationError{
			Field:   constants.VALIDATOR_PROFILES,
			Message: "<?> Error: At least one profile must be defined",
		})
	}

	for profileName, profile := range cfg.Profiles {
		profileErrors := ValidateProfiles(profileName, profile)
		errors = append(errors, profileErrors...)
	}

	return errors
}

// ValidateProfiles validates profile
//
// Parameters:
//   - key: field name
//   - val: field value
//
// Error possible causes:
//   - no platform configured
//
// Examples:
// errs := ValidateProfiles("prod", profile)
func ValidateProfiles(name string, profile *Profile) []error {
	var errors []error
	prefix := fmt.Sprintf("profiles.%s", name)

	if profile.Github == nil {
		errors = append(errors, &ValidationError{
			Field:   prefix,
			Message: "<?> Error: At least one platform must be configured",
		})

	}

	if profile.Github != nil {
		ghErrors := ValidateGithub(prefix+".github", profile.Github)
		errors = append(errors, ghErrors...)
	}

	return errors
}

// ValidateGithub validates github conf
//
// Parameters:
//   - prefix: prefix string
//   - cfg: github configuration struct
//
// Error possible causes:
//   - github token is not sat
//   - invalid url
//
// Examples:
// errs := Update(prefix, cfg)
func ValidateGithub(prefix string, cfg *GithubConfig) []error {
	var errors []error

	if cfg.Token == "" || strings.HasPrefix(cfg.Token, "${") {
		errors = append(errors, &ValidationError{
			Field:   prefix + ".token",
			Message: "<?> Error: Token is required (set via environment variable or directly)",
		})
	}

	if cfg.BaseURL == "" {
		errors = append(errors, &ValidationError{
			Field:   prefix + ".base_url",
			Message: "<?> Error: Must be a valid URL",
		})
	}

	if cfg.DefaultRepository == "" || helpers.IsValidRepoFormat(cfg.DefaultRepository) {
		errors = append(errors, &ValidationError{
			Field:   prefix + ".default_repository",
			Message: "<?> Error: Must be in format 'owner/repo'",
		})
	}

	return errors
}

// ValidateAndReport validates configuration
//
// Parameters:
//   - cfg: configuration struct
//
// Examples:
// err := ValidateAndReport(cfg)
func ValidateAndReport(cfg *Config) error {
	errors := cfg.Validate()
	if len(errors) == 0 {
		return nil
	}

	fmt.Println("==========================================")
	fmt.Println("=====Configuration validation failed:=====")
	fmt.Println("==========================================")

	for i, err := range errors {
		fmt.Printf("  %d. %s\n", i+1, err.Error())
	}

	fmt.Println("\nPlease fix these issues in your config file or use 'uniflow config set' to update values.")

	return fmt.Errorf("configuration validation failed with %d error(s)\n\n", len(errors))

}
