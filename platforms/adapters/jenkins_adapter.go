package platforms

import (
	"context"
	"fmt"
	"time"

	errorhandling "github.com/ignorant05/Uniflow/internal/errorHandling"
	"github.com/ignorant05/Uniflow/platforms/configurations/jenkins"
	jlogs "github.com/ignorant05/Uniflow/platforms/configurations/jenkins/logs"
	"github.com/ignorant05/Uniflow/platforms/constants"
	"github.com/ignorant05/Uniflow/types"
)

type JenkinsAdapter struct {
	Client  *jenkins.Client
	jobName string
}

// NewJenkinsAdapter creates an adapter object
//
// Parameters:
//   - client: jenkins client
//
// Example:
// adapter, err := NewJenkinsAdapter(client)
func NewJenkinsAdapter(client *jenkins.Client) (*JenkinsAdapter, error) {
	jobName, err := client.GetDefaultJob()
	if err != nil {
		return nil, err
	}

	return &JenkinsAdapter{
		Client:  client,
		jobName: jobName,
	}, nil
}

// TriggerWorkflow triggers a Jenkins job
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// resp, err := a.TriggerWorkflow(ctx, &types.TriggerRequest{ WorkflowName: "my-pipeline"})
func (a *JenkinsAdapter) TriggerWorkflow(ctx context.Context, req *types.TriggerRequest) (*types.TriggerResponse, error) {
	targetJob := req.WorkflowName
	if targetJob == "" {
		targetJob = a.jobName
	}

	queueID, err := a.Client.TriggerJob(targetJob, req.Inputs)
	if err != nil {
		return nil, &types.PlatformError{
			Code:     "trigger_failed",
			Message:  err.Error(),
			Platform: constants.JENKINS_PLATFORM,
		}
	}

	return &types.TriggerResponse{
		RunID:     queueID,
		RunNumber: int(queueID),
		Status:    "queued",
		QueuedAt:  time.Now(),
	}, nil
}

// GetStatus gets the status of a single job build
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// status, err := a.GetStatus(ctx, &types.StatusRequest{ RunID: 1})
func (a *JenkinsAdapter) GetStatus(ctx context.Context, req *types.StatusRequest) (*types.Status, error) {
	build, err := a.Client.GetJobBuild(a.jobName, req.RunID)
	if err != nil {
		return nil, &types.PlatformError{
			Code:     "status_failed",
			Message:  err.Error(),
			Platform: constants.JENKINS_PLATFORM,
		}
	}

	running, err := a.Client.IsJobRunning(a.jobName, req.RunID)
	if err != nil {
		return nil, err
	}

	status := &types.Status{
		RunID:     int64(build.GetBuildNumber()),
		RunNumber: int(build.GetBuildNumber()),
		Status:    getJenkinsStatus(running),
		StartedAt: build.GetTimestamp(),
		URL:       build.GetUrl(),
	}

	if !running {
		status.Conclusion = build.GetResult()
		status.CompletedAt = build.GetTimestamp().Add(time.Duration(build.GetDuration()) * time.Millisecond)
		status.Duration = time.Duration(build.GetDuration()) * time.Millisecond
	}

	return status, nil
}

// ListWorkflows lists all Jenkins jobs
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// jobs, err := a.ListWorkflows(ctx, &types.ListWorkflowsRequest{})
func (a *JenkinsAdapter) ListWorkflows(ctx context.Context, req *types.ListWorkflowsRequest) ([]*types.Workflow, error) {
	jobs, err := a.Client.ListJobs()
	if err != nil {
		return nil, &types.PlatformError{
			Code:     "forbidden",
			Message:  err.Error(),
			Platform: constants.JENKINS_PLATFORM,
		}
	}

	workflows := make([]*types.Workflow, 0, len(jobs))
	for _, job := range jobs {
		workflow := &types.Workflow{
			ID:   int64(job.GetDetails().LastBuild.Number),
			Name: job.GetName(),
			Path: job.GetName(),
			State: func() string {
				if ok, err := job.IsEnabled(ctx); ok && err != nil {
					return "disabled"
				}
				if err == nil {
					errMsg := fmt.Errorf("something went wrong: %w", err)
					errorhandling.HandleError(errMsg)
				}
				return "active"
			}(),
			URL: job.GetDetails().URL,
		}
		workflows = append(workflows, workflow)
	}

	return workflows, nil
}

