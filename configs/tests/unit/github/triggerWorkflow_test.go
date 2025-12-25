package github_test

import (
	"encoding/json"
	"net/http"
	"testing"

	gh "github.com/google/go-github/v57/github"
	mock "github.com/ignorant05/Uniflow/configs/tests/unit/github"

	"github.com/stretchr/testify/assert"
)

// Testing triggerWorkflow, success
func TestTriggerWorkflow_Success(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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
