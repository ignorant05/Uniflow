package github_test

import (
	"encoding/json"
	"net/http"
	"testing"

	mock "github.com/ignorant05/Uniflow/platforms/tests/unit/github"

	gh "github.com/google/go-github/v57/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Testing GetWorkflowRun, success
func TestGetWorkflowRun_Success(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repos/ignorant05/Uniflow/actions/runs/123456", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		response := gh.WorkflowRun{
			ID:         gh.Int64(123456),
			Status:     gh.String("completed"),
			Conclusion: gh.String("success"),
			HTMLURL:    gh.String("https://github.com/owner/repo/actions/runs/123456"),
			HeadBranch: gh.String("main"),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	defer server.Close()

	owner, repo, _ := client.GetDefaultRepository()
	run, _, err := client.Client.Actions.GetWorkflowRunByID(
		client.Ctx,
		owner,
		repo,
		123456,
	)

	require.NoError(t, err)
	assert.Equal(t, int64(123456), *run.ID)
	assert.Equal(t, "completed", *run.Status)
	assert.Equal(t, "success", *run.Conclusion)
}

// Testing GetWorkflowRun, in progress
func TestGetWorkflowRun_InProgress(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repos/ignorant05/Uniflow/actions/runs/123456", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		response := gh.WorkflowRun{
			ID:         gh.Int64(123456),
			Status:     gh.String("in_progress"),
			Conclusion: nil,
			HTMLURL:    gh.String("https://github.com/ignorant05/Uniflow/actions/runs/123456"),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	defer server.Close()

	owner, repo, _ := client.GetDefaultRepository()
	run, _, err := client.Client.Actions.GetWorkflowRunByID(
		client.Ctx,
		owner,
		repo,
		123456,
	)

	require.NoError(t, err)
	assert.Equal(t, int64(123456), *run.ID)
	assert.Equal(t, "in_progress", *run.Status)
	assert.Nil(t, run.Conclusion)
}

// Testing GetWorkflowRun, failure
func TestGetWorkflowRun_Failure(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repos/ignorant05/Uniflow/actions/runs/123456", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		response := gh.WorkflowRun{
			ID:      gh.Int64(123456),
			Status:  gh.String("failure"),
			HTMLURL: gh.String("https://github.com/ignorant05/Uniflow/actions/runs/123456"),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	defer server.Close()

	owner, repo, _ := client.GetDefaultRepository()
	run, _, err := client.Client.Actions.GetWorkflowRunByID(
		client.Ctx,
		owner,
		repo,
		123456,
	)

	require.NoError(t, err)
	assert.Equal(t, "failure", *run.Status)
}
