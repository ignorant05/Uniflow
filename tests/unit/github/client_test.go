package github_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	gh "github.com/google/go-github/v57/github"

	"github.com/ignorant05/Uniflow/configs/github"
	"github.com/ignorant05/Uniflow/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Setting up client with mock server
func setupTestClientWithMockServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *github.Client) {
	server := httptest.NewServer(handler)

	cfg := &config.GithubConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "ignorant05/Uniflow",
	}

	client, err := github.NewClient(context.Background(), cfg)

	require.NoError(t, err)
	baseURL, _ := url.Parse(server.URL + "/")
	client.Client.BaseURL = baseURL
	client.Client.UploadURL = baseURL

	return server, client
}

// Testing client creation with a valid token
func TestClientWithToken(t *testing.T) {
	cfg := &config.GithubConfig{
		Token: "random-gibbrich-as-token",
	}

	client, err := github.NewClient(context.Background(), cfg)

	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.Client)
}

// Testing client creation without (with invalid) token
func TestClientWithoutToken(t *testing.T) {
	// Remove the default token (if there is in the your env)
	token := os.Getenv("GITHUB_TOKEN")

	defer os.Setenv("GITHUB_TOKEN", token)
	os.Unsetenv("GITHUB_TOKEN")

	cfg := &config.GithubConfig{
		Token: "",
	}

	_, err := github.NewClient(context.Background(), cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "GITHUB_TOKEN")
}

// Testing client creation with token from env (as env variable)
// Make sure to set an env variable called "GITHUB_TOKEN" for this one before launching it or i'll fail
// use export GITHUB_TOKEN="ur token here"
func TestClientWithTokenFromEnv(t *testing.T) {
	os.Setenv("GITHUB_TOKEN", "")

	token := os.Getenv("GITHUB_TOKEN")

	cfg := &config.GithubConfig{
		Token: token,
	}

	client, err := github.NewClient(context.Background(), cfg)

	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.Client)
}

// Testing client creation with custom enterpise URL
func TestClientWithEnterpriseURL(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	baseURL := "https://github.enterprise.com"

	cfg := &config.GithubConfig{
		Token:   token,
		BaseURL: baseURL,
	}

	client, err := github.NewClient(context.Background(), cfg)

	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.Client)
}

// Testing client creation from profile
func TestClientFromProfile_Success(t *testing.T) {
	profile := &config.Profile{
		Github: &config.GithubConfig{
			Token: "random-gibbrich-as-token",
		},
	}

	client, err := github.NewClientFromProfile(context.Background(), profile)

	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.Client)
}

// Testing client creatino from profil with invalid github field
func TestClientFromProfile_Failure(t *testing.T) {
	profile := &config.Profile{
		Github: nil,
	}

	_, err := github.NewClientFromProfile(context.Background(), profile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Github isn't configured for this profile")
}

// Testing default repo configured
func TestDefaultRepo_Success(t *testing.T) {
	cfg := &config.GithubConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "ignorant05/Uniflow",
	}

	client, _ := github.NewClient(context.Background(), cfg)
	owner, repo, err := client.GetDefaultRepository()

	require.NoError(t, err)
	assert.Equal(t, owner, "ignorant05")
	assert.Equal(t, repo, "Uniflow")
}

// Testing default repo not configured
func TestDefaultRepo_Failure(t *testing.T) {
	cfg := &config.GithubConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "",
	}

	client, _ := github.NewClient(context.Background(), cfg)
	owner, repo, err := client.GetDefaultRepository()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "No default repository configured")
	assert.Empty(t, owner)
	assert.Empty(t, repo)
}

// Testing ListWorkflows func with valid setup
func TestListWorkflows_Success(t *testing.T) {
	server, client := setupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
	server, client := setupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
	server, client := setupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
	server, client := setupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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

// Testing getting default repository
func TestGetDefaultRepository_Success(t *testing.T) {
	cfg := &config.GithubConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "ignorant05/Uniflow",
	}

	client, _ := github.NewClient(context.Background(), cfg)

	owner, repo, err := client.GetDefaultRepository()

	require.NoError(t, err)
	assert.Equal(t, "ignorant05", owner)
	assert.Equal(t, "Uniflow", repo)
}

// Testing getting default repository (Failure: not configured)
func TestGetDefaultRepository_Failure(t *testing.T) {
	cfg := &config.GithubConfig{
		Token:             "random-gibbrich-as-token",
		DefaultRepository: "",
	}

	client, _ := github.NewClient(context.Background(), cfg)

	_, _, err := client.GetDefaultRepository()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "No default repository")
}

