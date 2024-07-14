// Package services_test provides test cases for the services package.
// It focuses on testing the DownloadFile function and its helper functions.
package services

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

// TestDownloadFile is a test suite for the DownloadFile function.
// It covers various scenarios including successful downloads, error handling,
// and edge cases to ensure the function behaves correctly under different conditions.
func TestDownloadFile(t *testing.T) {
	// createTestServer is a helper function that creates a test HTTP server
	// with the given handler function and returns its URL.
	createTestServer := func(handler func(http.ResponseWriter, *http.Request)) string {
		server := httptest.NewServer(http.HandlerFunc(handler))
		return server.URL
	}

	// Test case for successful download of a small file
	t.Run("Successful download of small file", func(t *testing.T) {
		content := "Hello, World!"
		url := createTestServer(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, content)
		})

		// Create a temporary file for the download
		tempFile, err := os.CreateTemp("", "test-download-*")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())
		tempFile.Close()

		// Attempt to download the file
		err = DownloadFile(context.Background(), url, tempFile.Name())
		if err != nil {
			t.Fatalf("DownloadFile failed: %v", err)
		}

		// Verify the downloaded content
		downloadedContent, err := os.ReadFile(tempFile.Name())
		if err != nil {
			t.Fatalf("Failed to read downloaded file: %v", err)
		}

		if string(downloadedContent) != content {
			t.Errorf("Downloaded content does not match. Got %s, want %s", string(downloadedContent), content)
		}
	})

	// Test case for successful download of a large file
	t.Run("Successful download of large file", func(t *testing.T) {
		largeContent := make([]byte, 3*1024*1024*1024) // 3 GB of data
		url := createTestServer(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", fmt.Sprint(len(largeContent)))
			w.Write(largeContent)
		})

		// Create a temporary file for the large download
		tempFile, err := os.CreateTemp("", "test-large-download-*")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())
		tempFile.Close()

		// Attempt to download the large file
		err = DownloadFile(context.Background(), url, tempFile.Name())
		if err != nil {
			t.Fatalf("DownloadFile failed: %v", err)
		}

		// Verify the size of the downloaded file
		info, err := os.Stat(tempFile.Name())
		if err != nil {
			t.Fatalf("Failed to stat downloaded file: %v", err)
		}

		if info.Size() != int64(len(largeContent)) {
			t.Errorf("Downloaded file size does not match. Got %d, want %d", info.Size(), len(largeContent))
		}
	})

	// Test case for handling non-200 HTTP status codes
	t.Run("Non-200 status code", func(t *testing.T) {
		url := createTestServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})

		err := DownloadFile(context.Background(), url, "nonexistent.txt")
		if err == nil {
			t.Fatal("Expected an error, but got nil")
		}
		if err.Error() != "bad status: 404 Not Found" {
			t.Errorf("Unexpected error message. Got %s", err.Error())
		}
	})

	// Test case for handling context cancellation
	t.Run("Context cancellation", func(t *testing.T) {
		url := createTestServer(func(w http.ResponseWriter, r *http.Request) {
			// Simulate a slow response
			<-r.Context().Done()
		})

		ctx, cancel := context.WithCancel(context.Background())
		errChan := make(chan error)

		go func() {
			errChan <- DownloadFile(ctx, url, "cancelled.txt")
		}()

		// Cancel the context immediately
		cancel()

		select {
		case err := <-errChan:
			if err == nil {
				t.Fatal("Expected an error due to cancellation, but got nil")
			}
		case <-time.After(5 * time.Second):
			t.Fatal("Test timed out")
		}
	})

	// Test case for handling invalid URLs
	t.Run("Invalid URL", func(t *testing.T) {
		err := DownloadFile(context.Background(), "http://invalid-url", "invalid.txt")
		if err == nil {
			t.Fatal("Expected an error for invalid URL, but got nil")
		}
	})
}

// TestDownloadLargeFile tests the downloadLargeFile function directly.
// This function is a helper for DownloadFile, used for streaming large files.
func TestDownloadLargeFile(t *testing.T) {
	content := "Large file content"

	// Create a temporary file for the test
	tempFile, err := os.CreateTemp("", "test-large-file-*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Test the downloadLargeFile function
	err = downloadLargeFile(strings.NewReader(content), tempFile.Name())
	if err != nil {
		t.Fatalf("downloadLargeFile failed: %v", err)
	}

	// Verify the downloaded content
	downloadedContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}

	if string(downloadedContent) != content {
		t.Errorf("Downloaded content does not match. Got %s, want %s", string(downloadedContent), content)
	}
}
