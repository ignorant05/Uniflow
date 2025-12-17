package helpers

import "strings"

func MaskSecret(val string, show bool) string {
	if show || val == "" || strings.HasPrefix(val, "${") {
		return val
	}

	if len(val) > 8 && show {
		return val[:4] + "..." + val[len(val)-4:]
	}

	return "*** (Use --show-secrets flag to see through)"
}

func ValueOrEmpty(val string) string {
	if val == "" {
		return "(not set)"
	}

	return val
}
