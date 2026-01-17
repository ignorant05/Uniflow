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

// Testing GetPipelineJobs, success
func TestGetPipelineJobs_Success(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipelines/123456/jobs", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		response := []*gl.Job{
			{
				ID:     123456,
				Name:   "build",
				Status: "success",
				Stage:  "build",
			},
			{
				ID:     123457,
				Name:   "test",
				Status: "running",
				Stage:  "test",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	owner, repo, _ := client.GetDefaultRepository()
	jobs, err := client.GetPipelineJobs(owner, repo, 123456)

	require.NoError(t, err)
	assert.Equal(t, 2, len(jobs))
	assert.Equal(t, 123456, int(jobs[0].ID))
	assert.Equal(t, "build", jobs[0].Name)
	assert.Equal(t, "success", jobs[0].Status)
	assert.Equal(t, "build", jobs[0].Stage)
	assert.Equal(t, 123457, int(jobs[1].ID))
	assert.Equal(t, "test", jobs[1].Name)
	assert.Equal(t, "running", jobs[1].Status)
}

// Testing GetPipelineJobs, empty response
func TestGetPipelineJobs_Empty(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipelines/123456/jobs", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		response := []*gl.Job{}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	owner, repo, _ := client.GetDefaultRepository()
	jobs, err := client.GetPipelineJobs(owner, repo, 123456)

	require.NoError(t, err)
	assert.Empty(t, jobs)
}

// Testing GetPipelineJobs, (Failure: pipeline not found)
func TestGetPipelineJobs_NotFound(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipelines/123456/jobs", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "404 Pipeline not found",
		})
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	owner, repo, _ := client.GetDefaultRepository()
	_, err := client.GetPipelineJobs(owner, repo, 123456)

	assert.Error(t, err)
}

// Testing GetPipelineJobs, (Failure: authentication error)
func TestGetPipelineJobs_AuthenticationFailure(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipelines/123456/jobs", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "401 Unauthorized",
		})
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	owner, repo, _ := client.GetDefaultRepository()
	_, err := client.GetPipelineJobs(owner, repo, 123456)

	assert.Error(t, err)
}

// Testing GetPipelineJobs with various job statuses
func TestGetPipelineJobs_VariousStatuses(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipelines/123456/jobs", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		response := []*gl.Job{
			{
				ID:     123456,
				Name:   "build",
				Status: "success",
				Stage:  "build",
			},
			{
				ID:     123457,
				Name:   "test",
				Status: "failed",
				Stage:  "test",
			},
			{
				ID:     123458,
				Name:   "deploy",
				Status: "running",
				Stage:  "deploy",
			},
			{
				ID:     123459,
				Name:   "notify",
				Status: "pending",
				Stage:  "notify",
			},
			{
				ID:     123460,
				Name:   "cleanup",
				Status: "canceled",
				Stage:  "cleanup",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	owner, repo, _ := client.GetDefaultRepository()
	jobs, err := client.GetPipelineJobs(owner, repo, 123456)

	require.NoError(t, err)
	assert.Equal(t, 5, len(jobs))
	assert.Equal(t, "success", jobs[0].Status)
	assert.Equal(t, "failed", jobs[1].Status)
	assert.Equal(t, "running", jobs[2].Status)
	assert.Equal(t, "pending", jobs[3].Status)
	assert.Equal(t, "canceled", jobs[4].Status)
}

// Testing GetPipelineJobs with pagination
func TestGetPipelineJobs_WithPagination(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipelines/123456/jobs", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		queryParams := r.URL.Query()
		assert.NotEmpty(t, queryParams.Get("per_page"))

		response := []*gl.Job{
			{
				ID:     123461,
				Name:   "page1-job",
				Status: "success",
				Stage:  "build",
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

	owner, repo, _ := client.GetDefaultRepository()
	jobs, err := client.GetPipelineJobs(owner, repo, 123456)

	require.NoError(t, err)
	assert.NotEmpty(t, jobs)
	assert.Equal(t, "page1-job", jobs[0].Name)
}

// Testing GetPipelineJobs, (Failure: rate limit)
func TestGetPipelineJobs_RateLimit(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipelines/123456/jobs", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("RateLimit-Remaining", "0")
		w.Header().Set("RateLimit-ResetTime", "1234567890")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "rate limit exceeded",
		})
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	owner, repo, _ := client.GetDefaultRepository()
	_, err := client.GetPipelineJobs(owner, repo, 123456)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rate limit")
}
