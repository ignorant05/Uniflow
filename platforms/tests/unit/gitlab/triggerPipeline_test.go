package gitlab_test

import (
	"encoding/json"
	"net/http"
	"testing"

	errorhandling "github.com/ignorant05/Uniflow/internal/errorHandling"
	mock "github.com/ignorant05/Uniflow/platforms/tests/unit/gitlab"
	"github.com/stretchr/testify/assert"
	gl "gitlab.com/gitlab-org/api/client-go"
)

// Testing TriggerPipeline, success
func TestTriggerPipeline_Success(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipeline", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			errorhandling.HandleError(err)
		}
		assert.Equal(t, "\"main\"", body["ref"])

		response := gl.Pipeline{
			ID:     131,
			Status: "pending",
			Ref:    "main",
			SHA:    "abc123",
			WebURL: "https://gitlab.com/ignorant05/Uniflow/-/pipelines/131",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	owner, repo, _ := client.GetDefaultRepository()
	err := client.TriggerPipeline(owner, repo, "main", nil)

	assert.NoError(t, err)
}

// Testing TriggerPipeline with variables, success
func TestTriggerPipelineWithVariables_Success(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipeline", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			errorhandling.HandleError(err)
		}
		assert.Equal(t, "\"main\"", body["ref"])

		variablesVal, ok := body["variables"]
		assert.True(t, ok, "variables should be in request body")

		// Variables might be an array of objects with Key/Value, validate it exists
		if variablesVal != nil {
			assert.NotNil(t, variablesVal)
		}

		response := gl.Pipeline{
			ID:     132,
			Status: "pending",
			Ref:    "main",
			SHA:    "def456",
			WebURL: "https://gitlab.com/ignorant05/Uniflow/-/pipelines/132",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	variables := map[string]interface{}{
		"VERSION":     "v1.2.3",
		"ENVIRONMENT": "production",
	}
	owner, repo, _ := client.GetDefaultRepository()
	err := client.TriggerPipeline(owner, repo, "main", variables)

	assert.NoError(t, err)
}

// Testing TriggerPipeline, (Failure: invalid ref/branch)
func TestTriggerPipeline_InvalidRef(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipeline", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			errorhandling.HandleError(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errResponse := map[string]interface{}{
			"error": "Invalid ref",
		}
		err = json.NewEncoder(w).Encode(errResponse)
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	owner, repo, _ := client.GetDefaultRepository()
	err := client.TriggerPipeline(owner, repo, "nonexistent-branch", nil)

	assert.Error(t, err)
}

// Testing TriggerPipeline, (Failure: authentication error)
func TestTriggerPipeline_AuthenticationFailure(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipeline", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		errResponse := map[string]interface{}{
			"error": "401 Unauthorized",
		}
		err := json.NewEncoder(w).Encode(errResponse)
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	owner, repo, _ := client.GetDefaultRepository()
	err := client.TriggerPipeline(owner, repo, "main", nil)

	assert.Error(t, err)
}

// Testing TriggerDefaultPipeline, success
func TestTriggerDefaultPipeline_Success(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipeline", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			errorhandling.HandleError(err)
		}
		assert.Equal(t, "\"main\"", body["ref"])

		response := gl.Pipeline{
			ID:     133,
			Status: "pending",
			Ref:    "main",
			SHA:    "ghi789",
			WebURL: "https://gitlab.com/ignorant05/Uniflow/-/pipelines/133",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	err := client.TriggerDefaultPipeline("main", nil)

	assert.NoError(t, err)
}

// Testing TriggerDefaultPipeline with variables, success
func TestTriggerDefaultPipelineWithVariables_Success(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v4/projects/ignorant05/Uniflow/pipeline", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			errorhandling.HandleError(err)
		}
		assert.Equal(t, "\"develop\"", body["ref"])

		variablesVal, ok := body["variables"]
		assert.True(t, ok, "variables should be in request body")

		// Variables might be an array of objects with Key/Value, validate it exists
		if variablesVal != nil {
			assert.NotNil(t, variablesVal)
		}

		response := gl.Pipeline{
			ID:     134,
			Status: "pending",
			Ref:    "develop",
			SHA:    "jkl012",
			WebURL: "https://gitlab.com/ignorant05/Uniflow/-/pipelines/134",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	variables := map[string]interface{}{
		"DEPLOY_ENV": "staging",
	}
	err := client.TriggerDefaultPipeline("develop", variables)

	assert.NoError(t, err)
}

// Testing TriggerPipeline, (Failure: rate limit)
func TestTriggerPipeline_RateLimit(t *testing.T) {
	server, client := mock.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("RateLimit-Remaining", "0")
		w.WriteHeader(http.StatusTooManyRequests)
		errResponse := map[string]interface{}{
			"error": "rate limit exceeded",
		}
		err := json.NewEncoder(w).Encode(errResponse)
		if err != nil {
			errorhandling.HandleError(err)
		}
	})
	defer server.Close()

	owner, repo, _ := client.GetDefaultRepository()
	err := client.TriggerPipeline(owner, repo, "main", nil)

	assert.Error(t, err)
}
