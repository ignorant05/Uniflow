package gitlab

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/fatih/color"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/ignorant05/Uniflow/internal/helpers"
	"github.com/ignorant05/Uniflow/platforms/configurations/gitlab/constants"
)

// Log Level type
type LogLevel int

// Log level types
const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelError
	LevelWarning
	LevelSuccess
)

// Streaming options struct
type StreamerOptions struct {
	Follow    bool
	TailLines int
	Colorize  bool
}

// Streamer struct
type Streamer struct {
	client     *gitlab.Client
	ctx        context.Context
	cancelFunc context.CancelFunc
	projectID  string
	pipelineID int

	follow    bool
	tailLines int
	colorize  bool

	seenJobs map[int]bool
}

// NewStreamer creates new streamer
//
// Parameters:
//   - client: gitlab client
//   - projectID: project ID in format "owner/repo"
//   - pipelineID: pipeline ID
//   - opts: Streamer options
func NewStreamer(client *gitlab.Client, projectID string, pipelineID int, opts StreamerOptions) *Streamer {
	ctx, cancel := context.WithCancel(context.Background())

	return &Streamer{
		client:     client,
		ctx:        ctx,
		cancelFunc: cancel,
		projectID:  projectID,
		pipelineID: pipelineID,
		follow:     opts.Follow,
		tailLines:  opts.TailLines,
		colorize:   opts.Colorize,
		seenJobs:   make(map[int]bool),
	}
}

// Stream function streams pipeline (search by ID)
func (s *Streamer) Stream() error {
	pipeline, _, err := s.client.Pipelines.GetPipeline(s.projectID, int64(s.pipelineID))
	if err != nil {
		return fmt.Errorf("<?> Error: Failed to get pipeline by ID: %d", int64(s.pipelineID))
	}

	s.printHeader(pipeline)

	if s.follow {
		return s.streamWithFollow(pipeline)
	}

	return s.streamOnce(pipeline)
}

// streamOnce displays logs
//
// Parameters :
//   - pipeline: gitlab pipeline
//
// Errors possible causes:
//   - no jobs for this pipeline
//
// Examples:
// err := s.streamOnce(pipeline)
func (s *Streamer) streamOnce(pipeline *gitlab.Pipeline) error {
	if pipeline.Status != "running" && pipeline.Status != "pending" {
		fmt.Println("  Waiting for pipeline to complete.")

		for range constants.MaxPollAttempts {
			select {
			case <-s.ctx.Done():
				return nil
			case <-time.After(constants.PollInterval):
				var err error
				pipeline, _, err := s.client.Pipelines.GetPipeline(s.projectID, int64(s.pipelineID))
				if err != nil {
					return err
				}
				if pipeline.Status != "running" && pipeline.Status != "pending" {
					return s.fetchAndDisplayLogs()
				}
			}
		}
	}
	return s.fetchAndDisplayLogs()
}

// streamWithFollow displays logs in real time
//
// Parameters :
//   - pipeline: gitlab pipeline
//
// Errors possible causes:
//   - no jobs for this pipeline
//
// Examples:
// err := s.streamWithFollow(pipeline)
func (s *Streamer) streamWithFollow(pipeline *gitlab.Pipeline) error {
	fmt.Println("  Following logs (press Ctrl+C to stop)...")

	ticker := time.NewTicker(constants.PollInterval)
	defer ticker.Stop()

	seenLines := make(map[string]bool)

	for {
		select {
		case <-s.ctx.Done():
			fmt.Println("\n\n  Log streaming stopped.")
			return nil

		case <-ticker.C:
			currentPipeline, _, err := s.client.Pipelines.GetPipeline(s.projectID, int64(s.pipelineID))
			if err != nil {
				return err
			}

			jobs, _, err := s.client.Jobs.ListPipelineJobs(s.projectID, int64(s.pipelineID), nil)
			if err != nil {
				return err
			}

			for _, job := range jobs {
				if err := s.streamJobLogs(job, seenLines); err != nil {
					if s.colorize {
						color.Red("<?> Warning: Failed to get logs for job %s: %v", job.Name, err)
					} else {
						fmt.Printf("<?> Warning: Failed to get logs for job %s: %v\n", job.Name, err)
					}
				}
			}

			if currentPipeline.Status != "running" && currentPipeline.Status != "pending" {
				s.formatCompletion(currentPipeline)
				return nil
			}
		}

	}
}

