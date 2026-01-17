package platforms

import (
	"context"
	"fmt"
	"maps"

	"github.com/ignorant05/Uniflow/platforms/configurations/gitlab"
	glClient "gitlab.com/gitlab-org/api/client-go"

	gllogs "github.com/ignorant05/Uniflow/platforms/configurations/gitlab/logs"
	"github.com/ignorant05/Uniflow/platforms/constants"
	"github.com/ignorant05/Uniflow/types"
)

type GitlabAdapter struct {
	Client    *gitlab.Client
	projectID string
	owner     string
	repo      string
}

// NewGitlabAdapter creates an adapter object
//
// Parameters:
//   - client: gitlab client
//
// Example:
// adapter, err := NewGitlabAdapter(client)
func NewGitlabAdapter(client *gitlab.Client) (*GitlabAdapter, error) {
	owner, repo, err := client.GetDefaultRepository()
	if err != nil {
		return nil, err
	}

	projectPath, err := client.GetRepositoryPath()
	if err != nil {
		return nil, err
	}

	return &GitlabAdapter{
		Client:    client,
		projectID: projectPath,
		owner:     owner,
		repo:      repo,
	}, nil
}

// TriggerPipeline triggers a pipeline
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// resp, err := a.TriggerPipeline(ctx, &types.TriggerRequest{ WorkflowName: "main",})
func (a *GitlabAdapter) TriggerPipeline(ctx context.Context, req *types.TriggerRequest) (*types.TriggerResponse, error) {
	variables := make(map[string]interface{})
	maps.Copy(variables, req.Inputs)

	targetRef := req.Branch
	if targetRef == "" {
		targetRef = "main"
	}

	err := a.Client.TriggerPipeline(
		a.owner,
		a.repo,
		targetRef,
		variables,
	)

	if err != nil {
		return nil, &types.PlatformError{
			Code:     "trigger_failed",
			Message:  err.Error(),
			Platform: constants.GITLAB_PLATFORM,
		}
	}

	pipelines, err := a.Client.ListPipelines(a.owner, a.repo)
	if err != nil {
		return nil, err
	}

	if len(pipelines) == 0 {
		return nil, fmt.Errorf("<?> Error: no pipelines found after trigger")
	}

	latestPipeline := pipelines[0]

	return &types.TriggerResponse{
		RunID:     int64(latestPipeline.ID),
		RunNumber: int(latestPipeline.ID),
		Status:    latestPipeline.Status,
		QueuedAt:  *latestPipeline.CreatedAt,
	}, nil
}

// GetStatus gets the status of a single pipeline
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// status, err := a.GetStatus(ctx, &types.StatusRequest{ RunID: 1,})
func (a *GitlabAdapter) GetStatus(ctx context.Context, req *types.StatusRequest) (*types.Status, error) {
	pipeline, err := a.Client.GetPipelineStatus(a.owner, a.repo, int(req.RunID))
	if err != nil {
		return nil, &types.PlatformError{
			Code:     "status_failed",
			Message:  err.Error(),
			Platform: constants.GITLAB_PLATFORM,
		}
	}

	status := &types.Status{
		RunID:     int64(pipeline.ID),
		RunNumber: int(pipeline.ID),
		Status:    pipeline.Status,
		StartedAt: *pipeline.CreatedAt,
		URL:       pipeline.WebURL,
	}

	if pipeline.UpdatedAt != nil {
		status.CompletedAt = *pipeline.UpdatedAt
		status.Duration = pipeline.UpdatedAt.Sub(*pipeline.CreatedAt)
	}

	return status, nil
}

