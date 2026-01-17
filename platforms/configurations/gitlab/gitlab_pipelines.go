package gitlab

import (
	"fmt"
	"strings"

	"github.com/ignorant05/Uniflow/platforms/configurations/gitlab/constants"
	"github.com/ignorant05/Uniflow/platforms/configurations/gitlab/helpers"
	"github.com/stretchr/testify/assert/yaml"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// TriggerPipeline triggers a GitLab CI/CD pipeline using the trigger API.
//
// Parameters:
//   - owner: Repository owner (username or group)
//   - repo: Repository name
//   - ref: Git reference (branch, tag, or commit SHA)
//   - variables: Pipeline variables as key-value pairs
//
// Returns an error if:
//   - The project doesn't exist
//   - The API request fails
//
// Example:
//
//	err := client.TriggerPipeline("owner", "repo", "main", map[string]interface{}{
//	  "ENVIRONMENT": "production",
//	})
func (c *Client) TriggerPipeline(owner, repo, ref string, variables map[string]interface{}) error {
	projectID := owner + "/" + repo

	ref = gitlab.Stringify(ref)
	opts := &gitlab.CreatePipelineOptions{
		Ref: &ref,
	}

	vars := make([]*gitlab.PipelineVariableOptions, 0)
	for k, v := range variables {
		val := fmt.Sprintf("%v", v)
		vars = append(vars, &gitlab.PipelineVariableOptions{
			Key:   &k,
			Value: &val,
		})
	}
	opts.Variables = &vars

	_, _, err := c.Pipelines.CreatePipeline(projectID, opts)
	if err != nil {
		return fmt.Errorf("<?> Error: Failed to trigger pipeline.\n<?> Error: %w", err)
	}

	return nil
}

// TriggerDefaultPipeline triggers a pipeline on the default repository.
//
// Parameters:
//   - ref: Git reference (branch, tag, or commit SHA)
//   - variables: Pipeline variables as key-value pairs
//
// Returns an error if:
//   - No default repository configured
//   - The project doesn't exist
//   - The API request fails
//
// Example:
//
//	err := client.TriggerDefaultPipeline("main", map[string]interface{}{
//	  "ENVIRONMENT": "production",
//	})
func (c *Client) TriggerDefaultPipeline(ref string, variables map[string]interface{}) error {
	owner, defRepo, err := c.GetDefaultRepository()
	if err != nil {
		return err
	}

	return c.TriggerPipeline(owner, defRepo, ref, variables)
}

// ListPipelines lists all pipelines in a project.
//
// Parameters:
//   - owner: Repository owner (username or group)
//   - repo: Repository name
//
// Returns an error if:
//   - The project doesn't exist
//   - The API request fails
//
// Example:
//
//	pipelines, err := client.ListPipelines("owner", "repo")
func (c *Client) ListPipelines(owner, repo string) ([]*gitlab.PipelineInfo, error) {
	projectID := owner + "/" + repo
	opts := &gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{PerPage: constants.DEFAULT_PER_PAGE},
	}

	pipelines, _, err := c.Pipelines.ListProjectPipelines(projectID, opts)
	if err != nil {
		return nil, err
	}

	return pipelines, nil
}

