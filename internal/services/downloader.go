package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownloadFile downloads a file from the specified URL and saves it to the specified filepath.
// If an error occurs during the process, it returns an error with a descriptive message.
func DownloadFile(url string, filepath string) error {
	// Send an HTTP GET request to the specified URL.
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to initiate request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 OK.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create a new file with the specified filepath.
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Copy the response body to the file.
	if _, err = io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to copy content to file: %w", err)
	}

	return nil
}