// ListPipelines lists all pipelines (filtered by dispatch if requested)
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// pipelines, err := a.ListPipelines(ctx, &types.ListWorkflowsRequest{ WithDispatch: true,})
func (a *GitlabAdapter) ListPipelines(ctx context.Context, req *types.ListWorkflowsRequest) ([]*types.Workflow, error) {
	var allPipelines []*glClient.PipelineInfo

	if req.WithDispatch {
		jobs, err := a.Client.ListJobsWithManualTriggerOnly(a.owner, a.repo)
		if err != nil {
			return nil, &types.PlatformError{
				Code:     "forbidden",
				Message:  err.Error(),
				Platform: constants.GITLAB_PLATFORM,
			}
		}

		pipelines, err := a.Client.ListPipelines(a.owner, a.repo)
		if err != nil {
			return nil, &types.PlatformError{
				Code:     "forbidden",
				Message:  err.Error(),
				Platform: constants.GITLAB_PLATFORM,
			}
		}

		// Filter pipelines that have manual trigger jobs
		filteredPipelines := make([]*glClient.PipelineInfo, 0)
		for _, pipeline := range pipelines {
			for _, job := range jobs {
				if job != "" {
					filteredPipelines = append(filteredPipelines, pipeline)
					break
				}
			}
		}
		allPipelines = filteredPipelines
	} else {
		pipelines, err := a.Client.ListPipelines(a.owner, a.repo)
		if err != nil {
			return nil, &types.PlatformError{
				Code:     "forbidden",
				Message:  err.Error(),
				Platform: constants.GITLAB_PLATFORM,
			}
		}
		allPipelines = pipelines
	}

	workflows := make([]*types.Workflow, 0, len(allPipelines))
	for _, p := range allPipelines {
		workflow := &types.Workflow{
			ID:           int64(p.ID),
			Name:         fmt.Sprintf("Pipeline %d", p.ID),
			Path:         p.Ref,
			State:        p.Status,
			URL:          p.WebURL,
			WithDispatch: req.WithDispatch,
		}
		workflows = append(workflows, workflow)
	}

	return workflows, nil
}

// ListPipelineJobs lists all pipeline jobs
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// jobs, err := a.ListPipelineJobs(ctx, &types.ListWokflowJobsRequest{ WorkflowName: "main"})
func (a *GitlabAdapter) ListPipelineJobs(ctx context.Context, req *types.ListWokflowJobsRequest) ([]*types.WorkflowJob, error) {
	// GitLab pipelines don't have workflow names, we use the latest pipeline
	pipelines, err := a.Client.ListPipelines(a.owner, a.repo)
	if err != nil {
		return nil, err
	}

	if len(pipelines) == 0 {
		return nil, &types.PlatformError{
			Code:     "not_found",
			Message:  "<?> Error: No pipelines found.",
			Platform: constants.GITLAB_PLATFORM,
		}
	}

	latestPipeline := pipelines[0]
	allJobs, err := a.Client.GetPipelineJobs(a.owner, a.repo, int(latestPipeline.ID))
	if err != nil {
		return nil, err
	}

	jobs := make([]*types.WorkflowJob, 0, len(allJobs))
	for _, job := range allJobs {
		if req.Status != "" && req.Status != job.Status {
			continue
		}

		jobs = append(jobs, &types.WorkflowJob{
			ID:           int64(job.ID),
			RunID:        int64(latestPipeline.ID),
			WorkflowName: latestPipeline.Ref,
			Name:         job.Name,
			Status:       job.Status,
			Conclusion:   job.FailureReason,
			RunURL:       latestPipeline.WebURL,
			URL:          job.WebURL,
			HTMLURL:      job.WebURL,
		})
	}

	return jobs, nil
}

// ListPipelineRuns lists all pipeline runs
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// runs, err := a.ListPipelineRuns(ctx, &types.ListWorkflowRunsRequest{ Branch: "main"})
func (a *GitlabAdapter) ListPipelineRuns(ctx context.Context, req *types.ListWorkflowRunsRequest) ([]*types.Run, error) {
	glPipelines, err := a.Client.ListPipelines(a.owner, a.repo)
	if err != nil {
		return nil, err
	}

	runs := make([]*types.Run, 0, len(glPipelines))
	for idx, p := range glPipelines {
		if req.Limit > 0 && idx >= req.Limit {
			break
		}

		if req.Status != "" && req.Status != p.Status {
			continue
		}

		if req.Branch != "" && req.Branch != p.Ref {
			continue
		}

		runs = append(runs, &types.Run{
			RunID:      int64(p.ID),
			RunNumber:  int(p.ID),
			Status:     p.Status,
			Conclusion: "",
			Branch:     p.Ref,
			Event:      "pipeline",
			CommitSHA:  p.SHA,
			CreatedAt:  *p.CreatedAt,
			UpdatedAt:  *p.UpdatedAt,
			URL:        p.WebURL,
		})
	}

	return runs, nil
}

