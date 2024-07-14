package services

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// DownloadFile downloads a file from the specified URL and saves it to the specified filepath.
// It removes any existing CSV files in the directory before downloading.
// It uses a context for cancellation and timeout support, and automatically switches to
// streaming for very large files.
//
// Parameters:
//   - ctx: A context.Context for cancellation and timeout control.
//   - url: The URL of the file to download.
//   - filepath: The local path where the downloaded file will be saved.
//
// Returns:
//   - error: An error if any step of the download process fails, nil otherwise.
//
// The function will use in-memory operations for files smaller than 2 GB,
// and will switch to streaming for larger files to conserve memory.
// This threshold is set considering modern hardware capabilities, where most systems
// have at least 8 GB of RAM, with 16 GB or more being common.
func DownloadFile(ctx context.Context, url string, filepath string) error {
	log.Printf("Starting download from %s to %s", url, filepath)

	// Remove existing CSV files
	if err := removeExistingCSVFiles(filepath); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return errors.Join(errors.New("failed to create request"), err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Join(errors.New("failed to send request"), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-OK status: %s", resp.Status)
		return errors.New("bad status: " + resp.Status)
	}

	// 2 GB threshold: 2 * 1024 * 1024 * 1024 bytes
	if resp.ContentLength > 2*1024*1024*1024 {
		log.Printf("File size (%d bytes) exceeds threshold. Using streaming download.", resp.ContentLength)
		return downloadLargeFile(resp.Body, filepath)
	}

	log.Printf("File size: %d bytes. Downloading into memory.", resp.ContentLength)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Join(errors.New("failed to read response body"), err)
	}

	err = os.WriteFile(filepath, body, 0644)
	if err != nil {
		return errors.Join(errors.New("failed to write file"), err)
	}

	log.Printf("Download completed successfully")
	return nil
}

// downloadLargeFile is a helper function that streams very large files (>2 GB) to disk.
// This function is called by DownloadFile when the file size exceeds 2 GB.
//
// Parameters:
//   - body: An io.Reader containing the file contents to be written.
//   - filepath: The local path where the file will be saved.
//
// Returns:
//   - error: An error if the file creation or copying process fails, nil otherwise.
func downloadLargeFile(body io.Reader, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return errors.Join(errors.New("failed to create file"), err)
	}
	defer out.Close()

	written, err := io.Copy(out, body)
	if err != nil {
		return errors.Join(errors.New("failed to copy content to file"), err)
	}

	log.Printf("Large file download completed. Bytes written: %d", written)
	return nil
}

// removeExistingCSVFiles removes any existing CSV files in the same directory as the filepath.
func removeExistingCSVFiles(filePath string) error {
	dir := filepath.Dir(filePath)
	files, err := os.ReadDir(dir)
	if err != nil {
		return errors.Join(errors.New("failed to read directory"), err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".csv") {
			fullPath := filepath.Join(dir, file.Name())
			if err := os.Remove(fullPath); err != nil {
				return errors.Join(errors.New("failed to remove existing CSV file"), err)
			}
			log.Printf("Removed existing CSV file: %s", fullPath)
		}
	}

	return nil
}
