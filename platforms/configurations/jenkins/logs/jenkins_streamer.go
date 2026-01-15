package jenkins

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bndr/gojenkins"
	"github.com/fatih/color"

	"github.com/ignorant05/Uniflow/internal/helpers"
	"github.com/ignorant05/Uniflow/platforms/configurations/jenkins/constants"
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
	ShowTime  bool
}

// Streamer struct
type Streamer struct {
	Client     *gojenkins.Jenkins
	ctx        context.Context
	cancelFunc context.CancelFunc
	jobName    string
	buildID    int64

	follow    bool
	tailLines int
	colorize  bool
	showTime  bool

	seenLogs map[string]bool
}

// NewStreamer creates new streamer
//
// Parameters:
//   - client: jenkins client
//   - jobName: Jenkins job name
//   - buildID: Jenkins build ID
//   - opts: Streamer options
func NewStreamer(client *gojenkins.Jenkins, jobName string, buildID int64, opts StreamerOptions) *Streamer {
	ctx, cancel := context.WithCancel(context.Background())

	return &Streamer{
		Client:     client,
		ctx:        ctx,
		cancelFunc: cancel,
		jobName:    jobName,
		buildID:    buildID,
		follow:     opts.Follow,
		tailLines:  opts.TailLines,
		colorize:   opts.Colorize,
		showTime:   opts.ShowTime,
		seenLogs:   make(map[string]bool),
	}
}

// Stream function streams build logs (search by Job Name and Build ID)
func (s *Streamer) Stream() error {
	build, err := s.Client.GetBuild(s.ctx, s.jobName, s.buildID)
	if err != nil {
		return fmt.Errorf("<?> Error: Failed to get build %d for job %s: %w", s.buildID, s.jobName, err)
	}

	s.printHeader(build)

	if s.follow {
		return s.streamWithFollow(build)
	}

	return s.streamOnce(build)
}

// streamOnce displays logs once (waits for completion if still building)
//
// Parameters:
//   - build: jenkins build
//
// Errors possible causes:
//   - can't poll build
//   - can't read logs
//
// Examples:
// err := s.streamOnce(build)
func (s *Streamer) streamOnce(build *gojenkins.Build) error {
	// Poll to get latest build info
	_, err := build.Poll(s.ctx)
	if err != nil {
		return fmt.Errorf("<?> Error: Failed to poll build %d for job %s: %w", s.buildID, s.jobName, err)
	}

	if build.IsRunning(s.ctx) {
		fmt.Println("  Waiting for build to complete...")

		for range constants.MaxPollAttempts {
			select {
			case <-s.ctx.Done():
				return nil
			case <-time.After(constants.PollInterval):
				_, err := build.Poll(s.ctx)
				if err != nil {
					return err
				}
				if !build.IsRunning(s.ctx) {
					return s.fetchAndDisplayLogs(build)
				}
			}
		}
	}

	return s.fetchAndDisplayLogs(build)
}

// streamWithFollow displays logs in real time
//
// Parameters:
//   - build: jenkins build
//
// Errors possible causes:
//   - can't poll build
//   - can't read logs
//
// Examples:
// err := s.streamWithFollow(build)
func (s *Streamer) streamWithFollow(build *gojenkins.Build) error {
	fmt.Println("  Following logs (press Ctrl+C to stop)...")

	ticker := time.NewTicker(constants.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			fmt.Println("\n\n  Log streaming stopped.")
			return nil

		case <-ticker.C:
			_, err := build.Poll(s.ctx)
			if err != nil {
				return err
			}

			logs := build.GetConsoleOutput(s.ctx)

			lines := strings.Split(logs, "\n")
			for _, line := range lines {
				if line == "" {
					continue
				}

				lineKey := fmt.Sprintf("%d:%s", s.buildID, line)
				if s.seenLogs[lineKey] {
					continue
				}
				s.seenLogs[lineKey] = true

				s.printLogLine(line)
			}

			if !build.IsRunning(s.ctx) {
				s.formatCompletion(build)
				return nil
			}
		}
	}
}

// fetchAndDisplayLogs fetches and displays formatted and colorized logs
//
// Parameters:
//   - build: jenkins build
//
// Errors possible causes:
//   - can't read logs
//
// Examples:
// err := s.fetchAndDisplayLogs(build)
func (s *Streamer) fetchAndDisplayLogs(build *gojenkins.Build) error {
	logs := build.GetConsoleOutput(s.ctx)

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

	return nil
}

// printLogLine prints colorized log line
//
// Parameters:
//   - line: log line
//
// Examples:
// s.printLogLine(line)
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

// colorizeContent colorizes content depending on log level
//
// Parameters:
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
// Parameters:
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
// Parameters:
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

// applyTail cuts the lines > s.tailLines
//
// Parameters:
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

// printHeader prints build header
//
// Parameters:
//   - build: jenkins build
//
// Examples:
// s.printHeader(build)
func (s *Streamer) printHeader(build *gojenkins.Build) {
	if s.colorize {
		_, err := color.New(color.Bold).Println("Jenkins Build")
		if err != nil {
			fmt.Printf("<?> Error: %v\n", err)
			return
		}
	} else {
		fmt.Println("Jenkins Build")
	}

	fmt.Println(strings.Repeat("─", 80))
	fmt.Printf("Job:        %s\n", s.jobName)
	fmt.Printf("Build #:    %d\n", build.GetBuildNumber())
	fmt.Printf("Status:     %s\n", s.formatBuildStatus(build))
	fmt.Printf("Result:     %s\n", build.GetResult())
	fmt.Printf("Duration:   %.5f ms\n", build.GetDuration())
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()
}

// formatBuildStatus formats status of build
//
// Parameters:
//   - build: jenkins build
//
// Examples:
// status := s.formatBuildStatus(build)
func (s *Streamer) formatBuildStatus(build *gojenkins.Build) string {
	if build.IsRunning(s.ctx) {
		return "Running"
	}
	return "Completed"
}

// formatCompletion formats completion process
//
// Parameters:
//   - build: jenkins build
//
// Examples:
// s.formatCompletion(build)
func (s *Streamer) formatCompletion(build *gojenkins.Build) {
	fmt.Println()
	fmt.Println(strings.Repeat("-", 80))
	result := build.GetResult()

	if s.colorize {
		switch result {
		case "SUCCESS":
			color.Green("Build completed successfully.")
		case "FAILURE":
			color.Red("Build failed.")
		case "ABORTED":
			color.Yellow("Build aborted.")
		case "UNSTABLE":
			color.Yellow("Build unstable.")
		default:
			color.White("Build completed: %s\n", result)
		}
	} else {
		fmt.Printf("Build completed: %s\n", result)
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
