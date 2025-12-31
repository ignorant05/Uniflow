package helpers

import (
	"fmt"

	"github.com/ignorant05/Uniflow/platforms"
	"github.com/ignorant05/Uniflow/platforms/github"
)

// ExtractGithubClient extracts github client from PlatformClient struct
func ExtractGithubClient(client platforms.PlatformClient) (*github.Client, error) {
	if g, ok := client.(interface {
		GetGithubClient() (*github.Client, error)
	}); ok {
		return g.GetGithubClient()
	}

	return nil, fmt.Errorf("<?> Error: Not a GitHub client\n")
}
