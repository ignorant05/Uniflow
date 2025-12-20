package helpers

import "fmt"

type WorkflowRunSummary struct {
	ID         int64
	Name       string
	Status     string
	Conclusion string
	CreatedAt  string
	UpdatedAt  string
	HTMLURL    string
}

type RepositoryInfo struct {
	Name          string
	FullName      string
	Description   string
	DefaultBranch string
	Private       bool
	HTMLURL       string
}

func ParseRepository(s string) (owner, repo string, err error) {
	parts := splitRepository(s)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("<?> Error: Invalid format.\n<.> Info: Must be in <owner/repo> format.\n")
	}

	return parts[0], parts[1], nil
}

func splitRepository(repo string) []string {
	var parts []string
	current := ""
	for _, c := range repo {
		if c == '/' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts
}
