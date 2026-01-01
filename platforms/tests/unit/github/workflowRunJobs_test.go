package github_test

import (
	"encoding/json"
	"net/http"
	"testing"

	gh "github.com/google/go-github/v57/github"
	errorhandling "github.com/ignorant05/Uniflow/internal/errorHandling"
	mock "github.com/ignorant05/Uniflow/platforms/tests/unit/github"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Testing ListWorkflowJobs, success
func TestListWorkflowRunJobs_Success(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repos/ignorant05/Uniflow/actions/runs/123456/jobs", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		response := gh.Jobs{
			TotalCount: gh.Int(2),
			Jobs: []*gh.WorkflowJob{
				{
					ID:         gh.Int64(123456),
					Name:       gh.String("build"),
					Status:     gh.String("completed"),
					Conclusion: gh.String("success"),
					RunID:      gh.Int64(123456),
				},
				{
					ID:     gh.Int64(123457),
					Name:   gh.String("test"),
					Status: gh.String("in_progress"),
					RunID:  gh.Int64(123456),
				},
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
	jobs, err := client.ListWorkflowJobs(owner, repo, 123456)

	require.NoError(t, err)
	assert.Equal(t, 2, len(jobs))
	assert.Equal(t, int64(123456), jobs[0].GetID())
	assert.Equal(t, "build", jobs[0].GetName())
	assert.Equal(t, "completed", jobs[0].GetStatus())
	assert.Equal(t, "success", jobs[0].GetConclusion())
	assert.Equal(t, int64(123457), jobs[1].GetID())
	assert.Equal(t, "test", jobs[1].GetName())
	assert.Equal(t, "in_progress", jobs[1].GetStatus())
}

// Testing ListWorkflowJobs, (Failure: not found)
func TestListWorkflowRunJobs_Failure(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repos/ignorant05/Uniflow/actions/runs/123456/jobs", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		w.WriteHeader(http.StatusNotFound)

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "no job found by id",
		})

		if err != nil {
			errorhandling.HandleError(err)
		}
	})

	defer server.Close()

	owner, repo, _ := client.GetDefaultRepository()
	_, err := client.ListWorkflowJobs(owner, repo, 123456)

	assert.Error(t, err)
}
