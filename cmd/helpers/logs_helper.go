package helpers

import (
	"fmt"

	"github.com/google/go-github/v57/github"
	"github.com/ignorant05/Uniflow/platforms"
)

// ExtractGithubClient extracts github client from PlatformClient struct
func ExtractGithubClient(client platforms.PlatformClient) (*github.Client, error) {
	if g, ok := client.(interface {
		GetGithubClient() (*github.Client, error)
	}); ok {
		return g.GetGithubClient()
	}

	return nil, fmt.Errorf("<?> Error: Not a GitHub client")
}
