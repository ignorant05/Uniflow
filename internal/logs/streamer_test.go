package logs

import (
	"context"
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/ignorant05/Uniflow/configs/github"
	"github.com/ignorant05/Uniflow/internal/config"
)

// Test streamer creation
func TestNewStreamer(t *testing.T) {
	tests := []struct {
		name    string
		owner   string
		repo    string
		runID   int64
		opts    StreamerOptions
		needErr bool
	}{
		{
			name:  "valid streamer creation",
			owner: "ignorant05",
			repo:  "Uniflow",
			runID: 123456,
			opts: StreamerOptions{
				Follow:    true,
				TailLines: 50,
				Colorize:  true,
			},
			needErr: false,
		},
		{
			name:    "streamer without opts",
			owner:   "ignorant05",
			repo:    "Uniflow",
			runID:   123455,
			opts:    StreamerOptions{},
			needErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.GithubConfig{
				Token:             "random-gibbrich-as-token",
				DefaultRepository: tt.repo,
			}

			client, _ := github.NewClient(context.Background(), cfg)

			streamer := NewStreamer(client, tt.owner, tt.repo, tt.runID, tt.opts)

			if streamer == nil {
				t.Error("New Streamer returned nil")
				return
			}

			if streamer.owner != tt.owner {
				t.Errorf("owner = %v, want %v", streamer.owner, tt.owner)
			}
			if streamer.repo != tt.repo {
				t.Errorf("repo = %v, want %v", streamer.repo, tt.repo)
			}
			if streamer.runID != tt.runID {
				t.Errorf("runID = %v, want %v", streamer.runID, tt.runID)
			}
			if streamer.follow != tt.opts.Follow {
				t.Errorf("follow = %v, want %v", streamer.follow, tt.opts.Follow)
			}
			if streamer.tailLines != tt.opts.TailLines {
				t.Errorf("tailLines = %v, want %v", streamer.tailLines, tt.opts.TailLines)
			}
			if streamer.showTime != tt.opts.ShowTime {
				t.Errorf("showTime = %v, want %v", streamer.showTime, tt.opts.ShowTime)
			}
			if streamer.colorize != tt.opts.Colorize {
				t.Errorf("colorize = %v, want %v", streamer.colorize, tt.opts.Colorize)
			}

			if streamer.ctx == nil {
				t.Error("context is nil")
			}
			if streamer.cancelFunc == nil {
				t.Error("cancelFunc is nil")
			}

			if streamer.seenJobs == nil {
				t.Error("seenJobs map is nil")
			}
		})
	}
}

// Test streamer stop
func TestStreamerStop(t *testing.T) {
	cfg := &config.GithubConfig{
		Token: "random-gibbrich-as-token",
	}

	client, _ := github.NewClient(context.Background(), cfg)

	streamer := NewStreamer(client, "ignorant05", "Uniflow", 12345, StreamerOptions{})
	select {
	case <-streamer.ctx.Done():
		t.Error("Context should not be done initially")
	default:
		t.Context().Err()
	}

	streamer.Stop()

	select {
	case <-streamer.ctx.Done():
		t.Context().Err()

	case <-time.After(100 * time.Millisecond):
		t.Error("Context should be cancelled after using streamer.Stop()")
	}
}