// streamJobLogs fetches and displays pipeline job logs in a formatted and colorized manner
//
// Parameters :
//   - job: pipeline job
//   - seenLines: traversed log lines
//
// Errors possible causes:
//   - can't read logs
//
// Examples:
// err := s.streamJobLogs(job, seenLines)
func (s *Streamer) streamJobLogs(job *gitlab.Job, seenLines map[string]bool) error {
	if job.Status == "queued" || job.Status == "waiting_for_resource" || job.Status == "preparing" {
		return nil
	}

	intJobID := int(job.ID)

	if !s.seenJobs[intJobID] {
		s.printJobHeader(job)
		s.seenJobs[intJobID] = true
	}

	logs, err := s.readJobLogs(intJobID)
	if err != nil {
		return err
	}

	lines := strings.SplitSeq(logs, "\n")
	for line := range lines {
		if line == "" {
			continue
		}

		lineKey := fmt.Sprintf("%d:%s", job.ID, line)
		if seenLines[lineKey] {
			continue
		}
		seenLines[lineKey] = true

		s.printLogLine(line)
	}

	return nil
}

// fetchAndDisplayLogs fetch and displays formatted and colorized logs
//
// Parameters :
//   - None
//
// Errors possible causes:
//   - can't read logs
//
// Examples:
// err := s.fetchAndDisplayLogs()
func (s *Streamer) fetchAndDisplayLogs() error {
	jobs, _, err := s.client.Jobs.ListPipelineJobs(s.projectID, int64(s.pipelineID), nil)
	if err != nil {
		return err
	}

	for _, job := range jobs {
		s.printJobHeader(job)

		logs, err := s.readJobLogs(int(job.ID))
		if err != nil {
			if s.colorize {
				color.Yellow("	No logs available for this job.")
			} else {
				fmt.Println("  No logs available for this job.")
			}

			fmt.Println()
			continue
		}

		if s.tailLines > 0 {
			logs = s.applyTail(logs)
		}

		lines := strings.SplitSeq(logs, "\n")
		for line := range lines {
			if line != "" {
				s.printLogLine(line)
			}
		}

		fmt.Println()

	}

	return nil
}