// GetPipelineJobs lists all jobs in a pipeline.
//
// Parameters:
//   - owner: Repository owner (username or group)
//   - repo: Repository name
//   - pipelineID: Pipeline ID
//
// Returns an error if:
//   - The pipeline doesn't exist
//   - The API request fails
//
// Example:
//
//	jobs, err := client.GetPipelineJobs("owner", "repo", 12345)
func (c *Client) GetPipelineJobs(owner, repo string, pipelineID int) ([]*gitlab.Job, error) {
	projectID := owner + "/" + repo
	opts := &gitlab.ListJobsOptions{
		ListOptions: gitlab.ListOptions{PerPage: constants.DEFAULT_PER_PAGE},
	}

	jobs, _, err := c.Jobs.ListPipelineJobs(projectID, int64(pipelineID), opts)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

// GetPipelineStatus retrieves the status of a specific pipeline.
//
// Parameters:
//   - owner: Repository owner (username or group)
//   - repo: Repository name
//   - pipelineID: Pipeline ID
//
// Returns an error if:
//   - The pipeline doesn't exist
//   - The API request fails
//
// Example:
//
//	pipeline, err := client.GetPipelineStatus("owner", "repo", 12345)
func (c *Client) GetPipelineStatus(owner, repo string, pipelineID int) (*gitlab.Pipeline, error) {
	projectID := owner + "/" + repo

	pipeline, _, err := c.Pipelines.GetPipeline(projectID, int64(pipelineID))
	if err != nil {
		return nil, fmt.Errorf("<?> Error: Failed to get pipeline status by pipelineID: %d.\n<?> Error: %w", pipelineID, err)
	}

	return pipeline, nil
}

// CancelPipeline cancels a running pipeline by ID.
//
// Parameters:
//   - owner: Repository owner (username or group)
//   - repo: Repository name
//   - pipelineID: Pipeline ID
//
// Returns an error if:
//   - The pipeline doesn't exist
//   - The API request fails
//
// Example:
//
//	err := client.CancelPipeline("owner", "repo", 12345)
func (c *Client) CancelPipeline(owner, repo string, pipelineID int) error {
	projectID := owner + "/" + repo

	_, _, err := c.Pipelines.CancelPipelineBuild(projectID, int64(pipelineID))
	if err != nil {
		return fmt.Errorf("<?> Error: Failed to cancel pipeline with pipelineID: %d.\n<?> Error: %w", pipelineID, err)
	}

	return nil
}

// GetPipelineSummary retrieves a summary of a pipeline.
//
// Parameters:
//   - owner: Repository owner (username or group)
//   - repo: Repository name
//   - pipelineID: Pipeline ID
//
// Returns an error if:
//   - The pipeline doesn't exist
//   - The API request fails
//
// Example:
//
//	summary, err := client.GetPipelineSummary("owner", "repo", 12345)
func (c *Client) GetPipelineSummary(owner, repo string, pipelineID int) (*helpers.PipelineSummary, error) {
	pipeline, err := c.GetPipelineStatus(owner, repo, pipelineID)
	if err != nil {
		return nil, err
	}

	summary := &helpers.PipelineSummary{
		ID:     pipeline.ID,
		Status: string(pipeline.Status),
		WebURL: pipeline.WebURL,
	}

	if pipeline.Ref != "" {
		summary.Ref = pipeline.Ref
	}

	if pipeline.CreatedAt != nil {
		summary.CreatedAt = pipeline.CreatedAt.String()
	}

	if pipeline.UpdatedAt != nil {
		summary.UpdatedAt = pipeline.UpdatedAt.String()
	}

	return summary, nil
}

// FilterJobsWithManualTriggerOnly filters jobs that have only:true trigger.
//
// Parameters:
//   - owner: Repository owner (username or group)
//   - repo: Repository name
//
// Returns an error if:
//   - The .gitlab-ci.yml doesn't exist
//   - The API request fails
//
// Example:
//
//	jobs, err := client.FilterJobsWithManualTriggerOnly("owner", "repo")
func (c *Client) FilterJobsWithManualTriggerOnly(owner, repo string) ([]string, error) {
	var jobsWithManualTrigger []string

	projectID := owner + "/" + repo

	file, _, err := c.RepositoryFiles.GetRawFile(projectID, ".gitlab-ci.yml", nil)
	if err != nil {
		return nil, fmt.Errorf("<?> Error: Failed to get .gitlab-ci.yml.\n<?> Error: %w", err)
	}

	var config map[string]interface{}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("<?> Error: Failed to parse .gitlab-ci.yml.\n<?> Error: %w", err)
	}

	for jobName, jobConfig := range config {
		// Skip reserved keywords
		if jobName == "stages" || jobName == "variables" || jobName == "default" || strings.HasPrefix(jobName, ".") {
			continue
		}

		if job, ok := jobConfig.(map[interface{}]interface{}); ok {
			if when, exists := job["when"]; exists && when == "manual" {
				jobsWithManualTrigger = append(jobsWithManualTrigger, jobName)
				fmt.Printf("✓ Found manual trigger in: %s\n", jobName)
			}
		}
	}

	return jobsWithManualTrigger, nil
}

// ListJobsWithManualTriggerOnly lists only jobs that have manual triggers.
//
// Parameters:
//   - owner: Repository owner (username or group)
//   - repo: Repository name
//
// Returns an error if:
//   - The .gitlab-ci.yml doesn't exist
//   - The API request fails
//
// Example:
//
//	jobs, err := client.ListJobsWithManualTriggerOnly("owner", "repo")
func (c *Client) ListJobsWithManualTriggerOnly(owner, repo string) ([]string, error) {
	return c.FilterJobsWithManualTriggerOnly(owner, repo)
}
