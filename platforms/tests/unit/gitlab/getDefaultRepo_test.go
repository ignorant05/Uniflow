package gitlab_test

import (
	"context"
	"testing"

	"github.com/ignorant05/Uniflow/internal/config"
	"github.com/ignorant05/Uniflow/platforms/configurations/gitlab"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Testing default repo configured
func TestDefaultRepo_Success(t *testing.T) {
	cfg := &config.GitlabConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "ignorant05/Uniflow",
	}
	client, _ := gitlab.NewClient(context.Background(), cfg)
	owner, repo, err := client.GetDefaultRepository()
	require.NoError(t, err)
	assert.Equal(t, owner, "ignorant05")
	assert.Equal(t, repo, "Uniflow")
}

// Testing default repo not configured
func TestDefaultRepo_Failure(t *testing.T) {
	cfg := &config.GitlabConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "",
	}
	client, _ := gitlab.NewClient(context.Background(), cfg)
	owner, repo, err := client.GetDefaultRepository()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no default repository configured")
	assert.Empty(t, owner)
	assert.Empty(t, repo)
}

// Testing getting default repository
func TestGetDefaultRepositorySuccess(t *testing.T) {
	cfg := &config.GitlabConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "ignorant05/Uniflow",
	}
	client, _ := gitlab.NewClient(context.Background(), cfg)
	owner, repo, err := client.GetDefaultRepository()
	require.NoError(t, err)
	assert.Equal(t, "ignorant05", owner)
	assert.Equal(t, "Uniflow", repo)
}

// Testing getting default repository (Failure: not configured)
func TestGetDefaultRepositoryFailure(t *testing.T) {
	cfg := &config.GitlabConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "",
	}
	client, _ := gitlab.NewClient(context.Background(), cfg)
	_, _, err := client.GetDefaultRepository()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no default repository")
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
			name:          "simple group/project",
			repository:    "leader/my-project",
			expectError:   false,
			expectedOwner: "leader",
			expectedRepo:  "my-project",
		},
		{
			name:        "empty repo",
			repository:  "",
			expectError: true,
		},
		{
			name:          "numeric group name",
			repository:    "one/numbers",
			expectError:   false,
			expectedOwner: "one",
			expectedRepo:  "numbers",
		},
		{
			name:          "mixed case names",
			repository:    "iamHehe/laughts",
			expectError:   false,
			expectedOwner: "iamHehe",
			expectedRepo:  "laughts",
		},
		{
			name:          "hyphenated project name",
			repository:    "my-group/my-project",
			expectError:   false,
			expectedOwner: "my-group",
			expectedRepo:  "my-project",
		},
		{
			name:          "underscore in names",
			repository:    "my_group/my_project",
			expectError:   false,
			expectedOwner: "my_group",
			expectedRepo:  "my_project",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.GitlabConfig{
				Token:             "random-gibbrich-as-token",
				DefaultRepository: tt.repository,
			}
			client, _ := gitlab.NewClient(context.Background(), cfg)
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

// Testing GetRepositoryPath
func TestGetRepositoryPathSuccess(t *testing.T) {
	cfg := &config.GitlabConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "ignorant05/Uniflow",
	}
	client, _ := gitlab.NewClient(context.Background(), cfg)
	path, err := client.GetRepositoryPath()
	require.NoError(t, err)
	assert.Equal(t, "ignorant05/Uniflow", path)
}

// Testing GetRepositoryPath failure
func TestGetRepositoryPathFailure(t *testing.T) {
	cfg := &config.GitlabConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "",
	}
	client, _ := gitlab.NewClient(context.Background(), cfg)
	path, err := client.GetRepositoryPath()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no default repository configured")
	assert.Empty(t, path)
}
