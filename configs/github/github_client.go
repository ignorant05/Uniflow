package github

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v57/github"
	"github.com/ignorant05/Uniflow/configs/github/helpers"
	"github.com/ignorant05/Uniflow/internal/config"
	constants "github.com/ignorant05/Uniflow/internal/constants/config"
	"golang.org/x/oauth2"
)

type Client struct {
	*github.Client
	Ctx    context.Context
	Config *config.GithubConfig
}

func NewClient(ctx context.Context, cfg *config.GithubConfig) (*Client, error) {
	if cfg.Token == "" {
		cfg.Token = os.Getenv(constants.GITHUB_TOKEN_ENV_VAR_NAME)
		if cfg.Token == "" {
			return nil, fmt.Errorf("<?> Error: No environment variable named %s found.\n<.> Please verify your ~/.zshrc (or ~/.bashrc) file.\n", constants.GITHUB_TOKEN_ENV_VAR_NAME)
		}
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.Token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	if cfg.BaseURL != "" && cfg.BaseURL != "https://api.github.com" {
		var err error

		client, err = github.NewClient(tc).WithEnterpriseURLs(cfg.BaseURL, cfg.BaseURL)
		if err != nil {
			return nil, fmt.Errorf("<?>Error: Failed to create enterprise client.\n<?> Error: %w.\n", err)
		}
	}

	return &Client{
			Client: client,
			Ctx:    ctx,
			Config: cfg},
		nil
}

func NewClientFromProfile(ctx context.Context, profile *config.Profile) (*Client, error) {
	if profile.Github == nil {
		return nil, fmt.Errorf("<?> Error: Github isn't configured for this profile.\n")
	}

	return NewClient(ctx, profile.Github)
}

func (c *Client) GetDefaultRepository() (owner, repo string, err error) {
	if c.Config.DefaultRepository == "" {
		return "", "", fmt.Errorf("<?> Error: No default repository configured.\n")
	}

	return helpers.ParseRepository(c.Config.DefaultRepository)
}
