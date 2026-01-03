package github_test

import (
	"context"
	"testing"

	"github.com/ignorant05/Uniflow/internal/config"
	"github.com/ignorant05/Uniflow/platforms/configurations/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Testing default repo configured
func TestDefaultRepo_Success(t *testing.T) {
	cfg := &config.GithubConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "ignorant05/Uniflow",
	}

	client, _ := github.NewClient(context.Background(), cfg)
	owner, repo, err := client.GetDefaultRepository()

	require.NoError(t, err)
	assert.Equal(t, owner, "ignorant05")
	assert.Equal(t, repo, "Uniflow")
}

// Testing default repo not configured
func TestDefaultRepo_Failure(t *testing.T) {
	cfg := &config.GithubConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "",
	}

	client, _ := github.NewClient(context.Background(), cfg)
	owner, repo, err := client.GetDefaultRepository()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "No default repository configured")
	assert.Empty(t, owner)
	assert.Empty(t, repo)
}

// Testing getting default repository
func TestGetDefaultRepository_Success(t *testing.T) {
	cfg := &config.GithubConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "ignorant05/Uniflow",
	}

	client, _ := github.NewClient(context.Background(), cfg)

	owner, repo, err := client.GetDefaultRepository()

	require.NoError(t, err)
	assert.Equal(t, "ignorant05", owner)
	assert.Equal(t, "Uniflow", repo)
}

// Testing getting default repository (Failure: not configured)
func TestGetDefaultRepository_Failure(t *testing.T) {
	cfg := &config.GithubConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "",
	}

	client, _ := github.NewClient(context.Background(), cfg)

	_, _, err := client.GetDefaultRepository()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "No default repository")
}

// Testing default repo, table driven
func TestGetDefaultRepository_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		repository    string
		expectError   bool
		expectedOwner string
		expectedRepo  string
	}{
		{
			name:          "lols",
			repository:    "leader/my-lols",
			expectError:   false,
			expectedOwner: "leader",
			expectedRepo:  "my-lols",
		},
		{
			name:        "empty repo",
			repository:  "",
			expectError: true,
		},
		{
			name:          "org of ones",
			repository:    "one/numbers",
			expectError:   false,
			expectedOwner: "one",
			expectedRepo:  "numbers",
		},
		{
			name:          "hehe",
			repository:    "iamHehe/laughts",
			expectError:   false,
			expectedOwner: "iamHehe",
			expectedRepo:  "laughts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.GithubConfig{
				Token:             "random-gibbrich-as-token",
				DefaultRepository: tt.repository,
			}
			client, _ := github.NewClient(context.Background(), cfg)

			owner, repo, err := client.GetDefaultRepository()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedOwner, owner)
				assert.Equal(t, tt.expectedRepo, repo)
			}
		})
	}
}
