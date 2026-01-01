package helpers

import "strings"

func IsError(content string) bool {
	if strings.Contains(content, "error") ||
		strings.Contains(content, "fatal") ||
		strings.Contains(content, "?") ||
		strings.Contains(content, "failed") {
		return true
	}

	return false
}

func IsWarning(content string) bool {
	if strings.Contains(content, "warning") ||
		strings.Contains(content, "!") ||
		strings.Contains(content, "warn") {
		return true
	}

	return false
}

func IsDebug(content string) bool {
	return strings.Contains(content, "debug")
}

func IsSuccess(content string) bool {
	if strings.Contains(content, "success") ||
		strings.Contains(content, "✓") ||
		strings.Contains(content, "passed") ||
		strings.Contains(content, "completed") {
		return true
	}

	return false
}

func FormatLogs(line string) (string, string) {
	if strings.Contains(line, "Z ") {
		parts := strings.SplitN(line, "Z ", 2)
		if len(parts) == 2 {
			return parts[0] + "Z", parts[1]
		}
	}

	return "", ""
}