// Testing default repo, table driven
func TestGetDefaultRepository_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		repository    string
		expectError   bool
		expectedOwner string
		expectedRepo  string
	}{
		{
			name:          "lols",
			repository:    "leader/my-lols",
			expectError:   false,
			expectedOwner: "leader",
			expectedRepo:  "my-lols",
		},
		{
			name:        "empty repo",
			repository:  "",
			expectError: true,
		},
		{
			name:          "org of ones",
			repository:    "one/numbers",
			expectError:   false,
			expectedOwner: "one",
			expectedRepo:  "numbers",
		},
		{
			name:          "hehe",
			repository:    "iamHehe/laughts",
			expectError:   false,
			expectedOwner: "iamHehe",
			expectedRepo:  "laughts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.GithubConfig{
				Token:             "random-gibbrich-as-token",
				DefaultRepository: tt.repository,
			}
			client, _ := github.NewClient(context.Background(), cfg)

			owner, repo, err := client.GetDefaultRepository()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedOwner, owner)
				assert.Equal(t, tt.expectedRepo, repo)
			}
		})
	}
}

// Testing triggerWorkflow, success
func TestTriggerWorkflow_Success(t *testing.T) {
	server, client := setupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repos/ignorant05/Uniflow/actions/workflows/ci.yml/dispatches", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "main", body["ref"])

		w.WriteHeader(http.StatusNoContent)
	})

	defer server.Close()

	owner, repo, _ := client.GetDefaultRepository()
	event := gh.CreateWorkflowDispatchEventRequest{
		Ref: "main",
	}

	_, err := client.Actions.CreateWorkflowDispatchEventByFileName(
		client.Ctx,
		owner,
		repo,
		"ci.yml",
		event,
	)

	assert.NoError(t, err)
}

// Testing triggerWorkflow with inputs, success
func TestTriggerWorkflowWithInputs_Success(t *testing.T) {
	server, client := setupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repos/ignorant05/Uniflow/actions/workflows/ci.yml/dispatches", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "main", body["ref"])

		inputs := body["inputs"].(map[string]interface{})
		assert.Equal(t, "v1.2.3", inputs["version"])
		assert.Equal(t, "production", inputs["environment"])

		w.WriteHeader(http.StatusNoContent)
	})

	defer server.Close()

	inputs := map[string]interface{}{
		"version":     "v1.2.3",
		"environment": "production",
	}

	owner, repo, _ := client.GetDefaultRepository()
	event := gh.CreateWorkflowDispatchEventRequest{
		Ref:    "main",
		Inputs: inputs,
	}

	_, err := client.Actions.CreateWorkflowDispatchEventByFileName(
		client.Ctx,
		owner,
		repo,
		"ci.yml",
		event,
	)

	assert.NoError(t, err)
}

// Testing triggerWorkflow, (Failure: non existent workflow file)
func TestTriggerWorkflow_NonExistentFile(t *testing.T) {
	server, client := setupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "main", body["ref"])

		w.WriteHeader(http.StatusNoContent)
	})

	defer server.Close()

	owner, repo, _ := client.GetDefaultRepository()
	event := gh.CreateWorkflowDispatchEventRequest{
		Ref: "main",
	}

	_, err := client.Actions.CreateWorkflowDispatchEventByFileName(
		client.Ctx,
		owner,
		repo,
		"nonexistent.yml",
		event,
	)

	assert.NoError(t, err)
}

// Testing workflow, (Failure: no content)
func TestTriggerWorkflow_Failure(t *testing.T) {
	server, client := setupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repos/ignorant05/Uniflow/actions/workflows/ci.yml/dispatches", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "main", body["ref"])

		w.WriteHeader(http.StatusNoContent)
	})

	defer server.Close()

	owner, repo, _ := client.GetDefaultRepository()
	event := gh.CreateWorkflowDispatchEventRequest{
		Ref: "main",
	}

	_, err := client.Actions.CreateWorkflowDispatchEventByFileName(
		client.Ctx,
		owner,
		repo,
		"ci.yml",
		event,
	)

	assert.NoError(t, err)
}

// Testing GetWorkflowRun, success
func TestGetWorkflowRun_Success(t *testing.T) {
	server, client := setupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
	server, client := setupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
	server, client := setupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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

// Testing ListWorkflowJobs, success
func TestListWorkflowRunJobs_Success(t *testing.T) {
	server, client := setupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
		json.NewEncoder(w).Encode(response)
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
	server, client := setupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/repos/ignorant05/Uniflow/actions/runs/123456/jobs", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		w.WriteHeader(http.StatusNotFound)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "no job found by id",
		})
	})

	defer server.Close()

	owner, repo, _ := client.GetDefaultRepository()
	_, err := client.ListWorkflowJobs(owner, repo, 123456)

	assert.Error(t, err)
}
