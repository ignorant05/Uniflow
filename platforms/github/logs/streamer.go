package github

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	gh "github.com/google/go-github/v57/github"

	"github.com/ignorant05/Uniflow/internal/helpers"
	"github.com/ignorant05/Uniflow/platforms/github"
	"github.com/ignorant05/Uniflow/platforms/github/constants"
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
	ShowTime  bool
	Colorize  bool
}

// Streamer struct
type Streamer struct {
	client     *github.Client
	ctx        context.Context
	cancelFunc context.CancelFunc
	owner      string
	repo       string
	runID      int64

	follow    bool
	tailLines int
	showTime  bool
	colorize  bool

	lastLogLine int
	seenJobs    map[int64]bool
}

// NewStreamer creates new streamer
//
// Parameters:
//   - client: github client
//   - owner: owner name
//   - repo: repository name
//   - runID: workflow runID
//   - opts: Streamer options
func NewStreamer(client *github.Client, owner, repo string, runID int64, opts StreamerOptions) *Streamer {
	ctx, cancel := context.WithCancel(context.Background())

	return &Streamer{
		client:     client,
		ctx:        ctx,
		cancelFunc: cancel,
		owner:      owner,
		repo:       repo,
		runID:      runID,
		follow:     opts.Follow,
		tailLines:  opts.TailLines,
		showTime:   opts.ShowTime,
		colorize:   opts.Colorize,
		seenJobs:   make(map[int64]bool),
	}
}

// Stream function streamns workflow (search by ID)
func (s *Streamer) Stream() error {
	run, _, err := s.client.Actions.GetWorkflowRunByID(s.ctx, s.owner, s.repo, s.runID)
	if err != nil {
		return fmt.Errorf("<?> Error: Failed to get workflow run by ID: %d", s.runID)
	}

	s.printHeader(run)

	if s.follow {
		return s.streamWithFollow(run)
	}

	return s.streamOnce(run)
}

