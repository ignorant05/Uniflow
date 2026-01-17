package jenkins

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ignorant05/Uniflow/platforms/configurations/github/constants"
)

type HttpClient interface {
	Get(url string) (*http.Response, error)
}

// Default client with default field timeout
var DefaultClient = &http.Client{
	Timeout: 30 * time.Second,
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
		return fmt.Errorf("invalid URL, url: %s", logsUrl)
	}

	path := constants.DEFAULT_DOWNLOAD_DIR_PATH + "/" + constants.DEFAULT_DOWNLOAD_FILE_NAME

	if downloadFileName != "" {
		if strings.HasPrefix(downloadFileName, "~/") ||
			strings.HasPrefix(downloadFileName, "/home/") {
			path = downloadFileName
		} else {
			path = constants.DEFAULT_DOWNLOAD_DIR_PATH + "/" + downloadFileName
		}
	}

	if !strings.HasSuffix(path, ".zip") {
		path += ".zip"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", logsUrl, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "Uniflow-CLI")

	resp, err := DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to Download: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download logs data.\n<?> Error: Status Code: %d", resp.StatusCode)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("failed to close file: %v", err)
		}
	}()

	limitedReader := io.LimitReader(resp.Body, constants.DATA_LOGS_MAX_SIZE)

	bytesWritten, err := io.Copy(file, limitedReader)
	if err != nil {
		return fmt.Errorf("failed to write logs data: %w", err)
	}

	// Verify that the zip file is valid
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = zip.NewReader(file, bytesWritten)
	if err != nil {
		return fmt.Errorf("failed to parse logs data.\n<?> Error: %w", err)
	}

	fmt.Printf("✓ Downloaded %d KB of logs to %s\n\n", bytesWritten/1024, downloadFileName)
	return nil
}
