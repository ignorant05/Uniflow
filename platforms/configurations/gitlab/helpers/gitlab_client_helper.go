package helpers

import (
	"fmt"
	"strings"
)

// WorkflowRunSummary represents summary structure.
// It contains pipeline information such as: ID, Name, Status, Conclusion, CreatedAt, UpdatedAt and HTMLURL.
type PipelineSummary struct {
	ID        int64
	Ref       string
	CreatedAt string
	UpdatedAt string
	Status    string
	WebURL    string
}

// RepositoryInfo represents repitory info structure.
// It contains repository information such as: Name, FullName, Description, DefaultBranch, Private and HTMLURL.
type RepositoryInfo struct {
	Name          string
	FullName      string
	Description   string
	DefaultBranch string
	Private       bool
	HTMLURL       string
}

// parseRepository splits "owner/repo" or "group/subgroup/repo" strings.
// For GitLab, it's common to have nested groups: "group/subgroup/project"
func ParseRepository(repoPath string) (owner, repo string, err error) {
	parts := strings.Split(repoPath, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid repository path format: %s", repoPath)
	}

	// Last part is the project name
	repo = parts[len(parts)-1]
	// Everything before last part is the namespace/owner
	owner = strings.Join(parts[:len(parts)-1], "/")

	return owner, repo, nil
}
