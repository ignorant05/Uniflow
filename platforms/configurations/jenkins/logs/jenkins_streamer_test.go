package jenkins

import (
	"context"
	"strings"
	"testing"

	"github.com/bndr/gojenkins"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockJenkinsClient is a mock implementation of the Jenkins client
type MockJenkinsClient struct {
	mock.Mock
}

func (m *MockJenkinsClient) GetBuild(ctx context.Context, jobName string, buildID int64) (*gojenkins.Build, error) {
	args := m.Called(ctx, jobName, buildID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gojenkins.Build), args.Error(1)
}

// TestNewStreamer tests streamer creation
func TestNewStreamer(t *testing.T) {
	client := &gojenkins.Jenkins{}
	opts := StreamerOptions{
		Follow:    true,
		TailLines: 10,
		Colorize:  true,
		ShowTime:  true,
	}

	streamer := NewStreamer(client, "test-job", 123, opts)

	assert.NotNil(t, streamer)
	assert.Equal(t, "test-job", streamer.jobName)
	assert.Equal(t, int64(123), streamer.buildID)
	assert.Equal(t, true, streamer.follow)
	assert.Equal(t, 10, streamer.tailLines)
	assert.Equal(t, true, streamer.colorize)
	assert.Equal(t, true, streamer.showTime)
	assert.NotNil(t, streamer.seenLogs)
}

// TestDetectLogLevel tests log level detection
func TestDetectLogLevel(t *testing.T) {
	streamer := &Streamer{}

	tests := []struct {
		name     string
		content  string
		expected LogLevel
	}{
		{"Error detection", "ERROR: Something went wrong", LevelError},
		{"Error lowercase", "error in system", LevelError},
		{"Debug detection", "DEBUG: Debugging info", LevelDebug},
		{"Warning detection", "WARNING: Be careful", LevelWarning},
		{"Success detection", "SUCCESS: Build passed", LevelSuccess},
		{"Info default", "Just some info", LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := streamer.detectLogLevel(tt.content)
			assert.Equal(t, tt.expected, level)
		})
	}
}

// TestFormatTimestamp tests timestamp formatting
func TestFormatTimestamp(t *testing.T) {
	streamer := &Streamer{}

	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			"Valid RFC3339Nano",
			"2024-01-14T15:30:45.123456789Z",
			false,
		},
		{
			"Invalid timestamp",
			"invalid-time",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := streamer.formatTimestamp(tt.input)
			assert.NotNil(t, result)
			if !tt.shouldErr {
				assert.NotEmpty(t, result)
			}
		})
	}
}

// TestApplyTail tests log tail functionality
func TestApplyTail(t *testing.T) {
	tests := []struct {
		name      string
		logs      string
		tailLines int
		expected  string
	}{
		{
			name:      "Tail with fewer lines than total",
			logs:      "line1\nline2\nline3\nline4\nline5",
			tailLines: 3,
			expected:  "line3\nline4\nline5",
		},
		{
			name:      "Tail with more lines than total",
			logs:      "line1\nline2\nline3",
			tailLines: 10,
			expected:  "line1\nline2\nline3",
		},
		{
			name:      "Tail zero lines",
			logs:      "line1\nline2\nline3",
			tailLines: 0,
			expected:  "",
		},
		{
			name:      "Empty logs",
			logs:      "",
			tailLines: 5,
			expected:  "",
		},
		{
			name:      "Single line log",
			logs:      "single line",
			tailLines: 1,
			expected:  "single line",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			streamer := &Streamer{tailLines: tt.tailLines}
			result := streamer.applyTail(tt.logs)

			// Compare the actual result with exected
			assert.Equal(t, tt.expected, result,
				"Test '%s' failed: expected %q, got %q",
				tt.name, tt.expected, result)

			// Additional check: verify line count if needed
			if tt.tailLines > 0 && result != "" {
				resultLines := strings.Count(result, "\n")
				if !strings.HasSuffix(result, "\n") {
					resultLines++ // Add 1 if result doesn't end with newline
				}
				expectedLines := strings.Count(tt.expected, "\n")
				if !strings.HasSuffix(tt.expected, "\n") && tt.expected != "" {
					expectedLines++
				}
				assert.Equal(t, expectedLines, resultLines,
					"Line count mismatch for test '%s'", tt.name)
			}
		})

	}
}

// TestColorizeContent tests content colorization
func TestColorizeContent(t *testing.T) {
	streamer := &Streamer{}

	tests := []struct {
		name     string
		content  string
		level    LogLevel
		expected string
	}{
		{"Error colorize", "ERROR occurred", LevelError, "ERROR occurred"},
		{"Warning colorize", "WARNING message", LevelWarning, "WARNING message"},
		{"Success colorize", "SUCCESS", LevelSuccess, "SUCCESS"},
		{"Debug colorize", "DEBUG info", LevelDebug, "DEBUG info"},
		{"Info colorize", "info message", LevelInfo, "info message"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := streamer.colorizeContent(tt.content, tt.level)
			assert.NotEmpty(t, result)
		})
	}
}

// TestStop tests graceful stop
func TestStop(t *testing.T) {
	client := &gojenkins.Jenkins{}
	streamer := NewStreamer(client, "test-job", 123, StreamerOptions{})

	select {
	case <-streamer.ctx.Done():
		t.Fatal("Context should not be cancelled initially")
	default:
	}

	streamer.Stop()

	select {
	case <-streamer.ctx.Done():
	default:
		t.Fatal("Context should be cancelled after Stop")
	}
}

// TestPrintLogLine tests log line printing
func TestPrintLogLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		colorize bool
		showTime bool
	}{
		{"Simple line", "test log line", false, false},
		{"With color", "ERROR: test error", true, false},
		{"With timestamp", "[2024-01-14T15:30:45Z] test log", false, true},
		{"Full options", "DEBUG: test debug", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			streamer := &Streamer{
				colorize: tt.colorize,
				showTime: tt.showTime,
			}

			assert.NotPanics(t, func() {
				streamer.printLogLine(tt.line)
			})
		})
	}
}

// TestLogLevelConstants tests log level constants are properly defined
func TestLogLevelConstants(t *testing.T) {
	assert.Equal(t, LogLevel(0), LevelDebug)
	assert.Equal(t, LogLevel(1), LevelInfo)
	assert.Equal(t, LogLevel(2), LevelError)
	assert.Equal(t, LogLevel(3), LevelWarning)
	assert.Equal(t, LogLevel(4), LevelSuccess)
}
