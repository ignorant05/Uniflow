package config

import (
	"fmt"
	"slices"
	"strings"

	constants "github.com/ignorant05/Uniflow/internal/constants/config"
	"github.com/ignorant05/Uniflow/internal/helpers"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

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

	return fmt.Errorf("configuration validation failed with %d error(s)", len(errors))

}
