package github

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v57/github"
	"github.com/ignorant05/Uniflow/internal/config"
	constants "github.com/ignorant05/Uniflow/internal/constants/config"
	"github.com/ignorant05/Uniflow/platforms/github/helpers"
	"golang.org/x/oauth2"
)

type Client struct {
	*github.Client
	Ctx    context.Context
	Config *config.GithubConfig
}

// NewClient creates new client from configuration.
//
// Parameters:
//   - ctx: context
//   - cfg: user's github configuration
//
// Returns an error if:
//   - invalid configuration
//   - github client creation failure
//   - invalid enterprise baseURL
//
// Example:
//
//	client, err, err := NewClient(context.Background(), cfg)
func NewClient(ctx context.Context, cfg *config.GithubConfig) (*Client, error) {
	if cfg.Token == "" {
		cfg.Token = os.Getenv(constants.GITHUB_TOKEN_ENV_VAR_NAME)
		if cfg.Token == "" {
			return nil, fmt.Errorf("<?> Error: No environment variable named %s found.\n<.> Please verify your ~/.zshrc (or ~/.bashrc) file", constants.GITHUB_TOKEN_ENV_VAR_NAME+"")
		}
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.Token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	if cfg.BaseURL != "" && cfg.BaseURL != "https://api.github.com" {
		var err error

		client, err = github.NewClient(tc).WithEnterpriseURLs(cfg.BaseURL, cfg.BaseURL)
		if err != nil {
			return nil, fmt.Errorf("<?>Error: Failed to create enterprise client.\n<?> Error: %w", err)
		}
	}

	return &Client{
			Client: client,
			Ctx:    ctx,
			Config: cfg},
		nil
}

// NewClientFromProfile creates new client from profile configuration.
//
// Parameters:
//   - ctx: context
//   - profile: user's profile configuration
//
// Returns an error if:
//   - invalid profile configuration
//   - github client creation failure (github isn't configured for this profile error)
//
// Example:
//
//	client, err, err := NewClientFromProfile(context.Background(), profile)
func NewClientFromProfile(ctx context.Context, profile *config.Profile) (*Client, error) {
	if profile.Github == nil {
		return nil, fmt.Errorf("<?> Error: Github isn't configured for this profile")
	}

	return NewClient(ctx, profile.Github)
}

// GetDefaultRepository retrieves owner and repo for the current configuration.
//
// Parameters:
//   - None
//
// Returns an error if:
//   - no default repo is configured
//
// Example:
//
//	owner, repo, err := client.GetDefaultRepository()
func (c *Client) GetDefaultRepository() (owner, repo string, err error) {
	if c.Config.DefaultRepository == "" {
		return "", "", fmt.Errorf("<?> Error: No default repository configured")
	}

	return helpers.ParseRepository(c.Config.DefaultRepository)
}
