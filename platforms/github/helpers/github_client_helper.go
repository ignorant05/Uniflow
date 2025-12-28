package helpers

import "fmt"

// WorkflowRunSummary represents summary structure.
// It contains workflow information such as: ID, Name, Status, Conclusion, CreatedAt, UpdatedAt and HTMLURL.
type WorkflowRunSummary struct {
	ID         int64
	Name       string
	Status     string
	Conclusion string
	CreatedAt  string
	UpdatedAt  string
	HTMLURL    string
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

// ParseRepository is a helper function that parses the owner and repo.
//
// Parameters:
//   - s: (eg: "owner/repo")
//
// Return an error if:
//    invalid string forma
//
// Example:
// owner, repo := helpers.ParseRepository("ignorant05/Uniflow")
func ParseRepository(s string) (owner, repo string, err error) {
	parts := splitRepository(s)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("<?> Error: Invalid format.\n</> Info: Must be in <owner/repo> format.\n\n")
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
