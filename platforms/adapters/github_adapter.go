package platforms

import (
	"context"
	"fmt"
	"maps"
	"strings"

	githubClient "github.com/google/go-github/v57/github"
	"github.com/ignorant05/Uniflow/platforms/constants"
	"github.com/ignorant05/Uniflow/platforms/github"
	ghlogs "github.com/ignorant05/Uniflow/platforms/github/logs"
	"github.com/ignorant05/Uniflow/types"
)

type GithubAdapter struct {
	Client *github.Client
	owner  string
	repo   string
}

// NewGithubAdapter creates an adapter object
//
// Parameters:
//   - client: github client
//
// Example:
// adapter, err := NewGithubAdapter(client)
func NewGithubAdapter(client *github.Client) (*GithubAdapter, error) {
	owner, repo, err := client.GetDefaultRepository()
	if err != nil {
		return nil, err
	}

	return &GithubAdapter{
		Client: client,
		owner:  owner,
		repo:   repo,
	}, nil
}

// TriggerWorkflow triggers a workflow
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// resp, err := a.TriggerWorkflow(ctx, &types.TriggerRequest{ WorkflowName: "deploy.yml",})
func (a *GithubAdapter) TriggerWorkflow(ctx context.Context, req *types.TriggerRequest) (*types.TriggerResponse, error) {
	inputs := make(map[string]interface{})
	maps.Copy(inputs, req.Inputs)

	targetWorkflow := req.WorkflowName

	if req.WorkflowName == "" {
		targetWorkflow = constants.DEFAULT_WORKFLOW
	}

	err := a.Client.TriggerWorkflow(
		a.owner,
		a.repo,
		targetWorkflow,
		req.Branch,
		inputs,
	)

	if err != nil {
		return nil, &types.PlatformError{
			Code:     "trigger_failed",
			Message:  err.Error(),
			Platform: constants.GITHUB_PLATFORM,
		}
	}

	workflows, err := a.Client.ListWorkflows(a.owner, a.repo)
	if err != nil {
		return nil, err
	}

	var workflowID int64
	for _, wf := range workflows {
		if strings.Contains(wf.GetPath(), req.WorkflowName) {
			workflowID = wf.GetID()
		}
	}

	if workflowID == 0 {
		return nil, fmt.Errorf("<?> Error: workflow not found: %s", req.WorkflowName)
	}

	runs, err := a.Client.GetWorkflowRuns(a.owner, a.repo, workflowID)
	if err != nil {
		return nil, err
	}

	latestRun := runs[0]

	return &types.TriggerResponse{
		RunID:     latestRun.GetID(),
		RunNumber: latestRun.GetRunNumber(),
		Status:    latestRun.GetStatus(),
		QueuedAt:  latestRun.GetCreatedAt().Time,
	}, nil
}

// GetStatus gets the status of a single workflow
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// workflows, err := a.GetStatus(ctx, &types.StatusRequest{ RunID: 1,})
func (a *GithubAdapter) GetStatus(ctx context.Context, req *types.StatusRequest) (*types.Status, error) {
	run, err := a.Client.GetWorkflowRunStatus(a.owner, a.repo, req.RunID)
	if err != nil {
		return nil, &types.PlatformError{
			Code:     "status_failed",
			Message:  err.Error(),
			Platform: constants.GITHUB_PLATFORM,
		}
	}

	status := &types.Status{
		RunID:     run.GetID(),
		RunNumber: run.GetRunNumber(),
		Status:    run.GetStatus(),
		StartedAt: run.GetRunStartedAt().Time,
		URL:       run.GetRerunURL(),
	}

	if run.GetConclusion() != "" {
		status.Conclusion = run.GetConclusion()
	}

	if run.GetUpdatedAt().GoString() != "" {
		completedAt := run.GetUpdatedAt().Time
		status.CompletedAt = completedAt
		status.Duration = completedAt.Sub(run.GetRunStartedAt().Time)
	}

	return status, nil
}

