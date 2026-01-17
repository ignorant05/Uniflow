package gitlab_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/ignorant05/Uniflow/internal/config"
	errorhandling "github.com/ignorant05/Uniflow/internal/errorHandling"
	"github.com/ignorant05/Uniflow/platforms/configurations/gitlab"
	mock "github.com/ignorant05/Uniflow/platforms/tests/unit/gitlab"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gl "gitlab.com/gitlab-org/api/client-go"
)

// Testing ListPipelines func with valid setup
func TestListPipelines_Success(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipelines", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		response := []*gl.PipelineInfo{
			{
				ID:     123,
				Status: "success",
				Ref:    "main",
				SHA:    "abc123",
			},
			{
				ID:     124,
				Status: "failed",
				Ref:    "develop",
				SHA:    "def456",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	pipelines, err := client.ListPipelines("ignorant05", "Uniflow")

	require.NoError(t, err)
	assert.Equal(t, 2, len(pipelines))
	assert.Len(t, pipelines, 2)
	assert.Equal(t, "success", pipelines[0].Status)
	assert.Equal(t, "failed", pipelines[1].Status)
}

// Testing listing pipelines authentication failure (Bad credentials)
func TestListPipelines_AuthenticationFailure(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(map[string]string{
			"error": "401 Unauthorized",
		})
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	_, err := client.ListPipelines("ignorant05", "Uniflow")

	assert.Error(t, err)
}

// Testing ListPipelines rate limit
func TestListPipelines_RateLimit(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("RateLimit-Remaining", "0")
		w.Header().Set("RateLimit-ResetTime", "1234567890")
		w.WriteHeader(http.StatusTooManyRequests)
		err := json.NewEncoder(w).Encode(map[string]string{
			"error": "rate limit exceeded",
		})
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	_, err := client.ListPipelines("ignorant05", "Uniflow")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rate limit")
}

// Testing ListPipelines empty response
func TestListPipelines_Empty(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		response := []*gl.PipelineInfo{}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	pipelines, err := client.ListPipelines("ignorant05", "Uniflow")

	require.NoError(t, err)
	assert.Empty(t, pipelines)
}

// Testing ListPipelines (Failure: network error)
func TestListPipelines_NetworkError(t *testing.T) {
	cfg := &config.GitlabConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "ignorant05/Uniflow",
	}
	client, _ := gitlab.NewClient(context.Background(), cfg)
	badURL, _ := url.Parse("http://localhost:12345")

	// Replace the client with a bad URL client
	glClient, _ := gl.NewClient(cfg.Token, gl.WithBaseURL(badURL.String()+"/api/v4"))
	client.Client = glClient

	_, err := client.ListPipelines("ignorant05", "Uniflow")

	assert.Error(t, err)
}

// Testing ListPipelines with pagination
func TestListPipelines_WithPagination(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipelines", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		queryParams := r.URL.Query()
		assert.NotEmpty(t, queryParams.Get("per_page"))

		response := []*gl.PipelineInfo{
			{
				ID:     125,
				Status: "success",
				Ref:    "main",
				SHA:    "ghi789",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Page", "1")
		w.Header().Set("X-Per-Page", "20")
		w.Header().Set("X-Total", "1")
		w.Header().Set("X-Total-Pages", "1")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	pipelines, err := client.ListPipelines("ignorant05", "Uniflow")

	require.NoError(t, err)

	require.NotNil(t, pipelines)
	assert.NotEmpty(t, pipelines)

	assert.True(t, len(pipelines) > 0, "pipelines list should not be empty")
}

// Testing ListPipelines with various statuses
func TestListPipelines_VariousStatuses(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		response := []*gl.PipelineInfo{
			{
				ID:     126,
				Status: "success",
				Ref:    "main",
				SHA:    "abc123",
			},
			{
				ID:     127,
				Status: "failed",
				Ref:    "feature",
				SHA:    "def456",
			},
			{
				ID:     128,
				Status: "running",
				Ref:    "develop",
				SHA:    "ghi789",
			},
			{
				ID:     129,
				Status: "pending",
				Ref:    "staging",
				SHA:    "jkl012",
			},
			{
				ID:     130,
				Status: "canceled",
				Ref:    "hotfix",
				SHA:    "mno345",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	pipelines, err := client.ListPipelines("ignorant05", "Uniflow")

	require.NoError(t, err)
	assert.Equal(t, 5, len(pipelines))
	assert.Equal(t, "success", pipelines[0].Status)
	assert.Equal(t, "failed", pipelines[1].Status)
	assert.Equal(t, "running", pipelines[2].Status)
	assert.Equal(t, "pending", pipelines[3].Status)
	assert.Equal(t, "canceled", pipelines[4].Status)
}
