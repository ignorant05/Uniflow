package types

import (
	"fmt"
	"strings"
	"time"
)

// Repository descriptive information
type RepositoryInfo struct {
	Name          string
	FullName      string
	Description   string
	DefaultBranch string
	Private       bool
	HTMLURL       string
}

// Request Response Types
type TriggerRequest struct {
	// WorkflowName is the name or path of the workflow you want to trigger
	// Example: "deploy.yaml"
	WorkflowName string

	// Branch is the git reference to run on (branchn, tag, etc...)
	// Example: "main"
	Branch string

	// Inputs contains workflow specific parameters
	Inputs map[string]interface{}

	// Wait indicates whether to work for completion or not
	Wait bool

	// Timeout indicates the timeout you want to set (maximum waiting time)
	Timeout time.Duration
}

type TriggerResponse struct {
	// RunID is the unique identifier for this run
	RunID int64

	// RunNumber is the sequence run number
	RunNumber int

	// URL is the web-URL to view the run
	URL string

	// Status is the run status
	// Example: "Queued", "Success"
	Status string

	// QueuedAt is for when the run is queued
	QueuedAt time.Time
}

type StatusRequest struct {
	Name  string
	RunID int64
}

// Status represents the workflow's status
type Status struct {
	// RunID is the unique identifier for this run
	RunID int64

	// RunNumber is the sequence run number
	RunNumber int

	// Status is the run status
	// Example: "Queued", "Success"
	Status string

	// Conclusion is the final result: "success", "failure", "cancelled", "skipped"
	// Only set when Status is "completed" else, it's an empty string
	Conclusion string

	// StartedAt is when the workflow is executed
	StartedAt time.Time

	// CompletedAt is when the workflow's execution completed
	CompletedAt time.Time

	// Duration is how much time did the workflow last during execution
	Duration time.Duration

	// URL is the web-URL to view the run
	URL string

	// QueuedAt is for when the run is queued
	QueuedAt time.Time

	// Metadata contains platform-specific additional information
	Metadata map[string]interface{}
}

// Run represents a workflow run in a list.
type Run struct {
	RunID       int64
	RunNumber   int
	Status      string
	Conclusion  string
	Branch      string
	Actor       string
	Event       string
	CommitSHA   string
	TriggeredBy string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	URL         string
}

type ListWorkflowsRequest struct {
	// WithDispatch is for whether to list only workflows containing "workflow_dispatch" trigger or all
	WithDispatch bool
}

// ListRunsRequest contains parameters for listing runs.
type ListWorkflowRunsRequest struct {
	RunID int64
	// WorkflowName is the name or path of the workflow you want to trigger
	// Example: "deploy.yaml"
	WorkflowName string

	// Status filters by status: "queued", "in_progress", "completed" (optional)
	Status string

	// Branch is the git reference to run on (branchn, tag, etc...)
	// Example: "main"
	Branch string

	// Limit is the maximum number to return of recent runs
	Limit int
}

type ListWokflowJobsRequest struct {
	// WorkflowName is the name or path of the workflow you want to trigger
	// Example: "deploy.yaml"
	WorkflowName string

	// Status filters by status: "queued", "in_progress", "completed" (optional)
	Status string

	// Branch is the git reference to run on (branchn, tag, etc...)
	// Example: "main"
	Branch string
}

// Workflow summary
type WorkflowRunSummary struct {
	ID         int64
	Name       string
	Status     string
	Conclusion string
	CreatedAt  string
	UpdatedAt  string
	HTMLURL    string
}

type LogsStreamRequest struct {
	// RunID is the unique identifier for this run
	RunID int64

	Follow bool

	NoColor bool

	// Tail retusn only the N lines of logs
	// Example: "tail 10" returns only the last 10 logs
	Tail int
}

type LogsRequest struct {
	// RunID is the unique identifier for this run
	RunID int64

	WorkflowName string

	DownloadPath string

	// Tail retusn only the N lines of logs
	// Example: "tail 10" returns only the last 10 logs
	Tail int
}

type LogsResponse struct {
	// Content is the log content
	Content string

	// URL is the web-URL to view logs
	URL string
}

type StreamLogsRequest struct {
	// RunID is the workflow run identifier
	RunID int64

	// JobName filters logs to a specific job (optional)
	JobName string

	// Follow continues streaming until run completes
	Follow bool
}

// LogCallback is called for each new log line during streaming
type LogCallback func(line *LogLine) error

type LogLine struct {
	// Content is the log message
	Content string

	// Timestamp is when the line was generated
	Timestamp time.Time

	// JobName is the job that generated this line
	JobName string

	// Level is the log level (info, error, etc...)
	Level string
}

// Workflow represents an available workflow/pipeline/job.
type Workflow struct {
	// RunID is the workflow run identifier
	ID int64

	// Name is the workflow name
	Name string

	// Path is the wofkflow file path
	Path string

	// State is whether the workflow is active or disabled
	State string

	// URL is the web-URL to view the workflow
	URL string

	// WithDispatch is for whether to list only workflows containing "workflow_dispatch" trigger or all
	WithDispatch bool
}

// WorkflowJob represents a repository action workflow job
type WorkflowJob struct {
	ID           int64
	RunID        int64
	WorkflowName string
	Name         string
	Status       string
	Conclusion   string
	RunURL       string
	URL          string
	HTMLURL      string
}

var (
	// ErrUnauthorized indicates authentication failure error
	ErrUnauthorized = &PlatformError{Code: "unauthorized", Message: "authentication failure"}

	// ErrNotFound indicates resource not found error
	ErrNotFound = &PlatformError{Code: "not_found", Message: "resource not found"}

	// ErrForbidden indicates permission failure for operation
	ErrForbidden = &PlatformError{Code: "forbidden", Message: "not permitted"}

	// ErrRateLimitExceeded indicates rate limit exceeded error
	ErrRateLimitExceeded = &PlatformError{Code: "rate_limited", Message: "rate limit exceeded"}

	// ErrTimeout indicated request timed out error
	ErrTimeout = &PlatformError{Code: "timeout", Message: "request timed out"}
)

// PlatformError is a standardized error across all platforms.
type PlatformError struct {
	Code       string
	Message    string
	StatusCode int
	Platform   string
	Details    map[string]interface{}
}

func (e *PlatformError) Error() string {
	if e.Platform != "" {
		return fmt.Sprintf("[%s] %s: %s\n", e.Platform, e.Code, e.Message)
	}
	return fmt.Sprintf("%s: %s\n", e.Code, e.Message)
}

// IsRunning checks if a status indicates the run is still active
func IsRunning(status string) bool {
	return strings.ToLower(status) == "queued" ||
		strings.ToLower(status) == "waiting" ||
		strings.ToLower(status) == "in_progress"
}

// IsRunning checks if a status indicates the completed
func IsCompleted(status string) bool {
	return strings.ToLower(status) == "completed"
}

// IsRunning checks if a status indicates the run failed
func IsFailed(status string) bool {
	return strings.ToLower(status) == "failed"
}

// IsRunning checks if a status indicates the run is completed successfully
func IsSuccess(status string) bool {
	return strings.ToLower(status) == "success"
}
