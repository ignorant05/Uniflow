package helpers

import "strings"

func IsValidRepoFormat(repo string) bool {
	parts := strings.Split(repo, "/")
	return len(parts) == 2 && parts[0] != "" && parts[1] != ""
}