// StreamLogs streams pipeline logs
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//   - callback: logs callback
//
// Example:
// err := a.StreamLogs(ctx, &types.LogsStreamRequest{ ID: 1,})
func (a *GitlabAdapter) StreamLogs(ctx context.Context, req *types.LogsStreamRequest, callback *types.LogCallback) error {
	_, err := a.Client.GetPipelineStatus(a.owner, a.repo, int(req.RunID))
	if err != nil {
		return &types.PlatformError{
			Code:     "logs_failed",
			Message:  err.Error(),
			Platform: constants.GITLAB_PLATFORM,
		}
	}

	return nil
}

// ListPipelineRunLogs lists the logs for a pipeline
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// logsResp, err := a.ListPipelineRunLogs(ctx, &types.LogsRequest{ RunID: 1,})
func (a *GitlabAdapter) ListPipelineRunLogs(ctx context.Context, req *types.LogsRequest) (*types.LogsResponse, error) {
	path := req.DownloadPath
	if req.DownloadPath == "" {
		path = fmt.Sprintf("pipeline-%d", req.RunID)
	}

	err := gllogs.DownloadLogs(a.projectID, path)
	if err != nil {
		return nil, &types.PlatformError{
			Code:     "logs_failed",
			Message:  err.Error(),
			Platform: constants.GITLAB_PLATFORM,
		}
	}

	pipeline, err := a.Client.GetPipelineStatus(a.owner, a.repo, int(req.RunID))
	if err != nil {
		return nil, err
	}

	return &types.LogsResponse{
		URL: pipeline.WebURL,
	}, nil
}

// GetPipelineSummary gets the summary of a pipeline
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// summary, err := a.GetPipelineSummary(ctx, &types.Workflow{ ID: 1,})
func (a *GitlabAdapter) GetPipelineSummary(ctx context.Context, req *types.Workflow) (*types.WorkflowRunSummary, error) {
	pipelineDetails, err := a.Client.GetPipelineSummary(a.owner, a.repo, int(req.ID))
	if err != nil {
		return nil, err
	}

	return &types.WorkflowRunSummary{
		ID:         int64(pipelineDetails.ID),
		Name:       fmt.Sprintf("Pipeline %d", pipelineDetails.ID),
		Status:     pipelineDetails.Status,
		Conclusion: "",
	}, nil
}

// CancelPipeline cancels a pipeline
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// err := a.CancelPipeline(ctx, &types.Run{ RunID: 1,})
func (a *GitlabAdapter) CancelPipeline(ctx context.Context, req *types.Run) error {
	return a.Client.CancelPipeline(a.owner, a.repo, int(req.RunID))
}

// GetRepository returns current repository elements (owner/repo)
//
// Parameters:
//   - ctx: the context variable
//
// Example:
// owner, repo := a.GetRepository(ctx)
func (a GitlabAdapter) GetRepository(ctx context.Context) (string, string) {
	return a.owner, a.repo
}

// GetRepositoryInfo returns repository information
//
// Parameters:
//   - ctx: the context variable
//
// Example:
// info, err := a.GetRepositoryInfo(ctx)
func (a *GitlabAdapter) GetRepositoryInfo(ctx context.Context) (*types.RepositoryInfo, error) {
	return a.Client.GetRepositoryInfo(a.owner, a.repo)
}

// GetUnderlyingClient returns the gitlab client from the GitlabAdapter struct as an interface
//
// Parameters:
//   - None
//
// Example:
// client := a.GetUnderlyingClient()
func (a *GitlabAdapter) GetUnderlyingClient() interface{} {
	return a.Client
}

// IsGitlab verifies that the current client is a gitlab client
//
// Parameters:
//   - None
//
// Example:
// valid := a.IsGitlab()
func (a *GitlabAdapter) IsGitlab() bool {
	return true
}
