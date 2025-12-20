package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v57/github"
	"github.com/ignorant05/Uniflow/configs/github/helpers"
	"github.com/ignorant05/Uniflow/internal/config"
)

func NewClientFromConfig(profileName string) (*Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("<?> Error: Failed to load configuration file.\n<?> Error: %w.\n", err)
	}

	profile, err := cfg.GetProfile(profileName)
	if err != nil {
		return nil, fmt.Errorf("<?> Error: Failed to get profile.\n<?> Error: %w.\n", err)
	}

	ctx := context.Background()
	client, err := NewClientFromProfile(ctx, profile)
	if err != nil {
		return nil, fmt.Errorf("<?> Error: Failed to create github client.\n<?> Error: %w.\n", err)
	}

	return client, nil
}

func NewDefaultClient() (*Client, error) {
	return NewClientFromConfig("default")
}

func (c *Client) TestConnection() error {
	user, _, err := c.Client.Users.Get(c.Ctx, "")
	if err != nil {
		if ghErr, ok := err.(*github.ErrorResponse); ok {
			fmt.Printf("GitHub API Error: %v (Status: %v)\n", ghErr.Message, ghErr.Response.StatusCode)
		}
		return fmt.Errorf("authentication failed: %w", err)
	}

	fmt.Printf("<✓> Successfully authenticated as: %s\n", user.GetLogin())
	return nil
}

func (c *Client) GetRepositoryInfo(owner, repo string) (*helpers.RepositoryInfo, error) {
	repository, _, err := c.Client.Repositories.Get(c.Ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("<?> Error: Failed to get repository %s info.\n<?> Error: %w.\n", repo, err)
	}

	return &helpers.RepositoryInfo{
		Name:          repository.GetName(),
		FullName:      repository.GetFullName(),
		Description:   repository.GetDescription(),
		DefaultBranch: repository.GetDefaultBranch(),
		Private:       repository.GetPrivate(),
		HTMLURL:       repository.GetHTMLURL(),
	}, nil

}