// ListWorkflowJobs lists all jobs for a workflow (builds in Jenkins)
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// jobs, err := a.ListWorkflowJobs(ctx, &types.ListWorkflowJobsRequest{ WorkflowName: "my-pipeline"})
func (a *JenkinsAdapter) ListWorkflowJobs(ctx context.Context, req *types.ListWokflowJobsRequest) ([]*types.WorkflowJob, error) {
	targetJob := req.WorkflowName
	if targetJob == "" {
		targetJob = a.jobName
	}

	builds, err := a.Client.GetJobBuilds(targetJob)
	if err != nil {
		return nil, err
	}

	jobs := make([]*types.WorkflowJob, 0, len(builds))
	for _, build := range builds {
		if req.Status != "" && req.Status != build.GetResult() {
			continue
		}

		jobs = append(jobs, &types.WorkflowJob{
			ID:           int64(build.GetBuildNumber()),
			RunID:        int64(build.GetBuildNumber()),
			WorkflowName: targetJob,
			Name:         fmt.Sprintf("%s #%d", targetJob, build.GetBuildNumber()),
			Status: func() string {
				running, _ := a.Client.IsJobRunning(targetJob, int64(build.GetBuildNumber()))
				if running {
					return "in_progress"
				}
				return "completed"
			}(),
			Conclusion: build.GetResult(),
			RunURL:     build.GetUrl(),
			URL:        build.GetUrl(),
			HTMLURL:    build.GetUrl(),
		})
	}

	return jobs, nil
}

// ListWorkflowRuns lists all job builds (runs in Jenkins)
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// runs, err := a.ListWorkflowRuns(ctx, &types.ListWorkflowRunsRequest{ WorkflowName: "my-pipeline"})
func (a *JenkinsAdapter) ListWorkflowRuns(ctx context.Context, req *types.ListWorkflowRunsRequest) ([]*types.Run, error) {
	targetJob := req.WorkflowName
	if targetJob == "" {
		targetJob = a.jobName
	}

	builds, err := a.Client.GetJobBuilds(targetJob)
	if err != nil {
		return nil, err
	}

	runs := make([]*types.Run, 0, len(builds))
	for idx, build := range builds {
		if req.Limit > 0 && idx >= req.Limit {
			break
		}

		if req.Status != "" && req.Status != build.GetResult() {
			continue
		}

		runs = append(runs, &types.Run{
			RunID:     int64(build.GetBuildNumber()),
			RunNumber: int(build.GetBuildNumber()),
			Status: func() string {
				running, _ := a.Client.IsJobRunning(targetJob, int64(build.GetBuildNumber()))
				if running {
					return "in_progress"
				}
				return "completed"
			}(),
			Conclusion:  build.GetResult(),
			Actor:       "jenkins",
			Event:       "manual",
			TriggeredBy: "jenkins",
			CreatedAt:   build.GetTimestamp(),
			UpdatedAt:   build.GetTimestamp().Add(time.Duration(build.GetDuration()) * time.Millisecond),
			URL:         build.GetUrl(),
		})
	}

	return runs, nil
}

// StreamLogs streams logs from a Jenkins build
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//   - callback: logs callback
//
// Example:
// err := a.StreamLogs(ctx, &types.LogsStreamRequest{ RunID: 1})
func (a *JenkinsAdapter) StreamLogs(ctx context.Context, req *types.LogsStreamRequest, callback *types.LogCallback) error {
	_, err := a.Client.GetBuildLogs(a.jobName, req.RunID)
	if err != nil {
		return &types.PlatformError{
			Code:     "logs_failed",
			Message:  err.Error(),
			Platform: constants.JENKINS_PLATFORM,
		}
	}

	return nil
}

