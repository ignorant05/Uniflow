package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ignorant05/Uniflow/internal/config"
	"github.com/ignorant05/Uniflow/platforms/github"
	"github.com/stretchr/testify/require"
)

// Setting up client with mock server
func SetupTestClientWithMockServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *github.Client) {
	server := httptest.NewServer(handler)

	cfg := &config.GithubConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "ignorant05/Uniflow",
	}

	client, err := github.NewClient(context.Background(), cfg)

	require.NoError(t, err)
	baseURL, _ := url.Parse(server.URL + "/")
	client.BaseURL = baseURL
	client.UploadURL = baseURL

	return server, client
}
