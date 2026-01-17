package gitlab

import (
	"context"
	"fmt"

	"github.com/ignorant05/Uniflow/internal/config"
	"github.com/ignorant05/Uniflow/types"
)

// NewClientFromConfig creates new GitLab client from configuration.
//
// Parameters:
//   - profileName: user's profile name
//
// Returns an error if:
//   - invalid configuration
//   - profile not found
//   - gitlab client creation failure
//
// Example:
//
//	client, err := NewClientFromConfig("default")
func NewClientFromConfig(profileName string) (*Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("<?> Error: Failed to load configuration file.\n<?> Error: %w", err)
	}

	profile, err := cfg.GetProfile(profileName)
	if err != nil {
		return nil, fmt.Errorf("<?> Error: Failed to get profile.\n<?> Error: %w", err)
	}

	ctx := context.Background()
	client, err := NewClientFromProfile(ctx, profile)
	if err != nil {
		return nil, fmt.Errorf("<?> Error: Failed to create gitlab client.\n<?> Error: %w", err)
	}

	return client, nil
}

// NewDefaultClient creates new client using default profileName.
//
// Parameters:
//   - None
//
// Returns an error if:
//   - invalid configuration
//   - gitlab client creation failure
//
// Example:
//
//	client, err := NewDefaultClient()
func NewDefaultClient() (*Client, error) {
	return NewClientFromConfig("default")
}

// GetRepositoryInfo retrieves the repository info of a certain repo.
//
// Parameters:
//   - owner: repository owner/group
//   - repo: repository name
//
// Returns an error if:
//   - repository doesn't exist
//
// Example:
//
//	info, err := client.GetRepositoryInfo("owner", "repo")
func (c *Client) GetRepositoryInfo(owner, repo string) (*types.RepositoryInfo, error) {
	projectID := owner + "/" + repo
	project, _, err := c.Projects.GetProject(projectID, nil)
	if err != nil {
		return nil, fmt.Errorf("<?> Error: Failed to get repository %s info.\n<?> Error: %w", repo, err)
	}

	return &types.RepositoryInfo{
		Name:          project.Name,
		FullName:      project.PathWithNamespace,
		Description:   project.Description,
		DefaultBranch: project.DefaultBranch,
		Private:       project.Visibility != "public",
		HTMLURL:       project.WebURL,
	}, nil
}