// streamOnce displays logs
//
// Parameters :
//   - run: github workflow run
//
// Errors possible causes:
//   - no jobs for this workflow
//
// Examples:
// err := s.streamOnce(run)
func (s *Streamer) streamOnce(run *gh.WorkflowRun) error {
	if run.GetStatus() == "completed" {
		fmt.Println("  Waiting for workflow to complete.")

		for range constants.MaxPollAttempts {
			select {
			case <-s.ctx.Done():
				return nil
			case <-time.After(constants.PollInterval):
				var err error
				run, _, err := s.client.Actions.GetWorkflowRunByID(s.ctx, s.owner, s.repo, s.runID)
				if err != nil {
					return err
				}
				if run.GetStatus() == "completed" {
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
//   - run: github workflow run
//
// Errors possible causes:
//   - no jobs for this workflow
//
// Examples:
// err := s.streamWithFollow(run)
func (s *Streamer) streamWithFollow(run *gh.WorkflowRun) error {
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
			currentRun, _, err := s.client.Actions.GetWorkflowRunByID(s.ctx, s.owner, s.repo, run.GetID())
			if err != nil {
				return err
			}

			jobs, _, err := s.client.Actions.ListWorkflowJobs(s.ctx, s.owner, s.repo, run.GetID(), nil)
			if err != nil {
				return err
			}

			for _, job := range jobs.Jobs {
				if err := s.streamJobLogs(job, seenLines); err != nil {
					if s.colorize {
						color.Red("<?> Warning: Failed to get logs for job %s: %v", job.GetName(), err)
					} else {
						fmt.Printf("<?> Warning: Failed to get logs for job %s: %v\n", job.GetName(), err)
					}
				}
			}

			if currentRun.GetStatus() == "completed" {
				s.formatCompletion(currentRun)
				return nil
			}
		}

	}
}

// streamJobLogs fetches and displays workflow job logs in a formatted and colorized manner
//
// Parameters :
//   - job: workflow job
//   - seenLines: traversed log lines
//
// Errors possible causes:
//   - can't read logs
//
// Examples:
// err := s.streamJobLogs(run)
func (s *Streamer) streamJobLogs(job *gh.WorkflowJob, seenLines map[string]bool) error {
	if job.GetStatus() == "queued" ||
		job.GetStatus() == "waiting" {
		return nil
	}

	if !s.seenJobs[job.GetID()] {
		s.printJobHeader(job)
		s.seenJobs[job.GetID()] = true
	}

	logURL, _, err := s.client.Actions.GetWorkflowJobLogs(s.ctx, s.owner, s.repo, job.GetID(), constants.MAX_REDIRECTS)
	if err != nil {
		return err
	}

	logs, err := s.readLogs(logURL.String())
	if err != nil {
		return err
	}

	lines := strings.SplitSeq(logs, "\n")
	for line := range lines {
		if line == "" {
			continue
		}

		lineKey := fmt.Sprintf("%d:%s", job.GetID(), line)
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
	jobs, _, err := s.client.Actions.ListWorkflowJobs(s.ctx, s.owner, s.repo, s.runID, nil)
	if err != nil {
		return err
	}

	for _, job := range jobs.Jobs {
		s.printJobHeader(job)

		logURL, _, err := s.client.Actions.GetWorkflowJobLogs(s.ctx, s.owner, s.repo, job.GetID(), constants.MAX_REDIRECTS)
		if err != nil {
			if s.colorize {
				color.Yellow("	No logs available for this job.")
			} else {
				fmt.Println("  No logs available for this job.")
			}

			fmt.Println()
			continue
		}

		logs, err := s.readLogs(logURL.String())
		if err != nil {
			return err
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

// printLogLine prints colorized log line
//
// Parameters :
//   - line: log line
//
// Examples:
// s.printLogLine(dine)
func (s *Streamer) printLogLine(line string) {
	var output string

	timestamp, content := helpers.FormatLogs(line)

	level := s.detectLogLevel(content)

	if s.showTime && timestamp != "" {
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

// colorizeContent colorizez content depending on log level
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

// printJobHeader prints workflow job header
//
// Parameters :
//   - job: workflow job
//
// Examples:
// s.printJobHeader(job)
func (s *Streamer) printJobHeader(job *gh.WorkflowJob) {
	if s.colorize {
		color.New(color.Bold, color.FgCyan).Printf("\nJob: %s\n", job.GetName())
	} else {
		fmt.Printf("\nJob: %s\n", job.GetName())
	}

	fmt.Printf("\nJob: %s\n", s.formatStatus(job.GetStatus()))
	if job.GetConclusion() != "" {
		fmt.Printf("\nResult: %s\n", s.formatConclusion(job.GetConclusion()))
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
//   - run: github workflow run
//
// Examples:
// s.printHeader(run)
func (s *Streamer) printHeader(run *gh.WorkflowRun) {
	if s.colorize {
		color.New(color.Bold).Println("Workflow Run")
	} else {
		fmt.Println("Workflow Run")
	}

	fmt.Println(strings.Repeat("─", 80))
	fmt.Printf("Name:       %s\n", run.GetName())
	fmt.Printf("Run #:      %d\n", run.GetRunNumber())
	fmt.Printf("Status:     %s\n", s.formatStatus(run.GetStatus()))
	fmt.Printf("Branch:     %s\n", run.GetHeadBranch())
	fmt.Printf("Commit:     %.7s\n", run.GetHeadSHA())
	fmt.Printf("Actor:      %s\n", run.GetActor().GetLogin())
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()
}

// formatStatus formats status of workflow
//
// Parameters :
//   - status: status string
//
// Examples:
// status := s.formatStatus("completed")
func (s *Streamer) formatStatus(status string) string {
	switch status {
	case "queued":
		return "Queued"
	case "in_progress":
		return "In Progress"
	case "completed":
		return "Completed"
	case "waiting":
		return "Waiting"
	default:
		return status
	}
}

// formatConclusion formats conclusion of workflow run
//
// Parameters :
//   - conc: conclusion content
//
// Examples:
// conclusion := s.formatConclusion("success")
func (s *Streamer) formatConclusion(conc string) string {
	switch conc {
	case "success":
		return "Success"
	case "failure":
		return "Failure"
	case "cancelled":
		return "Cancelled"
	case "skipped":
		return "Skipped"
	default:
		return conc
	}
}

// formatCompletion formats completion process
//
// Parameters :
//   - run: github workflow run
//
// Examples:
// s.formatCompletion(run)
func (s *Streamer) formatCompletion(run *gh.WorkflowRun) {
	fmt.Println()
	fmt.Println(strings.Repeat("-", 80))
	conclusion := run.GetConclusion()

	if s.colorize {
		switch conclusion {
		case "success":
			color.Green("Workflow completed successfully.")
		case "failure":
			color.Red("Workflow failed.")
		case "cancelled":
			color.Yellow("Workflow canceled.")
		default:
			color.White("Workflow completed: %s\n", conclusion)
		}
	} else {
		fmt.Printf("Workflow completed: %s\n", conclusion)
	}

	fmt.Println(strings.Repeat("-", 80))
}

// Stop cancels gracefully
//
// Examples:
// s.stop()
func (s *Streamer) Stop() {
	s.cancelFunc()
}
