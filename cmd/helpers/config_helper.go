package helpers

import "strings"

func MaskSecret(val string, show bool, force bool) string {
	if (show || val == "" || strings.HasPrefix(val, "${")) || (len(val) > 8 && show && force) {
		return val
	}

	if len(val) > 8 && show && !force {
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
