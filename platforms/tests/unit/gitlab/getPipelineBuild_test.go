package gitlab_test

import (
	"encoding/json"
	"net/http"
	"testing"

	errorhandling "github.com/ignorant05/Uniflow/internal/errorHandling"
	mock "github.com/ignorant05/Uniflow/platforms/tests/unit/gitlab"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gl "gitlab.com/gitlab-org/api/client-go"
)

// Testing GetPipeline, success
func TestGetPipeline_Success(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipelines/123456", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		response := gl.Pipeline{
			ID:     123456,
			Status: "success",
			Ref:    "main",
			SHA:    "abcdef1234567890",
			WebURL: "https://gitlab.com/ignorant05/Uniflow/-/pipelines/123456",
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	projectPath, _ := client.GetRepositoryPath()
	pipeline, _, err := client.Pipelines.GetPipeline(projectPath, 123456)

	require.NoError(t, err)
	assert.Equal(t, 123456, int(pipeline.ID))
	assert.Equal(t, "success", pipeline.Status)
	assert.Equal(t, "main", pipeline.Ref)
	assert.Equal(t, "abcdef1234567890", pipeline.SHA)
}

// Testing GetPipeline, running
func TestGetPipeline_Running(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipelines/123456", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		response := gl.Pipeline{
			ID:     123456,
			Status: "running",
			Ref:    "main",
			SHA:    "abcdef1234567890",
			WebURL: "https://gitlab.com/ignorant05/Uniflow/-/pipelines/123456",
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	projectPath, _ := client.GetRepositoryPath()
	pipeline, _, err := client.Pipelines.GetPipeline(projectPath, 123456)

	require.NoError(t, err)
	assert.Equal(t, 123456, int(pipeline.ID))
	assert.Equal(t, "running", pipeline.Status)
}

// Testing GetPipeline, pending
func TestGetPipeline_Pending(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipelines/123456", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		response := gl.Pipeline{
			ID:     123456,
			Status: "pending",
			Ref:    "develop",
			SHA:    "abcdef1234567890",
			WebURL: "https://gitlab.com/ignorant05/Uniflow/-/pipelines/123456",
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	projectPath, _ := client.GetRepositoryPath()
	pipeline, _, err := client.Pipelines.GetPipeline(projectPath, 123456)

	require.NoError(t, err)
	assert.Equal(t, 123456, int(pipeline.ID))
	assert.Equal(t, "pending", pipeline.Status)
}

// Testing GetPipeline, failed
func TestGetPipeline_Failed(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipelines/123456", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		response := gl.Pipeline{
			ID:     123456,
			Status: "failed",
			Ref:    "main",
			SHA:    "abcdef1234567890",
			WebURL: "https://gitlab.com/ignorant05/Uniflow/-/pipelines/123456",
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	projectPath, _ := client.GetRepositoryPath()
	pipeline, _, err := client.Pipelines.GetPipeline(projectPath, 123456)

	require.NoError(t, err)
	assert.Equal(t, 123456, int(pipeline.ID))
	assert.Equal(t, "failed", pipeline.Status)
}

// Testing GetPipeline, cancelled
func TestGetPipeline_Cancelled(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipelines/123456", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		response := gl.Pipeline{
			ID:     123456,
			Status: "canceled",
			Ref:    "main",
			SHA:    "abcdef1234567890",
			WebURL: "https://gitlab.com/ignorant05/Uniflow/-/pipelines/123456",
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	projectPath, _ := client.GetRepositoryPath()
	pipeline, _, err := client.Pipelines.GetPipeline(projectPath, 123456)

	require.NoError(t, err)
	assert.Equal(t, 123456, int(pipeline.ID))
	assert.Equal(t, "canceled", pipeline.Status)
}

// Testing GetPipeline, skipped
func TestGetPipeline_Skipped(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipelines/123456", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		response := gl.Pipeline{
			ID:     123456,
			Status: "skipped",
			Ref:    "main",
			SHA:    "abcdef1234567890",
			WebURL: "https://gitlab.com/ignorant05/Uniflow/-/pipelines/123456",
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	projectPath, _ := client.GetRepositoryPath()
	pipeline, _, err := client.Pipelines.GetPipeline(projectPath, 123456)

	require.NoError(t, err)
	assert.Equal(t, 123456, int(pipeline.ID))
	assert.Equal(t, "skipped", pipeline.Status)
}