// ListWorkflows lists all workflows
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// workflows, err := a.ListWorkflows(ctx, &types.ListWorkflowsRequest{ WithDispatch: true,})
func (a *GithubAdapter) ListWorkflows(ctx context.Context, req *types.ListWorkflowsRequest) ([]*types.Workflow, error) {
	var allworkflows []*githubClient.Workflow

	if req.WithDispatch {
		ws, err := a.Client.ListWorkflowsWithDispatchOnly(a.owner, a.repo)
		if err != nil {
			return nil, &types.PlatformError{
				Code:     "forbidden",
				Message:  err.Error(),
				Platform: constants.GITHUB_PLATFORM,
			}
		}
		allworkflows = ws
	} else {
		ws, err := a.Client.ListWorkflows(a.owner, a.repo)
		if err != nil {
			return nil, &types.PlatformError{
				Code:     "forbidden",
				Message:  err.Error(),
				Platform: constants.GITHUB_PLATFORM,
			}
		}
		allworkflows = ws
	}

	workflows := make([]*types.Workflow, 0, len(allworkflows))
	for _, wf := range allworkflows {
		workflow := &types.Workflow{
			ID:           wf.GetID(),
			Name:         wf.GetName(),
			Path:         wf.GetPath(),
			State:        wf.GetState(),
			URL:          wf.GetURL(),
			WithDispatch: req.WithDispatch,
		}
		workflows = append(workflows, workflow)
	}

	return workflows, nil
}

// ListWorkflowJobs lists all workflow job
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// runs, err := a.ListWorkflowJobs(ctx, &types.ListWorkflowJobsRequest{ WorkflowName: "deploy.yml"})
func (a *GithubAdapter) ListWorkflowJobs(ctx context.Context, req *types.ListWokflowJobsRequest) ([]*types.WorkflowJob, error) {
	var workflowID int64
	if req.WorkflowName != "" {
		workflows, err := a.Client.ListWorkflows(a.owner, a.repo)
		if err != nil {
			return nil, err
		}

		if len(workflows) <= 0 {
			return nil, &types.PlatformError{
				Code:     "not_found",
				Message:  "<?> Error: No workflows found.",
				Platform: constants.GITHUB_PLATFORM,
			}
		}

		for _, wf := range workflows {
			if strings.Contains(wf.GetPath(), req.WorkflowName) {
				workflowID = wf.GetID()
				break
			}
		}

		if workflowID == 0 {
			return nil, &types.PlatformError{
				Code:     "not_found",
				Message:  "<?> Error: No workflows found.",
				Platform: constants.GITHUB_PLATFORM,
			}
		}
	}

	alljobs, err := a.Client.ListWorkflowJobs(a.owner, a.repo, workflowID)
	if err != nil {
		return nil, err
	}

	jobs := make([]*types.WorkflowJob, 0, len(alljobs))
	for _, job := range alljobs {
		if req.Status != "" && req.Status != job.GetStatus() {
			continue
		}

		if req.Branch != "" && req.Branch != job.GetHeadBranch() {
			continue
		}

		jobs = append(jobs, &types.WorkflowJob{
			ID:           job.GetID(),
			RunID:        job.GetRunID(),
			WorkflowName: job.GetWorkflowName(),
			Name:         job.GetName(),
			Status:       job.GetStatus(),
			Conclusion:   job.GetConclusion(),
			RunURL:       job.GetRunURL(),
			URL:          job.GetURL(),
			HTMLURL:      job.GetHTMLURL(),
		})
	}

	return jobs, nil
}

// ListWorkflowRuns lists all workflow runs (all that is <= limit)
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// runs, err := a.ListWorkflowRuns(ctx, &types.ListWorkflowRunsRequest{ WorkflowName: "deploy.yml"})
func (a *GithubAdapter) ListWorkflowRuns(ctx context.Context, req *types.ListWorkflowRunsRequest) ([]*types.Run, error) {
	var workflowID int64
	if req.WorkflowName != "" {
		workflows, err := a.Client.ListWorkflows(a.owner, a.repo)
		if err != nil {
			return nil, err
		}

		if len(workflows) <= 0 {
			return nil, &types.PlatformError{
				Code:     "not_found",
				Message:  "<?> Error: No workflows found.",
				Platform: constants.GITHUB_PLATFORM,
			}
		}

		for _, wf := range workflows {
			if strings.Contains(wf.GetPath(), req.WorkflowName) {
				workflowID = wf.GetID()
				break
			}
		}

		if workflowID == 0 {
			return nil, &types.PlatformError{
				Code:     "not_found",
				Message:  "<?> Error: No workflows found.",
				Platform: constants.GITHUB_PLATFORM,
			}
		}
	}
	ghruns, err := a.Client.GetWorkflowRuns(a.owner, a.repo, workflowID)
	if err != nil {
		return nil, err
	}

	runs := make([]*types.Run, 0, len(ghruns))
	for idx, r := range ghruns {
		if req.Limit > 0 && idx >= req.Limit {
			break
		}

		if req.Status != "" && req.Status != r.GetStatus() {
			continue
		}

		if req.Branch != "" && req.Branch != r.GetHeadBranch() {
			continue
		}

		runs = append(runs, &types.Run{
			RunID:       r.GetID(),
			RunNumber:   r.GetRunNumber(),
			Status:      r.GetStatus(),
			Conclusion:  r.GetConclusion(),
			Branch:      r.GetHeadBranch(),
			Actor:       r.GetActor().GetLogin(),
			Event:       r.GetEvent(),
			CommitSHA:   r.GetHeadCommit().GetSHA(),
			TriggeredBy: r.GetActor().GetLogin(),
			CreatedAt:   r.GetCreatedAt().Time,
			UpdatedAt:   r.GetUpdatedAt().Time,
			URL:         r.GetURL(),
		})
	}

	return runs, nil
}

