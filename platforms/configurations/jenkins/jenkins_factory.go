package jenkins

import (
	"context"
	"fmt"

	"github.com/ignorant05/Uniflow/internal/config"
)

// NewClientFromConfig creates new client from configuration.
//
// Parameters:
//   - profileName: user's profile name
//
// Returns an error if:
//   - invalid configuration
//   - github client creation failure
//
// Example:
//
//	client, err := NewClientFromConfig("default")
func NewClientFromConfig(profileName string) (*Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration file.\n<?> Error: %w", err)
	}

	profile, err := cfg.GetProfile(profileName)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile.\n<?> Error: %w", err)
	}

	if profile.Jenkins == nil {
		return nil, fmt.Errorf("profile '%s' does not have Jenkins configured", profileName)
	}

	ctx := context.Background()
	client, err := NewClientFromProfile(ctx, profile)
	if err != nil {
		return nil, fmt.Errorf("failed to create jenkins client.\n<?> Error: %w", err)
	}

	return client, nil
}

// NewDefaultClient creates new Jenkins client using default profile.
//
// Parameters:
//   - None
//
// Returns an error if:
//   - invalid configuration
//   - default profile not found
//   - Jenkins not configured for default profile
//   - Jenkins client creation failure
//
// Example:
//
//	client, err := NewDefaultClient()
func NewDefaultClient() (*Client, error) {
	return NewClientFromConfig("default")
}

// NewClientFromProfileName is a convenience wrapper.
//
// Parameters:
//   - profileName: user's profile name
//
// Returns:
//   - Jenkins client
//   - error if any
func NewClientFromProfileName(profileName string) (*Client, error) {
	return NewClientFromConfig(profileName)
}
