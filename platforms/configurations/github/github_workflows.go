package github

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v57/github"
	"github.com/ignorant05/Uniflow/platforms/configurations/github/constants"
	"github.com/ignorant05/Uniflow/platforms/configurations/github/helpers"
)

// TriggerWorkflow triggers a GitHub Actions workflow using the workflow_dispatch event.
//
// Parameters:
//   - owner: Repository owner (username or organization)
//   - repo: Repository name
//   - workflowFileName: Name of the workflow file (e.g., "deploy.yml")
//   - ref: Git reference (branch, tag, or commit SHA)
//   - inputs: Workflow inputs as key-value pairs
//
// Returns an error if:
//   - The workflow file doesn't exist
//   - The workflow doesn't have workflow_dispatch trigger
//   - The API request fails
//
// Example:
//
//	err := client.TriggerWorkflow("owner", "repo", "deploy.yml", "main", map[string]interface{}{
//	  "environment": "production",
//	})
func (c *Client) TriggerWorkflow(owner, repo, workflowFileName, ref string, inputs map[string]interface{}) error {
	event := github.CreateWorkflowDispatchEventRequest{
		Ref:    ref,
		Inputs: inputs,
	}

	_, err := c.Actions.CreateWorkflowDispatchEventByFileName(
		c.Ctx,
		owner,
		repo,
		workflowFileName,
		event,
	)

	if err != nil {
		return fmt.Errorf("<?> Error: Failed to trigger workflow.\n<?> Error: %w", err)
	}

	return nil
}

// TriggerDefaultWorkflow triggers the default GitHub Actions workflow using the workflow_dispatch event.
//
// Parameters:
//   - workflowFileName: Name of the workflow file (e.g., "deploy.yml")
//   - ref: Git reference (branch, tag, or commit SHA)
//   - inputs: Workflow inputs as key-value pairs
//
// Returns an error if:
//   - The workflow file doesn't exist
//   - The workflow doesn't have workflow_dispatch trigger
//   - The API request fails
//
// Example:
//
//	err := client.TriggerWorkflow(deploy.yml", "main", map[string]interface{}{
//	  "environment": "production",
//	})
func (c *Client) TriggerDefaultWorkflow(workflowFileName, ref string, inputs map[string]interface{}) error {
	owner, defRepo, err := c.GetDefaultRepository()
	if err != nil {
		return err
	}

	return c.TriggerWorkflow(owner, defRepo, workflowFileName, ref, inputs)
}

func (c *Client) ListWorkflows(owner, repo string) ([]*github.Workflow, error) {
	opts := &github.ListOptions{PerPage: constants.DEFAULT_PER_PAGE}

	workflows, _, err := c.Actions.ListWorkflows(c.Ctx, owner, repo, opts)
	if err != nil {
		return nil, err
	}

	return workflows.Workflows, nil
}

// ListWorkflowJobs lists all workflow jobs
//
// Parameters:
//   - owner: Repository owner (username or organization)
//   - repo: Repository name
//   - runID: workflow ID
//
// Returns an error if:
//   - The workflow doesn't exist (or invalid runID)
//   - The API request fails
//
// Example:
//
//	jobs, err := client.ListWorkflowJobs("owner", "repo", 12345)
func (c *Client) ListWorkflowJobs(owner, repo string, runID int64) ([]*github.WorkflowJob, error) {
	opts := &github.ListWorkflowJobsOptions{
		ListOptions: github.ListOptions{PerPage: constants.DEFAULT_PER_PAGE},
	}

	jobs, _, err := c.Actions.ListWorkflowJobs(c.Ctx, owner, repo, runID, opts)
	if err != nil {
		return nil, err
	}

	return jobs.Jobs, nil
}

// GetWorkflowRuns lists all running workflows.
//
// Parameters:
//   - owner: Repository owner (username or organization)
//   - repo: Repository name
//   - runID: workflow ID
//
// Returns an error if:
//   - The workflow doesn't exist (or invalid runID)
//   - The API request fails
//
// Example:
//
//	runs, err := client.GetWorkflowRuns("owner", "repo", 12345)
func (c *Client) GetWorkflowRuns(owner, repo string, workflowID int64) ([]*github.WorkflowRun, error) {
	opts := &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{PerPage: constants.DEFAULT_PER_PAGE},
	}
	runs, _, err := c.Actions.ListWorkflowRunsByID(c.Ctx, owner, repo, workflowID, opts)
	if err != nil {
		return nil, fmt.Errorf("<?> Error: Failed to get workflow runs by ID: %d.\n<?> Error: %w", workflowID, err)
	}

	return runs.WorkflowRuns, nil
}

// GetWorkflowRunStatus show the status of workflows present in the repo.
//
// Parameters:
//   - owner: Repository owner (username or organization)
//   - repo: Repository name
//   - runID: workflow id
//
// Returns an error if:
//   - The workflow doesn't exist (or invalid runID)
//   - The API request fails
//
// Example:
//
//	workflows, err := client.GetWorkflowRunStatus("owner", "repo", 12345)
func (c *Client) GetWorkflowRunStatus(owner, repo string, runID int64) (*github.WorkflowRun, error) {
	run, _, err := c.Actions.GetWorkflowRunByID(c.Ctx, owner, repo, runID)
	if err != nil {
		return nil, fmt.Errorf("<?> Error: Failed to get workflow run status by runID: %d.\n<?> Error: %w", runID, err)
	}

	return run, nil
}

