package gitlab

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/ignorant05/Uniflow/internal/config"
	"github.com/ignorant05/Uniflow/platforms/configurations/gitlab"
)

// Test streamer creation
func TestNewStreamer(t *testing.T) {
	tests := []struct {
		name       string
		projectID  string
		pipelineID int
		opts       StreamerOptions
		needErr    bool
	}{
		{
			name:       "valid streamer creation",
			projectID:  "ignorant05/Uniflow",
			pipelineID: 123456,
			opts: StreamerOptions{
				Follow:    true,
				TailLines: 50,
				Colorize:  true,
			},
			needErr: false,
		},
		{
			name:       "streamer without opts",
			projectID:  "ignorant05/Uniflow",
			pipelineID: 123455,
			opts:       StreamerOptions{},
			needErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.GitlabConfig{
				Token:             "random-gibbrich-as-token",
				DefaultRepository: tt.projectID,
			}

			client, _ := gitlab.NewClient(context.Background(), cfg)

			streamer := NewStreamer(client.Client, tt.projectID, tt.pipelineID, tt.opts)

			if streamer == nil {
				t.Error("New Streamer returned nil")
				return
			}

			if streamer.projectID != tt.projectID {
				t.Errorf("projectID = %v, want %v", streamer.projectID, tt.projectID)
			}
			if streamer.pipelineID != tt.pipelineID {
				t.Errorf("pipelineID = %v, want %v", streamer.pipelineID, tt.pipelineID)
			}
			if streamer.follow != tt.opts.Follow {
				t.Errorf("follow = %v, want %v", streamer.follow, tt.opts.Follow)
			}
			if streamer.tailLines != tt.opts.TailLines {
				t.Errorf("tailLines = %v, want %v", streamer.tailLines, tt.opts.TailLines)
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
	cfg := &config.GitlabConfig{
		Token: "random-gibbrich-as-token",
	}

	client, _ := gitlab.NewClient(context.Background(), cfg)

	streamer := NewStreamer(client.Client, "ignorant05/Uniflow", 12345, StreamerOptions{})
	select {
	case <-streamer.ctx.Done():
		t.Error("Context should not be done initially")
	default:
		err := t.Context().Err()
		if err != nil {
			fmt.Printf("<?> Error: %v\n", err)
		}
	}

	streamer.Stop()

	select {
	case <-streamer.ctx.Done():
		err := t.Context().Err()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

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
		tailLines int
		wantLines int
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
			streamer.tailLines = tt.tailLines
			got := streamer.applyTail(tt.logs)

			lines := len(splitLines(got))
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
			name:   "created status",
			status: "created",
			want:   "Created",
		},
		{
			name:   "pending status",
			status: "pending",
			want:   "Pending",
		},
		{
			name:   "running status",
			status: "running",
			want:   "Running",
		},
		{
			name:   "success status",
			status: "success",
			want:   "Success",
		},
		{
			name:   "failed status",
			status: "failed",
			want:   "Failed",
		},
		{
			name:   "canceled status",
			status: "canceled",
			want:   "Canceled",
		},
		{
			name:   "skipped status",
			status: "skipped",
			want:   "Skipped",
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

// Test colorize content
func TestColorizeContent(t *testing.T) {
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

// Helper function to split lines
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
