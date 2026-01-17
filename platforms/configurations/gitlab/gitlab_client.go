package gitlab

import (
	"context"
	"fmt"
	"os"

	"github.com/ignorant05/Uniflow/internal/config"
	"github.com/ignorant05/Uniflow/platforms/configurations/gitlab/constants"
	"github.com/ignorant05/Uniflow/platforms/configurations/gitlab/helpers"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type Client struct {
	*gitlab.Client
	Ctx    context.Context
	Config *config.GitlabConfig
}

// NewClient creates new GitLab client from configuration.
func NewClient(ctx context.Context, cfg *config.GitlabConfig) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration cannot be nil")
	}

	token := cfg.Token
	if token == "" {
		token = os.Getenv(constants.GITLAB_TOKEN_ENV_VAR_NAME)
	}

	if token == "" {
		return nil, fmt.Errorf(
			"no GitLab token found. Set it in config or %s environment variable",
			constants.GITLAB_TOKEN_ENV_VAR_NAME,
		)
	}

	apiClient, err := gitlab.NewClient(token)
	if err != nil {
		return nil, fmt.Errorf(
			"error while trying to crate client with token: %s,\nerror: %w",
			token,
			err,
		)
	}

	return &Client{
		Client: apiClient,
		Ctx:    ctx,
		Config: cfg,
	}, nil
}

// NewClientFromProfile creates new client from profile configuration.
func NewClientFromProfile(ctx context.Context, profile *config.Profile) (*Client, error) {
	if profile.Gitlab == nil {
		return nil, fmt.Errorf("GitLab isn't configured for this profile")
	}

	return NewClient(ctx, profile.Gitlab)
}

// GetDefaultRepository parses owner/repo from the default repository string.
func (c *Client) GetDefaultRepository() (owner, repo string, err error) {
	if c.Config.DefaultRepository == "" {
		return "", "", fmt.Errorf("no default repository configured")
	}

	return helpers.ParseRepository(c.Config.DefaultRepository)
}

// GetRepositoryPath returns the full repository path
func (c *Client) GetRepositoryPath() (string, error) {
	if c.Config.DefaultRepository == "" {
		return "", fmt.Errorf("no default repository configured")
	}
	return c.Config.DefaultRepository, nil
}