// GetWorkflowRunLogs lists all workflow logs by ID.
//
// Parameters:
//   - owner: Repository owner (username or organization)
//   - repo: Repository name
//   - runID: workflow ID
//
// Returns an error if:
//   - The workflow doesn't exist (or invalid runID)
//   - The API request fails
//
// Example:
//
//	logs, err := client.GetWorkflowRunLogs("owner", "repo", 12345)
func (c *Client) GetWorkflowRunLogs(owner, repo string, runID int64) (string, error) {
	url, _, err := c.Actions.GetWorkflowRunLogs(c.Ctx, owner, repo, runID, constants.GITHUB_LOGS_MAX_INDIRECT)
	if err != nil {
		return "", fmt.Errorf("<?> Error: Failed to get workflow run logs by runID: %d.\n<?> Error: %w", runID, err)
	}

	return url.String(), nil
}

// CancelWorkflowRun cancels a workflow by it's ID
//
// Parameters:
//   - owner: Repository owner (username or organization)
//   - repo: Repository name
//   - runID: workflow ID
//
// Returns an error if:
//   - The workflow doesn't exist (or invalid runID)
//   - The API request fails
//
// Example:
//
//	err := client.CancelWorkflowRun("owner", "repo", 12345)
func (c *Client) CancelWorkflowRun(owner, repo string, runID int64) error {
	_, err := c.Actions.CancelWorkflowRunByID(c.Ctx, owner, repo, runID)
	if err != nil {
		return fmt.Errorf("<?> Error: Failed to cancel workflow run with runID: %d.\n<?> Error: %w", runID, err)
	}

	return nil
}

// GetWorkflowRunSummary retrieves the workflow summary of the running workflow.
//
// Parameters:
//   - owner: Repository owner (username or organization)
//   - repo: Repository name
//   - runID: workflow ID
//
// Returns an error if:
//   - The workflow doesn't exist (or invalid runID)
//   - The API request fails
//
// Example:
//
//	summary, err := client.GetWorkflowRunSummary("owner", "repo", 12345)
func (c *Client) GetWorkflowRunSummary(owner, repo string, runID int64) (*helpers.WorkflowRunSummary, error) {
	run, err := c.GetWorkflowRunStatus(owner, repo, runID)
	if err != nil {
		return nil, err
	}

	summary := &helpers.WorkflowRunSummary{
		ID:      run.GetID(),
		Name:    run.GetName(),
		Status:  run.GetStatus(),
		HTMLURL: run.GetHTMLURL(),
	}

	if run.Conclusion != nil {
		summary.Conclusion = *run.Conclusion
	}

	if run.CreatedAt != nil {
		summary.CreatedAt = run.CreatedAt.String()
	}

	if run.UpdatedAt != nil {
		summary.UpdatedAt = run.UpdatedAt.String()
	}

	return summary, nil
}

// FilterWorkflowsWithDispatchOnly filters the workflows to only get the ones that have the workflow_dispatch filter.
//
// Parameters:
//   - owner: Repository owner (username or organization)
//   - repo: Repository name
//
// Returns an error if:
//   - The workflow doesn't exist (or invalid runID)
//   - The API request fails
//
// Example:
//
//	workflows, err := client.FilterWorkflowsWithDispatchOnly("owner", "repo")
func (c *Client) FilterWorkflowsWithDispatchOnly(owner, repo string) ([]string, error) {
	var workflowsWithDispatch []string

	_, dirContent, _, err := c.Repositories.GetContents(
		c.Ctx,
		owner,
		repo,
		".github/workflows",
		&github.RepositoryContentGetOptions{},
	)

	if err != nil {
		return nil, fmt.Errorf("<?> Error: Failed to get workflows.\n<?> Error: %w", err)
	}

	for _, content := range dirContent {
		if strings.HasSuffix(*content.Name, ".yaml") || strings.HasSuffix(*content.Name, ".yml") {
			fileContent, _, _, err := c.Repositories.GetContents(
				c.Ctx,
				owner,
				repo,
				fmt.Sprintf("./github/workflows/%s", *content.Name),
				&github.RepositoryContentGetOptions{},
			)

			if err != nil {
				continue
			}

			contentStr, err := fileContent.GetContent()
			if err != nil {
				continue
			}

			if strings.Contains(contentStr, "workflow_dispatch:") {
				workflowsWithDispatch = append(workflowsWithDispatch, *content.Name)

				fmt.Printf("✓ Found workflow_dispatch in: %s\n", *content.Name)
			}
		}
	}

	return workflowsWithDispatch, nil
}

// ListWorkflowsWithDispatchOnly lists only the workflows that contains workflow_dispatch trigger.
//
// Parameters:
//   - owner: Repository owner (username or organization)
//   - repo: Repository name
//
// Returns an error if:
//   - The workflow doesn't exist (or invalid runID)
//   - The API request fails
//
// Example:
//
//	workflows, err := client.ListWorkflowsWithDispatchOnly("owner", "repo")
func (c *Client) ListWorkflowsWithDispatchOnly(owner, repo string) ([]*github.Workflow, error) {
	var workflowsWithDispatch []*github.Workflow

	filteredWorkflows, err := c.FilterWorkflowsWithDispatchOnly(owner, repo)
	if err != nil {
		return nil, err
	}

	allWorkflows, err := c.ListWorkflows(owner, repo)
	if err != nil {
		return nil, err
	}

	for _, wf := range allWorkflows {
		for _, fwf := range filteredWorkflows {
			if fwf == *wf.Name {
				workflowsWithDispatch = append(workflowsWithDispatch, wf)
			}
		}
	}

	return workflowsWithDispatch, nil
}
