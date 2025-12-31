package platforms

import (
	"context"

	"github.com/ignorant05/Uniflow/types"
)

// Global abstraction interface for every platform to come
type PlatformClient interface {
	// Triggers a workflow
	TriggerWorkflow(ctx context.Context, req *types.TriggerRequest) (*types.TriggerResponse, error)

	// Lists workflow runs
	ListWorkflowRuns(ctx context.Context, req *types.ListWorkflowRunsRequest) ([]*types.Run, error)

	// Lists workflow jobs
	ListWorkflowJobs(ctx context.Context, req *types.ListWokflowJobsRequest) ([]*types.WorkflowJob, error)

	// Lists workflows
	ListWorkflows(ctx context.Context, req *types.ListWorkflowsRequest) ([]*types.Workflow, error)

	// Retrieves the status of a specific workflow
	GetStatus(ctx context.Context, req *types.StatusRequest) (*types.Status, error)

	// Retrieves the run summary of a specific workflow
	GetWorkflowRunSummary(ctx context.Context, req *types.Workflow) (*types.WorkflowRunSummary, error)

	// Streams log in real time
	StreamLogs(ctx context.Context, req *types.LogsStreamRequest, callback *types.LogCallback) error

	// List all workflow run logs (for a specific workflow)
	ListWorkflowRunLogs(ctx context.Context, req *types.LogsRequest) (*types.LogsResponse, error)

	// Cancels streaming
	Cancel(ctx context.Context, req *types.Run) error

	// Retrieves the repository corresponding to the current working dir if exists (format: owner/repo)
	GetRepository(ctx context.Context) (string, string)

	// Retrieves current repository info
	GetRepositoryInfo(ctx context.Context) (*types.RepositoryInfo, error)

	// Retrieves client
	GetUnderlyingClient() interface{}

	// Verifies if the current client is a github client
	IsGithub() bool
}
