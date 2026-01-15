package jenkins

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ignorant05/Uniflow/internal/config"
	"github.com/ignorant05/Uniflow/platforms/configurations/jenkins"
	"github.com/stretchr/testify/require"
)

// Setting up client with mock server
func SetupTestClientWithMockServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *jenkins.Client) {
	server := httptest.NewServer(handler)

	cfg := &config.JenkinsConfig{
		BaseURL:  server.URL,
		Username: "test-user",
		APIToken: "random-gibberish-token",
		JobName:  "test-job",
	}

	client, err := jenkins.NewTestClient(context.Background(), cfg)
	require.NoError(t, err)

	require.NotNil(t, client, "Client should not be nil")
	require.NotNil(t, client.Client, "Client.Client should not be nil")
	require.NotNil(t, client.Config, "Client.Config should not be nil")

	return server, client
}
