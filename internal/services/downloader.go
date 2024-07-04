package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(url string, filepath string) error {
	// Send an HTTP GET request to the specified URL
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to initiate request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create a new file with the specified filepath
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Copy the response body to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy content to file: %w", err)
	}

	return nil
}
