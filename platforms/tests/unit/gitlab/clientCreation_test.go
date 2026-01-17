package gitlab_test

import (
	"context"
	"testing"

	"github.com/ignorant05/Uniflow/internal/config"
	"github.com/ignorant05/Uniflow/platforms/configurations/gitlab"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Testing client creation with a valid token
func TestClientWithToken(t *testing.T) {
	cfg := &config.GitlabConfig{
		Token: "random-gibbrich-as-token",
	}
	client, err := gitlab.NewClient(context.Background(), cfg)
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.Client)
}

// // Testing client creation with token from env (as env variable)
// // Make sure to set an env variable called "GITLAB_TOKEN" for this one before launching it or it will fail
// // use export GITLAB_TOKEN="your token here"
// func TestClientWithTokenFromEnv(t *testing.T) {
// 	// Uncomment the next line and set your token here for testing and place it in a replacements.txt (if you want to re-push this file again)
// 	// os.Setenv("GITLAB_TOKEN", "your token here")
//
// 	token := os.Getenv("GITLAB_TOKEN")
//
// 	cfg := &config.GitlabConfig{
// 		Token: token,
// 	}
//
// 	client, err := gitlab.NewClient(context.Background(), cfg)
//
// 	require.NoError(t, err)
// 	assert.NotNil(t, client)
// 	assert.NotNil(t, client.Client)
// }
//
// // Testing client creation with custom GitLab URL
// func TestClientWithCustomURL(t *testing.T) {
// 	token := os.Getenv("GITLAB_TOKEN")
// 	baseURL := "https://gitlab.enterprise.com"
//
// 	cfg := &config.GitlabConfig{
// 		Token:   token,
// 		BaseURL: baseURL,
// 	}
//
// 	client, err := gitlab.NewClient(context.Background(), cfg)
//
// 	require.NoError(t, err)
// 	assert.NotNil(t, client)
// 	assert.NotNil(t, client.Client)
// }

// Testing client creation from profile
func TestClientFromProfile_Success(t *testing.T) {
	profile := &config.Profile{
		Gitlab: &config.GitlabConfig{
			Token: "random-gibbrich-as-token",
		},
	}
	client, err := gitlab.NewClientFromProfile(context.Background(), profile)
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.Client)
}

// Testing client creation from profile with invalid gitlab field
func TestClientFromProfile_Failure(t *testing.T) {
	profile := &config.Profile{
		Gitlab: nil,
	}
	_, err := gitlab.NewClientFromProfile(context.Background(), profile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "GitLab isn't configured for this profile")
}

// Testing client creation with nil configuration
func TestClientWithNilConfig(t *testing.T) {
	_, err := gitlab.NewClient(context.Background(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "configuration cannot be nil")
}

// // Testing client creation with missing token
// func TestClientWithMissingToken(t *testing.T) {
// 	// NOTE:
// 	// The following lines are just a workaround for invalid token test
// 	// i didn't find another way, if you do, please open an issue
//
// 	// Retrieving the original token (if there is)
// 	originalToken := os.Getenv("GITLAB_TOKEN")
//
// 	// Unset the env var to simulate missing token
// 	err := os.Unsetenv("GITLAB_TOKEN")
// 	if err != nil {
// 		return
// 	}
//
// 	// Restore it after test completes
// 	t.Cleanup(func() {
// 		if originalToken != "" {
// 			err = os.Setenv("GITLAB_TOKEN", originalToken)
// 			if err != nil {
// 				return
// 			}
// 		}
// 	})
//
// 	cfg := &config.GitlabConfig{
// 		Token: "",
// 	}
// 	_, err = gitlab.NewClient(context.Background(), cfg)
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "no GitLab token found")
// }

// Testing GetDefaultRepository
func TestGetDefaultRepository_Success(t *testing.T) {
	cfg := &config.GitlabConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "ignorant05/Uniflow",
	}
	client, err := gitlab.NewClient(context.Background(), cfg)
	require.NoError(t, err)

	owner, repo, err := client.GetDefaultRepository()
	require.NoError(t, err)
	assert.Equal(t, "ignorant05", owner)
	assert.Equal(t, "Uniflow", repo)
}

// Testing GetDefaultRepository with no default repository configured
func TestGetDefaultRepository_Failure(t *testing.T) {
	cfg := &config.GitlabConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "",
	}
	client, err := gitlab.NewClient(context.Background(), cfg)
	require.NoError(t, err)

	_, _, err = client.GetDefaultRepository()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no default repository configured")
}

// Testing GetRepositoryPath
func TestGetRepositoryPath_Success(t *testing.T) {
	cfg := &config.GitlabConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "ignorant05/Uniflow",
	}
	client, err := gitlab.NewClient(context.Background(), cfg)
	require.NoError(t, err)

	path, err := client.GetRepositoryPath()
	require.NoError(t, err)
	assert.Equal(t, "ignorant05/Uniflow", path)
}

// Testing GetRepositoryPath with no default repository configured
func TestGetRepositoryPath_Failure(t *testing.T) {
	cfg := &config.GitlabConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "",
	}
	client, err := gitlab.NewClient(context.Background(), cfg)
	require.NoError(t, err)

	_, err = client.GetRepositoryPath()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no default repository configured")
}