// GetGithubClient streams logs
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//   - callback: logs callback
//
// Example:
// err := a.StreamLogs(ctx, &types.LogsStreamRequest{ ID: 1,})
func (a *GithubAdapter) StreamLogs(ctx context.Context, req *types.LogsStreamRequest, callback *types.LogCallback) error {
	_, err := a.Client.GetWorkflowRunLogs(a.owner, a.repo, req.RunID)
	if err != nil {
		return &types.PlatformError{
			Code:     "logs_failed",
			Message:  err.Error(),
			Platform: constants.GITHUB_PLATFORM,
		}
	}

	return nil
}

// ListWorkflowRunLogs lists the logs for a workflow
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// summary, err := a.ListWorkflowRunLogs(ctx, &types.LogsRequest{ RunID: 1,})
func (a *GithubAdapter) ListWorkflowRunLogs(ctx context.Context, req *types.LogsRequest) (*types.LogsResponse, error) {
	logsURL, err := a.Client.GetWorkflowRunLogs(a.owner, a.repo, req.RunID)
	if err != nil {
		return nil, &types.PlatformError{
			Code:     "logs_failed",
			Message:  err.Error(),
			Platform: constants.GITHUB_PLATFORM,
		}
	}

	path := req.DownloadPath
	if req.DownloadPath == "" {
		path = req.WorkflowName
	}

	err = ghlogs.DownloadLogs(logsURL, path)
	if err != nil {
		return nil, err
	}

	return &types.LogsResponse{
		URL: logsURL,
	}, nil
}

// GetGithubClient returns the github client from the GithubAdapter struct
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// summary, err := a.GetWorkflowRunSummary(ctx, &types.Workflow{ ID: 1,})
func (a *GithubAdapter) GetWorkflowRunSummary(ctx context.Context, req *types.Workflow) (*types.WorkflowRunSummary, error) {
	workflowDetails, err := a.Client.GetWorkflowRunSummary(a.owner, a.repo, req.ID)
	if err != nil {
		return nil, err
	}

	return &types.WorkflowRunSummary{
		ID:         workflowDetails.ID,
		Name:       workflowDetails.Name,
		Status:     workflowDetails.Status,
		Conclusion: workflowDetails.Conclusion,
	}, nil
}

// GetGithubClient returns the github client from the GithubAdapter struct
//
// Parameters:
//   - ctx: the context variable
//   - req: the request body
//
// Example:
// err := a.Cancel(ctx, &types.Run{ RunID: 1,})
func (a *GithubAdapter) Cancel(ctx context.Context, req *types.Run) error {
	return a.Client.CancelWorkflowRun(a.owner, a.repo, req.RunID)
}

// GetGithubClient returns current repository elements (owner/repo)
//
// Parameters:
//   - ctx: the context variable
//
// Example:
// owner, repo := a.GetRepository(ctx)
func (a GithubAdapter) GetRepository(ctx context.Context) (string, string) {
	return a.owner, a.repo
}

// GetGithubClient returns the github client from the GithubAdapter struct
//
// Parameters:
//   - ctx: the context variable
//
// Example:
// info, err := a.GetRepositoryInfo(ctx)
func (a *GithubAdapter) GetRepositoryInfo(ctx context.Context) (*types.RepositoryInfo, error) {
	return a.Client.GetRepositoryInfo(a.owner, a.repo)
}

// GetUnderlyingClient returns the github client from the GithubAdapter struct but as an interface
//
// Parameters:
//   - None
//
// Example:
// client := a.GetUnderlyingClient()
func (a *GithubAdapter) GetUnderlyingClient() interface{} {
	return a.Client
}

// IsGithub verifies that the current client is a github client
//
// Parameters:
//   - None
//
// Example:
// valid:= a.IsGithub()
func (a *GithubAdapter) IsGithub() bool {
	return true
}
