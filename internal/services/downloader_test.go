package services

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// TestDownloadFile tests the DownloadFile function
func TestDownloadFile(t *testing.T) {
	// Create a test server that serves a sample file
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("This is a test file."))
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	// Temporary directory for the downloaded file
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "testfile.txt")

	// Download the file from the test server
	err := DownloadFile(ts.URL, filePath)
	if err != nil {
		t.Fatalf("DownloadFile failed: %v", err)
	}

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("Downloaded file does not exist")
	}

	// Check the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}
	expectedContent := "This is a test file."
	if string(content) != expectedContent {
		t.Fatalf("File content mismatch. Got %q, expected %q", string(content), expectedContent)
	}
}
