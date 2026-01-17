package gitlab

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ignorant05/Uniflow/internal/config"
	"github.com/ignorant05/Uniflow/platforms/configurations/gitlab"
	"github.com/stretchr/testify/require"
	glClient "gitlab.com/gitlab-org/api/client-go"
)

// Setting up client with mock server
func SetupTestClientWithMockServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *gitlab.Client) {
	server := httptest.NewServer(handler)
	cfg := &config.GitlabConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "ignorant05/Uniflow",
	}
	client, err := gitlab.NewClient(context.Background(), cfg)
	require.NoError(t, err)

	// Create a new GitLab client with the mock server URL
	mockClient, err := glClient.NewClient(cfg.Token, glClient.WithBaseURL(server.URL+"/api/v4"))
	require.NoError(t, err)

	// Replace the client with the mock client
	client.Client = mockClient

	return server, client
}
