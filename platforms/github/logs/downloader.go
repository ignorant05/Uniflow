package github

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ignorant05/Uniflow/platforms/github/constants"
)

type HttpClient interface {
	Get(url string) (*http.Response, error)
}

// Default client with default field timeout
var DefaultClient = &http.Client{
	Timeout: 30 * time.Second,
}

// readLogs reads and downloads logs
//
// Parameters :
//   - logsUrl: logs url
//
// Errors possible causes:
//   - invalid url
//   - internal error
//
// Example:
// body, err := s.readLogs(logsUrl)
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

// DownloadLogs reads and downloads logs
//
// Parameters :
//   - logsUrl: logs url
//
// Errors possible causes:
//   - invalid url
//   - internal error
//
// Example:
// body, err := DownloadLogs(logsUrl)
func DownloadLogs(logsUrl, downloadFileName string) error {
	if logsUrl == "" {
		return fmt.Errorf("<?> Error: Invalid URL.\n")
	}

	var path string
	if downloadFileName == "" {
		downloadFileName = constants.DEFAULT_DOWNLOAD_FILE_NAME
		path = constants.DEFAULT_DOWNLOAD_DIR_PATH + "/" + downloadFileName
	}

	if strings.HasPrefix(downloadFileName, "~/") || strings.HasPrefix(downloadFileName, "/home/") {
		path = downloadFileName
	}

	if !strings.HasSuffix(path, ".zip") {
		path += ".zip"
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

	file, err := os.Create(downloadFileName)
	if err != nil {
		return fmt.Errorf("<?> Error: Create output file: %w", err)
	}
	defer file.Close()

	limitedReader := io.LimitReader(resp.Body, constants.DATA_LOGS_MAX_SIZE)

	bytesWritten, err := io.Copy(file, limitedReader)
	if err != nil {
		return fmt.Errorf("<?> Error: Failed to write logs data: %w", err)
	}

	// Verify that the zip file is valid
	file.Seek(0, 0)
	_, err = zip.NewReader(file, bytesWritten)
	if err != nil {
		return fmt.Errorf("<?> Error: Failed to parse logs data.\n<?> Error: %w\n", err)
	}

	fmt.Printf("✓ Downloaded %d KB of logs to %s\n\n", bytesWritten/1024, downloadFileName)
	return nil
}
