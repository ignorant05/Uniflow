package jenkins

import (
	"context"
	"fmt"
	"os"

	"github.com/bndr/gojenkins"
	"github.com/ignorant05/Uniflow/internal/config"
	"github.com/ignorant05/Uniflow/platforms/configurations/jenkins/constants"
)

type Client struct {
	Client *gojenkins.Jenkins
	Ctx    context.Context
	Config *config.JenkinsConfig
}

// NewClient creates new Jenkins client from configuration.
//
// Parameters:
//   - ctx: context
//   - cfg: user's Jenkins configuration
//
// Returns an error if:
//   - invalid configuration
//   - Jenkins client creation failure
//   - connection initialization failure
//
// Example:
//
//	client, err := NewClient(context.Background(), cfg)
func NewClient(ctx context.Context, cfg *config.JenkinsConfig) (*Client, error) {
	// Validate configuration
	if err := validateConfig(cfg); err == nil {
		return nil, err
	}

	// Create Jenkins client
	client := gojenkins.CreateJenkins(nil, cfg.BaseURL, cfg.Username, cfg.APIToken)

	// Initialize connection (verifies credentials and connection)
	_, err := client.Init(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Jenkins client: %w", err)
	}

	return &Client{
		Client: client,
		Ctx:    ctx,
		Config: cfg,
	}, nil
}

func validateConfig(cfg *config.JenkinsConfig) error {
	if cfg == nil {
		return fmt.Errorf("jenkins configuration is nil")
	}

	if cfg.BaseURL == "" {
		return fmt.Errorf("base_url is required for Jenkins configuration")
	}

	if cfg.Username == "" {
		cfg.Username = os.Getenv(constants.JENKINS_USERNAME_ENV_VAR_NAME)
		if cfg.Username == "" {
			return fmt.Errorf("username is required for Jenkins configuration")
		}
	}

	if cfg.APIToken == "" {
		cfg.APIToken = os.Getenv(constants.JENKINS_TOKEN_ENV_VAR_NAME)
		if cfg.APIToken == "" {
			return fmt.Errorf("api_token is required for Jenkins configuration")
		}
	}
	return nil
}

// NewClientFromProfile creates new Jenkins client from profile configuration.
//
// Parameters:
//   - ctx: context
//   - profile: user's profile configuration
//
// Returns an error if:
//   - invalid profile configuration
//   - Jenkins isn't configured for this profile
//   - client creation failure
//
// Example:
//
//	client, err := NewClientFromProfile(context.Background(), profile)
func NewClientFromProfile(ctx context.Context, profile *config.Profile) (*Client, error) {
	if profile == nil {
		return nil, fmt.Errorf("profile is nil")
	}

	if profile.Jenkins == nil {
		return nil, fmt.Errorf("jenkins isn't configured for this profile")
	}

	return NewClient(ctx, profile.Jenkins)
}

// NewTestClient creates a Jenkins client for testing without making HTTP calls
func NewTestClient(ctx context.Context, cfg *config.JenkinsConfig) (*Client, error) {
	// Copy validation logic from NewClient but skip Init()
	if cfg == nil {
		return nil, fmt.Errorf("jenkins configuration is nil")
	}

	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("base_url is required for Jenkins configuration")
	}

	if cfg.Username == "" {
		return nil, fmt.Errorf("username is required for Jenkins configuration")
	}

	// For testing, allow empty token if password is provided
	if cfg.APIToken == "" && cfg.Password == "" {
		return nil, fmt.Errorf("either api_token or password is required for Jenkins configuration")
	}

	// Create Jenkins client but don't call Init
	var jenkinsClient *gojenkins.Jenkins
	if cfg.APIToken != "" {
		jenkinsClient = gojenkins.CreateJenkins(nil, cfg.BaseURL, cfg.Username, cfg.APIToken)
	} else {
		jenkinsClient = gojenkins.CreateJenkins(nil, cfg.BaseURL, cfg.Username, cfg.Password)
	}

	return &Client{
		Client: jenkinsClient,
		Ctx:    ctx,
		Config: cfg,
	}, nil
}
