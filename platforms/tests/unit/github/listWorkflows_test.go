package github_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	gh "github.com/google/go-github/v57/github"
	"github.com/ignorant05/Uniflow/internal/config"
	"github.com/ignorant05/Uniflow/platforms/github"
	mock "github.com/ignorant05/Uniflow/platforms/tests/unit/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Testing ListWorkflows func with valid setup
func TestListWorkflows_Success(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repos/ignorant05/Uniflow/actions/workflows", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		response := gh.Workflows{
			TotalCount: gh.Int(2),
			Workflows: []*gh.Workflow{
				{
					ID:   gh.Int64(1),
					Name: gh.String("CI"),
					Path: gh.String(".github/workflows/ci.yaml"),
				},
				{
					ID:   gh.Int64(2),
					Name: gh.String("Deploy"),
					Path: gh.String(".github/workflows/deploy.yaml"),
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	})

	defer server.Close()

	workflows, err := client.ListWorkflows("ignorant05", "Uniflow")

	require.NoError(t, err)
	assert.Equal(t, 2, len(workflows))
	assert.Len(t, workflows, 2)
	assert.Equal(t, "CI", workflows[0].GetName())
	assert.Equal(t, "Deploy", workflows[1].GetName())
}

// Testing listing workflows authentication failure (Bad credentials (no credentials on server))
func TestListWorkflows_AuthenticationFailure(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid credentials",
		})

	})

	defer server.Close()

	_, err := client.ListWorkflows("ignorant05", "Uniflow")

	assert.Error(t, err)
}

// Testing ListWorkflows rate limit
func TestListWorkflowsRateLimit(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Remaining", "0")
		w.Header().Set("X-RateLimit-Reset", "1234567890")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "api rate limit exceeded",
		})
	})

	defer server.Close()

	_, err := client.ListWorkflows("ignorant05", "Uniflow")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rate limit")
}

// Testing ListWorkflows empty response
func TestListWorkflows_Empty(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		response := gh.Workflows{
			TotalCount: gh.Int(0),
			Workflows:  []*gh.Workflow{},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	defer server.Close()

	workflows, err := client.ListWorkflows("ignorant05", "Uniflow")

	require.NoError(t, err)
	assert.Empty(t, workflows)
}

// Testing ListWorkflows (Failure: network error)
func TestListWorkflows_NetworkError(t *testing.T) {
	cfg := &config.GithubConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "ignorant05/Uniflow",
	}

	client, _ := github.NewClient(context.Background(), cfg)

	badURL, _ := url.Parse("http://localhost:12345")
	client.Client.BaseURL = badURL

	_, err := client.ListWorkflows("ignorant05", "Uniflow")

	assert.Error(t, err)
}
