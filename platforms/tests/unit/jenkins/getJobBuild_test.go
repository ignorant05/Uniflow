package jenkins_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ignorant05/Uniflow/platforms/tests/unit/jenkins"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type JenkinsBuild struct {
	Number int64  `json:"number"`
	Result string `json:"result"`
	URL    string `json:"url"`
}

// Testing GetWorkflowRun, in progress
func TestGetJobBuild_InProgress(t *testing.T) {
	server, client := jenkins.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		response := JenkinsBuild{
			Number: 42,
			Result: "IN_PROGRESS",
			URL:    "/job/test-job/42/",
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	})

	defer server.Close()

	client.Config.BaseURL = server.URL

	jobName := "test-job"
	buildNumber := int64(42)
	build, err := client.GetJobBuild(jobName, buildNumber)

	require.NoError(t, err)
	assert.Equal(t, int64(42), build.GetBuildNumber())
	assert.Equal(t, "IN_PROGRESS", build.GetResult())
	assert.Equal(t, "/job/test-job/42/", build.GetUrl())

}

func TestGetJobBuild_Failure(t *testing.T) {
	server, client := jenkins.SetupTestClientWithMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		response := JenkinsBuild{
			Number: 42,
			Result: "FAILURE",
			URL:    "/job/test-job/42/",
		}

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(response))
	})
	defer server.Close()

	client.Config.BaseURL = server.URL

	jobName := "test-job"
	buildNumber := int64(42)
	build, err := client.GetJobBuild(jobName, buildNumber)

	require.NoError(t, err)
	assert.Equal(t, int64(42), build.GetBuildNumber())
	assert.Equal(t, "FAILURE", build.GetResult())
	assert.Equal(t, "/job/test-job/42/", build.GetUrl())
}