// readJobLogs reads logs from a GitLab job
//
// Parameters:
//   - jobID: Job ID
//
// Returns:
//   - logs as string
//   - error if job logs can't be retrieved
func (s *Streamer) readJobLogs(jobID int) (string, error) {
	logs, _, err := s.client.Jobs.GetTraceFile(s.projectID, int64(jobID))
	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(logs)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// printLogLine prints colorized log line
//
// Parameters :
//   - line: log line
//
// Examples:
// s.printLogLine(line)
func (s *Streamer) printLogLine(line string) {
	var output string

	timestamp, content := helpers.FormatLogs(line)

	level := s.detectLogLevel(content)

	if timestamp != "" {
		timeStr := s.formatTimestamp(timestamp)
		if s.colorize {
			output = color.New(color.FgHiBlack).Sprint(timeStr) + " "
		} else {
			output = timeStr + " "
		}
	}

	if s.colorize {
		output += s.colorizeContent(content, level)
	} else {
		output = content
	}

	fmt.Println(output)
}

// colorizeContent colorizes content depending on log level
//
// Parameters :
//   - content: logs content
//   - loglvl: log level
//
// Examples:
// colorizedContent := s.colorizeContent(content, lvl)
func (s *Streamer) colorizeContent(content string, loglvl LogLevel) string {
	switch loglvl {
	case LevelError:
		return color.RedString(content)
	case LevelWarning:
		return color.YellowString(content)
	case LevelSuccess:
		return color.GreenString(content)
	case LevelDebug:
		return color.New(color.FgHiBlack).Sprint(content)
	default:
		return content
	}
}

// detectLogLevel detects log level to format
//
// Parameters :
//   - content: log content
//
// Examples:
// lvl := s.detectLogLevel(content)
func (s *Streamer) detectLogLevel(content string) LogLevel {
	contentLower := strings.ToLower(content)

	if helpers.IsError(contentLower) {
		return LevelError
	}

	if helpers.IsDebug(contentLower) {
		return LevelDebug
	}

	if helpers.IsWarning(contentLower) {
		return LevelWarning
	}

	if helpers.IsSuccess(contentLower) {
		return LevelSuccess
	}

	return LevelInfo
}

// formatTimestamp formats timestamp in "15:04:05" format
//
// Parameters :
//   - timestamp: time string
//
// Examples:
// timestamp := s.formatTimestamp(time)
func (s *Streamer) formatTimestamp(timestamp string) string {
	t, err := time.Parse(time.RFC3339Nano, timestamp)
	if err != nil {
		return timestamp
	}

	return t.Format("15:04:05")
}

// printJobHeader prints pipeline job header
//
// Parameters :
//   - job: pipeline job
//
// Examples:
// s.printJobHeader(job)
func (s *Streamer) printJobHeader(job *gitlab.Job) {
	if s.colorize {
		_, err := color.New(color.Bold, color.FgCyan).Printf("\nJob: %s\n", job.Name)
		if err != nil {
			fmt.Printf("<?> Error: %v\n", err)
			return
		}
	} else {
		fmt.Printf("\nJob: %s\n", job.Name)
	}

	fmt.Printf("Status: %s\n", s.formatStatus(job.Status))
	if job.FailureReason != "" {
		fmt.Printf("Failure Reason: %s\n", job.FailureReason)
	}

	fmt.Println(strings.Repeat("─", 80))
}

// applyTail cuts the lines > s.taillines
//
// Parameters :
//   - logs: logs data
//
// Examples:
// tailedLogs := s.applyTail(logs)
func (s *Streamer) applyTail(logs string) string {
	lines := strings.Split(logs, "\n")
	if len(lines) <= s.tailLines {
		return logs
	}

	return strings.Join(lines[len(lines)-s.tailLines:], "\n")
}

// printHeader prints header
//
// Parameters :
//   - pipeline: gitlab pipeline
//
// Examples:
// s.printHeader(pipeline)
func (s *Streamer) printHeader(pipeline *gitlab.Pipeline) {
	if s.colorize {
		_, err := color.New(color.Bold).Println("Pipeline")
		if err != nil {
			fmt.Printf("<?> Error: %v\n", err)
			return
		}
	} else {
		fmt.Println("Pipeline")
	}

	fmt.Println(strings.Repeat("─", 80))
	fmt.Printf("ID:         %d\n", pipeline.ID)
	fmt.Printf("Status:     %s\n", s.formatStatus(pipeline.Status))
	fmt.Printf("Ref:        %s\n", pipeline.Ref)
	fmt.Printf("Commit:     %.7s\n", pipeline.SHA)
	fmt.Printf("Web URL:    %s\n", pipeline.WebURL)
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()
}

// formatStatus formats status of pipeline
//
// Parameters :
//   - status: status string
//
// Examples:
// status := s.formatStatus("running")
func (s *Streamer) formatStatus(status string) string {
	switch status {
	case "created":
		return "Created"
	case "pending":
		return "Pending"
	case "running":
		return "Running"
	case "success":
		return "Success"
	case "failed":
		return "Failed"
	case "canceled":
		return "Canceled"
	case "skipped":
		return "Skipped"
	default:
		return status
	}
}

// formatCompletion formats completion process
//
// Parameters :
//   - pipeline: gitlab pipeline
//
// Examples:
// s.formatCompletion(pipeline)
func (s *Streamer) formatCompletion(pipeline *gitlab.Pipeline) {
	fmt.Println()
	fmt.Println(strings.Repeat("-", 80))

	if s.colorize {
		switch pipeline.Status {
		case "success":
			color.Green("Pipeline completed successfully.")
		case "failed":
			color.Red("Pipeline failed.")
		case "canceled":
			color.Yellow("Pipeline canceled.")
		default:
			color.White("Pipeline completed: %s\n", pipeline.Status)
		}
	} else {
		fmt.Printf("Pipeline completed: %s\n", pipeline.Status)
	}

	fmt.Println(strings.Repeat("-", 80))
}

// Stop cancels gracefully
//
// Examples:
// s.Stop()
func (s *Streamer) Stop() {
	s.cancelFunc()
}
