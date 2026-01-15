package jenkins_test

import (
	"context"
	"testing"

	"github.com/ignorant05/Uniflow/internal/config"
	"github.com/ignorant05/Uniflow/platforms/configurations/jenkins"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Testing client creation with username + API token
func TestClientWithAPIToken(t *testing.T) {
	cfg := &config.JenkinsConfig{
		BaseURL:  "https://jenkins.example.com",
		Username: "test-user",
		APIToken: "random-gibberish-token",
	}

	// Use NewTestClient instead of NewClient
	client, err := jenkins.NewTestClient(context.Background(), cfg)

	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.Client)
}

// Testing client creation with password authentication
func TestClientWithPassword(t *testing.T) {
	cfg := &config.JenkinsConfig{
		BaseURL:  "https://jenkins.example.com",
		Username: "test-user",
		Password: "test-password",
	}

	client, err := jenkins.NewTestClient(context.Background(), cfg)

	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.Client)
}

// Testing client creation with insecure TLS (self-signed certs)
func TestClientWithInsecureTLS(t *testing.T) {
	cfg := &config.JenkinsConfig{
		BaseURL:            "https://jenkins.example.com",
		Username:           "test-user",
		APIToken:           "random-gibberish-token",
		InsecureSkipVerify: true,
		TimeoutSeconds:     10,
	}

	client, err := jenkins.NewTestClient(context.Background(), cfg)

	require.NoError(t, err)
	assert.NotNil(t, client)
}

// Testing client creation with custom CA certificate
func TestClientWithCustomCACert(t *testing.T) {
	cfg := &config.JenkinsConfig{
		BaseURL:    "https://jenkins.example.com",
		Username:   "test-user",
		APIToken:   "random-gibberish-token",
		CACertPath: "/path/to/ca.pem",
	}

	client, err := jenkins.NewTestClient(context.Background(), cfg)

	require.NoError(t, err)
	assert.NotNil(t, client)
}

// Testing failure when BaseURL is missing
func TestClientFailure_MissingBaseURL(t *testing.T) {
	cfg := &config.JenkinsConfig{
		Username: "test-user",
		APIToken: "token",
	}

	_, err := jenkins.NewTestClient(context.Background(), cfg)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "base_url")
}

// Testing failure when no authentication method is provided
func TestClientFailure_MissingAuth(t *testing.T) {
	cfg := &config.JenkinsConfig{
		BaseURL: "https://jenkins.example.com",
	}

	_, err := jenkins.NewTestClient(context.Background(), cfg)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "username is required")
}
