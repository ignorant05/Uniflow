package constants

// Defaults
const (
	// validation patform default field name
	VALIDATOR_PLATFORM = "default_platform"

	// validation profiles default field name
	VALIDATOR_PROFILES = "profiles"
)

var (
	// supported platforms
	// NOTE: only github is supported right now
	// NOTE: any other platform must be included here
	ValidPlarforms = []string{"github"}
)
