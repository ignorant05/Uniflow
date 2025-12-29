package github_test

import (
	"context"
	"os"
	"testing"

	"github.com/ignorant05/Uniflow/internal/config"
	"github.com/ignorant05/Uniflow/platforms/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Testing client creation with a valid token
func TestClientWithToken(t *testing.T) {
	cfg := &config.GithubConfig{
		Token: "random-gibbrich-as-token",
	}

	client, err := github.NewClient(context.Background(), cfg)

	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.Client)
}

// Testing client creation without (with invalid) token
func TestClientWithoutToken(t *testing.T) {
	// Remove the default token (if there is in the your env)
	token := os.Getenv("GITHUB_TOKEN")

	defer os.Setenv("GITHUB_TOKEN", token)
	os.Unsetenv("GITHUB_TOKEN")

	cfg := &config.GithubConfig{
		Token: "",
	}

	_, err := github.NewClient(context.Background(), cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "GITHUB_TOKEN")
}

// Testing client creation with token from env (as env variable)
// Make sure to set an env variable called "GITHUB_TOKEN" for this one before launching it or i'll fail
// use export GITHUB_TOKEN="ur token here"
func TestClientWithTokenFromEnv(t *testing.T) {
	// Uncomment the next line and set ur token here for testing and place it in a replacements.txt (if you want to re-push this file again)
	// os.Setenv("GITHUB_TOKEN", "your token here")

	token := os.Getenv("GITHUB_TOKEN")

	cfg := &config.GithubConfig{
		Token: token,
	}

	client, err := github.NewClient(context.Background(), cfg)

	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.Client)
}

// Testing client creation with custom enterpise URL
func TestClientWithEnterpriseURL(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	baseURL := "https://github.enterprise.com"

	cfg := &config.GithubConfig{
		Token:   token,
		BaseURL: baseURL,
	}

	client, err := github.NewClient(context.Background(), cfg)

	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.Client)
}

// Testing client creation from profile
func TestClientFromProfile_Success(t *testing.T) {
	profile := &config.Profile{
		Github: &config.GithubConfig{
			Token: "random-gibbrich-as-token",
		},
	}

	client, err := github.NewClientFromProfile(context.Background(), profile)

	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.Client)
}

// Testing client creatino from profil with invalid github field
func TestClientFromProfile_Failure(t *testing.T) {
	profile := &config.Profile{
		Github: nil,
	}

	_, err := github.NewClientFromProfile(context.Background(), profile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Github isn't configured for this profile")
}
