package logs

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	constants "github.com/ignorant05/Uniflow/internal/constants/logs"
)

type HttpClient interface {
	Get(url string) (*http.Response, error)
}

var DefaultClient = &http.Client{
	Timeout: 30 * time.Second,
}

func (s *Streamer) readLogs(logsUrl string) (string, error) {
	if logsUrl == "" {
		return "", fmt.Errorf("<?> Error: Invalid URL.\n")
	}

	resp, err := http.Get(logsUrl)
	if err != nil {
		return "", fmt.Errorf("<?> Error: Failed to download logs from urls: %s\n<?> Error: %w\n", logsUrl, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("<?> Error: Failed to download logs data.\n<?> Error: Status Code: %d\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("<?> Error: Failed to read logs.\n<?> Error: %w\b", err)
	}

	return string(body), nil
}

func DownloadLogs(logsUrl string) error {
	if logsUrl == "" {
		return fmt.Errorf("<?> Error: Invalid URL.\n")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", logsUrl, nil)
	if err != nil {
		return fmt.Errorf("<?> Error: Create request: %w", err)
	}
	req.Header.Set("User-Agent", "Uniflow-CLI")

	resp, err := DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("<?> Error: Download: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("<?> Error: Failed to download logs data.\n<?> Error: Status Code: %d\n", resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, constants.DATA_LOGS_MAX_SIZE))
	if err != nil {
		return fmt.Errorf("<?> Error: Failed to read logs data.\n<?> Error: %w\n", err)
	}

	_, err = zip.NewReader(bytes.NewReader(data), int64(len(data)/1024))
	if err != nil {
		return fmt.Errorf("<?> Error: Failed to parse logs data.\n<?> Error: %w\n", err)
	}

	fmt.Printf("<✓> Downloaded %d KB of logs\n\n", len(data)/1024)

	return nil
}
