package constants

import (
	"time"
)

// Streamer constants
const (
	MAX_REDIRECTS = 10

	// logs max size
	DATA_LOGS_MAX_SIZE = 10 * 1024 * 1024
)

const (
	PollInterval    = 3 * time.Second
	MaxPollAttempts = 600
	JobCheckIntervl = 5 * time.Second
)
