package jenkins

import (
	"fmt"

	"github.com/bndr/gojenkins"
)

// GetDefaultJob retrieves the default job name from configuration.
//
// Jenkins doesn't have "repositories" like GitHub - it has "jobs".
// This method provides similar interface to GitHub client.
//
// Parameters:
//   - None
//
// Returns an error if:
//   - no default job is configured
//
// Example:
//
//	jobName, err := client.GetDefaultJob()
func (c *Client) GetDefaultJob() (string, error) {
	if c.Config.JobName == "" {
		return "", fmt.Errorf("no jenkins default job configured")
	}

	return c.Config.JobName, nil
}

// GetJobInfo retrieves information about a specific job.
//
// Parameters:
//   - jobName: name of the Jenkins job
//
// Returns:
//   - job information
//   - error if any
//
// Example:
//
//	job, err := client.GetJobInfo("my-pipeline")
func (c *Client) GetJobInfo(jobName string) (*gojenkins.Job, error) {
	if jobName == "" {
		// Try to get default job
		defaultJob, err := c.GetDefaultJob()
		if err != nil {
			return nil, fmt.Errorf("no job specified and no default job configured,\nerror: %w", err)
		}
		jobName = defaultJob
	}
	return c.Client.GetJob(c.Ctx, jobName)
}

// TriggerJob triggers a Jenkins job with parameters.
//
// Parameters:
//   - jobName: name of the Jenkins job
//   - params: map of parameters for the job
//
// Example:
//
//	queueID, err := client.TriggerJob("default", params)
func (c *Client) TriggerJob(jobName string, params map[string]interface{}) (int64, error) {
	if jobName == "" {
		return 0, fmt.Errorf("no default job found with this name: %s", jobName)
	}

	stringParams := make(map[string]string)
	for key, val := range params {
		if strVal, ok := val.(string); ok {
			stringParams[key] = strVal
		} else {
			stringParams[key] = fmt.Sprintf("%v", val)
		}
	}

	return c.Client.BuildJob(c.Ctx, jobName, stringParams)
}

// TriggerDefaultJob triggers Jenkins default job with parameters.
//
// Parameters:
//   - params: map of parameters for the job
//
// Example:
//
//	queueID, err := client.TriggerDefaultJob(params)
func (c *Client) TriggerDefaultJob(params map[string]interface{}) (int64, error) {
	if c.Config.JobName == "" {
		return 0, fmt.Errorf("no default job found, please consider updating your defaults")
	}

	stringParams := make(map[string]string)
	for key, val := range params {
		if strVal, ok := val.(string); ok {
			stringParams[key] = strVal
		} else {
			stringParams[key] = fmt.Sprintf("%v", val)
		}
	}

	return c.Client.BuildJob(c.Ctx, c.Config.JobName, stringParams)
}

// ListJobs lists all Jenkins jobs present.
//
// Parameters:
//   - None
//
// Example:
//
//	jobs, err := client.ListJobs()
func (c *Client) ListJobs() ([]*gojenkins.Job, error) {
	return c.Client.GetAllJobs(c.Ctx)
}

// GetJobBuilds lists all job builds.
//
// Parameters:
//   - jobName: target job name
//
// Returns an error if:
//   - The job doesn't exist
//   - The API request fails
//
// Example:
//
//	runs, err := client.GetJobBuilds("deploy-prod")
func (c *Client) GetJobBuilds(jobName string) ([]*gojenkins.Build, error) {
	job, err := c.Client.GetJob(c.Ctx, jobName)
	if err != nil {
		return nil, fmt.Errorf("failed to get job by job name: %s\nerror: %w", jobName, err)
	}

	buildsIDs, err := job.GetAllBuildIds(c.Ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get job builds IDs by job name: %s\nerror: %w", jobName, err)
	}

	builds := make([]*gojenkins.Build, 0, len(buildsIDs))
	for _, buildID := range buildsIDs {
		build, err := c.Client.GetBuild(c.Ctx, jobName, buildID.Number)
		if err != nil {
			fmt.Printf("failed to get job builds IDs by job name: %s\nerror: %v", jobName, err)
			continue
		}
		builds = append(builds, build)
	}

	return builds, nil
}

// GetJobBuildInfo retrieves job build information.
//
// Parameters:
//   - jobName: target job name
//   - buildID: target build ID
//
// Returns an error if:
//   - The job doesn't exist
//   - Invalid buildID
//
// Example:
//
//	build, err := client.GetJobBuild("deploy-prod", 12334)
func (c *Client) GetJobBuild(jobName string, buildID int64) (*gojenkins.Build, error) {
	build, err := c.Client.GetBuild(c.Ctx, jobName, buildID)
	if err != nil {
		return nil, fmt.Errorf("failed to get build by job name: %s and build ID: %d\nerror: %w", jobName, buildID, err)
	}

	return build, nil
}