// ListWorkflowRunLogs gets the logs for a job build
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// logs, err := a.ListWorkflowRunLogs(ctx, &types.LogsRequest{ RunID: 1})
func (a *JenkinsAdapter) ListWorkflowRunLogs(ctx context.Context, req *types.LogsRequest) (*types.LogsResponse, error) {
	_, err := a.Client.GetJobBuild(a.jobName, req.RunID)
	if err != nil {
		return nil, &types.PlatformError{
			Code:     "logs_failed",
			Message:  err.Error(),
			Platform: constants.JENKINS_PLATFORM,
		}
	}

	logsURL := fmt.Sprintf("%s/consoleText", getJobBuildURL(a.jobName, req.RunID))

	path := req.DownloadPath
	if req.DownloadPath == "" {
		path = req.WorkflowName
	}

	err = jlogs.DownloadLogs(logsURL, path)
	if err != nil {
		return nil, err
	}

	return &types.LogsResponse{
		URL: logsURL,
	}, nil
}

// GetWorkflowRunSummary gets summary information about a job build
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// summary, err := a.GetWorkflowRunSummary(ctx, &types.Workflow{ ID: 1})
func (a *JenkinsAdapter) GetWorkflowRunSummary(ctx context.Context, req *types.Workflow) (*types.WorkflowRunSummary, error) {
	build, err := a.Client.GetJobBuild(a.jobName, req.ID)
	if err != nil {
		return nil, err
	}

	running, _ := a.Client.IsJobRunning(a.jobName, req.ID)

	return &types.WorkflowRunSummary{
		ID:   int64(build.GetBuildNumber()),
		Name: fmt.Sprintf("%s #%d", a.jobName, build.GetBuildNumber()),
		Status: func() string {
			if running {
				return "in_progress"
			}
			return "completed"
		}(),
		Conclusion: build.GetResult(),
	}, nil
}

// Cancel cancels a running job build
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// err := a.Cancel(ctx, &types.Run{ RunID: 1})
func (a *JenkinsAdapter) Cancel(ctx context.Context, req *types.Run) error {
	return a.Client.CancelExecution(a.jobName, req.RunID)
}

// GetRepository returns the job name (Jenkins equivalent of owner/repo)
//
// Parameters:
//   - ctx: the context variable
//
// Example:
// jobName := a.GetRepository(ctx)
func (a JenkinsAdapter) GetRepository(ctx context.Context) (string, string) {
	return a.jobName, ""
}

// GetRepositoryInfo returns job information
//
// Parameters:
//   - ctx: the context variable
//
// Example:
// info, err := a.GetRepositoryInfo(ctx)
func (a *JenkinsAdapter) GetRepositoryInfo(ctx context.Context) (*types.RepositoryInfo, error) {
	job, err := a.Client.GetJobInfo(a.jobName)
	if err != nil {
		return nil, err
	}

	return &types.RepositoryInfo{
		Name:    job.GetName(),
		HTMLURL: job.GetDetails().URL,
	}, nil
}

// GetUnderlyingClient returns the jenkins client as an interface
//
// Parameters:
//   - None
//
// Example:
// client := a.GetUnderlyingClient()
func (a *JenkinsAdapter) GetUnderlyingClient() interface{} {
	return a.Client
}

// IsGithub verifies that the current client is a jenkins client
//
// Parameters:
//   - None
//
// Example:
// valid:= a.IsJenkins()
func (a *JenkinsAdapter) IsJenkins() bool {
	return true
}

// Helper functions
func getJenkinsStatus(running bool) string {
	if running {
		return "in_progress"
	}
	return "completed"
}

func getJobBuildURL(jobName string, buildID int64) string {
	return fmt.Sprintf("/job/%s/%d", jobName, buildID)
}