func TestDetectLogLevel(t *testing.T) {
	streamer := &Streamer{}

	tests := []struct {
		name    string
		content string
		want    LogLevel
	}{
		{
			name:    "error level",
			content: "ERROR: something went wrong",
			want:    LevelError,
		},
		{
			name:    "error with failed keyword",
			content: "Build failed with exit code 1",
			want:    LevelError,
		},
		{
			name:    "fatal error",
			content: "FATAL: cannot continue",
			want:    LevelError,
		},
		{
			name:    "warning level",
			content: "WARNING: deprecated function",
			want:    LevelWarning,
		},
		{
			name:    "warn keyword",
			content: "Warn: unused variable",
			want:    LevelWarning,
		},
		{
			name:    "success level",
			content: "✓ Build completed successfully",
			want:    LevelSuccess,
		},
		{
			name:    "success keyword",
			content: "Tests passed",
			want:    LevelSuccess,
		},
		{
			name:    "debug level",
			content: "DEBUG: checking values",
			want:    LevelDebug,
		},
		{
			name:    "info level",
			content: "Starting deployment process",
			want:    LevelInfo,
		},
		{
			name:    "case insensitive error",
			content: "Error: file not found",
			want:    LevelError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := streamer.detectLogLevel(tt.content)
			if got != tt.want {
				t.Errorf("detectLogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test format timestamp
func TestFormatTimestamp(t *testing.T) {
	streamer := &Streamer{}

	tests := []struct {
		name      string
		timestamp string
		want      string
	}{
		{
			name:      "valid timestamp",
			timestamp: "2024-01-15T10:30:45.1234567Z",
			want:      "10:30:45",
		},
		{
			name:      "midnight",
			timestamp: "2024-01-15T00:00:00.0000000Z",
			want:      "00:00:00",
		},
		{
			name:      "noon",
			timestamp: "2024-01-15T12:00:00.0000000Z",
			want:      "12:00:00",
		},
		{
			name:      "invalid timestamp",
			timestamp: "invalid",
			want:      "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := streamer.formatTimestamp(tt.timestamp)
			if got != tt.want {
				t.Errorf("formatTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Testing Apply tail
func TestApplyTail(t *testing.T) {
	streamer := &Streamer{}

	tests := []struct {
		name      string
		logs      string
		tailLines int64
		wantLines int64
	}{
		{
			name:      "tail 3 lines from 5",
			logs:      "line1\nline2\nline3\nline4\nline5",
			tailLines: 3,
			wantLines: 3,
		},
		{
			name:      "tail more than available",
			logs:      "line1\nline2",
			tailLines: 10,
			wantLines: 2,
		},
		{
			name:      "empty logs",
			logs:      "",
			tailLines: 5,
			wantLines: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			streamer.tailLines = int(tt.tailLines)
			got := streamer.applyTail(tt.logs)

			lines := int64(len(splitLines(got)))
			if lines != tt.wantLines {
				t.Errorf("TailLines() = %v, want %v", lines, tt.wantLines)
			}
		})
	}
}

// Test format status
func TestFormatStatus(t *testing.T) {
	streamer := Streamer{}

	tests := []struct {
		name   string
		status string
		want   string
	}{
		{
			name:   "queued status",
			status: "queued",
			want:   "Queued",
		},
		{
			name:   "in_progress status",
			status: "in_progress",
			want:   "In Progress",
		},
		{
			name:   "completed status",
			status: "completed",
			want:   "Completed",
		},
		{
			name:   "waiting status",
			status: "waiting",
			want:   "Waiting",
		},
		{
			name:   "unknown status",
			status: "unknown",
			want:   "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := streamer.formatStatus(tt.status)
			if got != tt.want {
				t.Errorf("formatStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test format conclusion
func TestFormatConclusion(t *testing.T) {
	streamer := Streamer{}

	tests := []struct {
		name       string
		conclusion string
		want       string
	}{
		{
			name:       "success",
			conclusion: "success",
			want:       "Success",
		},
		{
			name:       "failure",
			conclusion: "failure",
			want:       "Failure",
		},
		{
			name:       "cancelled",
			conclusion: "cancelled",
			want:       "Cancelled",
		},
		{
			name:       "skipped",
			conclusion: "skipped",
			want:       "Skipped",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := streamer.formatConclusion(tt.conclusion)
			if got != tt.want {
				t.Errorf("formatConclusion() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test colorize content
func TestColorizeContent(t *testing.T) {
	// Enable colors for this test
	color.NoColor = false
	defer func() { color.NoColor = true }()

	streamer := Streamer{}
	tests := []struct {
		name    string
		content string
		level   LogLevel
	}{
		{
			name:    "error level",
			content: "Error occurred",
			level:   LevelError,
		},
		{
			name:    "warning level",
			content: "Warning message",
			level:   LevelWarning,
		},
		{
			name:    "success level",
			content: "Success message",
			level:   LevelSuccess,
		},
		{
			name:    "debug level",
			content: "Debug message",
			level:   LevelDebug,
		},
		{
			name:    "info level",
			content: "Info message",
			level:   LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := streamer.colorizeContent(tt.content, tt.level)

			if got == "" {
				t.Error("colorizeContent() returned empty string")
			}

			if tt.level != LevelInfo {
				if got == tt.content {
					t.Error("colorizeContent() should add color codes")
				}
			}
		})
	}
}

// Test colorize color: no color
func TestColorizeContentNoColor(t *testing.T) {
	streamer := &Streamer{colorize: false}

	content := "Test message"
	level := LevelError

	got := streamer.colorizeContent(content, level)

	if got != content {
		t.Errorf("colorizeContent() with colorize=false = %v, want %v", got, content)
	}
}

// This is very much needed to repsplit the logs
func splitLines(s string) []string {
	if s == "" {
		return []string{""}
	}
	lines := []string{}
	current := ""
	for _, c := range s {
		if c == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

// Benchmark testing detect log level
func BenchmarkDetectLogLevel(b *testing.B) {
	streamer := Streamer{}

	content := "Error: Something went wrong"

	b.ResetTimer()
	for b.Loop() {
		streamer.detectLogLevel(content)
	}
}

// Benchmark testing format timestamp
func BenchmarkFormatTimestamp(b *testing.B) {
	streamer := Streamer{}

	timestamp := "2024-01-15T10:30:45.1234567Z"

	b.ResetTimer()
	for b.Loop() {
		streamer.formatTimestamp(timestamp)
	}
}

// Benchmark testing colorize content
func BenchmarkColorizeContent(b *testing.B) {
	streamer := Streamer{}
	context := "Error: Something went wrong"
	level := LevelError

	b.ResetTimer()
	for b.Loop() {
		streamer.colorizeContent(context, level)
	}
}