// GetBuildLogs retrieves the console output/logs of a build with fresh data.
//
// Parameters:
//   - jobName: target job name
//   - buildID: target build ID
//
// Returns an error if:
//   - The job doesn't exist
//   - Invalid buildID
//   - The API request fails
//
// Example:
//
//	logs, err := client.GetBuildLogs("deploy-prod", 12334)
func (c *Client) GetBuildLogs(jobName string, buildID int64) (string, error) {
	build, err := c.GetJobBuild(jobName, buildID)
	if err != nil {
		return "", err
	}

	_, err = build.Poll(c.Ctx)
	if err != nil {
		return "", fmt.Errorf("failed to poll build data for job: %s and build ID: %d\nerror: %w", jobName, buildID, err)
	}

	return build.GetConsoleOutput(c.Ctx), nil
}

// GetBuildStatus retrieves the current status of a build with fresh data.
//
// Parameters:
//   - jobName: target job name
//   - buildID: target build ID
//
// Returns status string and an error if:
//   - The job doesn't exist
//   - Invalid buildID
//   - The API request fails
//
// Example:
//
//	status, err := client.GetBuildStatus("deploy-prod", 12334)
func (c *Client) GetBuildStatus(jobName string, buildID int64) (string, error) {
	build, err := c.Client.GetBuild(c.Ctx, jobName, buildID)
	if err != nil {
		return "", fmt.Errorf("failed to get build by job name: %s and build ID: %d\nerror: %w", jobName, buildID, err)
	}

	_, err = build.Poll(c.Ctx)
	if err != nil {
		return "", fmt.Errorf("failed to poll build data for job: %s and build ID: %d\nerror: %w", jobName, buildID, err)
	}

	return build.GetResult(), nil
}

// IsJobRunning checks if a build is currently running.
//
// Parameters:
//   - jobName: target job name
//   - buildID: target build ID
//
// Returns boolean indicating if the build is running and an error if:
//   - The job doesn't exist
//   - Invalid buildID
//   - The API request fails
//
// Example:
//
//	running, err := client.IsJobRunning("deploy-prod", 12334)
func (c *Client) IsJobRunning(jobName string, buildID int64) (bool, error) {
	build, err := c.Client.GetBuild(c.Ctx, jobName, buildID)
	if err != nil {
		return false, fmt.Errorf("failed to get build by job name: %s and build ID: %d\nerror: %w", jobName, buildID, err)
	}

	_, err = build.Poll(c.Ctx)
	if err != nil {
		return false, fmt.Errorf("failed to poll build data for job: %s and build ID: %d\nerror: %w", jobName, buildID, err)
	}

	return build.IsRunning(c.Ctx), nil
}

// ListAllNodes lists all Jenkins nodes/agents.
//
// Parameters:
//   - None
//
// Example:
//
//	nodes, err := client.ListAllNodes()
func (c *Client) ListAllNodes() ([]*gojenkins.Node, error) {
	nodes, err := c.Client.GetAllNodes(c.Ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Nodes.\nerror: %w", err)
	}

	return nodes, nil
}

// GetNodeStatus returns the status of a given node (by name)
// Parameters:
//   - nodeName: name of the node
//
// Returns an error if:
//   - The node doesn't exist
//   - The API request fails
//
// Example:
//
//	online, err := client.GetNodeStatus("agent-1")
func (c *Client) GetNodeStatus(nodeName string) (bool, error) {
	node, err := c.Client.GetNode(c.Ctx, nodeName)
	if err != nil {
		return false, fmt.Errorf("failed to get node: %s\nerror: %w", nodeName, err)
	}

	_, err = node.Poll(c.Ctx)
	if err != nil {
		return false, fmt.Errorf("failed to poll node data: %s\nerror: %w", nodeName, err)
	}

	return node.IsOnline(c.Ctx)
}

// GetCurrentQueue retrieves all tasks currently in the Jenkins queue.
//
// Parameters:
//   - None
//
// Example:
//
//	queue, err := client.GetCurrentQueue()
func (c *Client) GetCurrentQueue() (*gojenkins.Queue, error) {
	tasks, err := c.Client.GetQueue(c.Ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Jenkins queue.\nerror %w", err)
	}

	return tasks, nil
}

// CancelQueueTask cancels a task in the queue by its ID.
//
// Parameters:
//   - jobName: target job name
//   - taskID: ID of the task/build in queue
//
// Returns an error if:
//   - The task doesn't exist
//   - The API request fails
//
// Example:
//
//	err := client.CancelExecution("build-prod", 12345)
func (c *Client) CancelExecution(jobName string, executionID int64) error {
	tasks, err := c.GetCurrentQueue()
	if err != nil {
		return err
	}

	task := tasks.GetTaskById(executionID)
	if ok, err := task.Cancel(c.Ctx); ok && err == nil {
		fmt.Printf("✓ Cancelled queued task %d for job %s", executionID, jobName)
		return nil
	}

	build, _ := c.Client.GetBuild(c.Ctx, jobName, executionID)
	if ok, err := build.Stop(c.Ctx); ok && err == nil {
		fmt.Printf("✓ Cancelled running build %d for job %s", executionID, jobName)
		return nil
	}

	build, buildErr := c.Client.GetBuild(c.Ctx, jobName, executionID)
	if buildErr == nil && !build.Raw.Building {
		return fmt.Errorf("build %d for job %s already completed with status: %s",
			executionID, jobName, build.GetResult())
	}

	return fmt.Errorf("failed to cancel execution %d for job %s", executionID, jobName)
}
