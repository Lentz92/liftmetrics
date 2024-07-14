// Package services_test provides test cases for the services package.
// It focuses on testing the ExtractCSVFromZip function and its behavior under various conditions.
package services

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

// TestExtractCSVFromZip is a test suite for the ExtractCSVFromZip function.
// It covers various scenarios including successful extraction, error handling,
// and edge cases to ensure the function behaves correctly under different conditions.
func TestExtractCSVFromZip(t *testing.T) {
	// createTestZip is a helper function that creates a test zip file with the given contents.
	// It returns the path to the created zip file.
	//
	// Parameters:
	//   - files: A map where keys are filenames and values are file contents.
	//
	// Returns:
	//   - string: The path to the created zip file.
	createTestZip := func(files map[string]string) string {
		// Create a temporary file for the zip
		tmpfile, err := os.CreateTemp("", "test*.zip")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer tmpfile.Close()

		// Create a new zip writer
		zipWriter := zip.NewWriter(tmpfile)
		defer zipWriter.Close()

		// Add each file to the zip
		for name, content := range files {
			f, err := zipWriter.Create(name)
			if err != nil {
				t.Fatalf("Failed to create file in zip: %v", err)
			}
			_, err = f.Write([]byte(content))
			if err != nil {
				t.Fatalf("Failed to write content to zip file: %v", err)
			}
		}

		return tmpfile.Name()
	}

	// Test case: Successful extraction of CSV
	t.Run("Successful extraction of CSV", func(t *testing.T) {
		// Create a test zip file with a CSV and another file
		zipPath := createTestZip(map[string]string{
			"test.csv":  "col1,col2\nvalue1,value2",
			"other.txt": "Some other content",
		})
		defer os.Remove(zipPath)

		// Create a temporary directory for extraction
		destDir, err := os.MkdirTemp("", "extract")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(destDir)

		// Attempt to extract the CSV
		extractedPath, err := ExtractCSVFromZip(zipPath, destDir)
		if err != nil {
			t.Fatalf("ExtractCSVFromZip failed: %v", err)
		}

		// Read the contents of the extracted file
		content, err := os.ReadFile(extractedPath)
		if err != nil {
			t.Fatalf("Failed to read extracted file: %v", err)
		}

		// Verify the contents of the extracted file
		expectedContent := "col1,col2\nvalue1,value2"
		if string(content) != expectedContent {
			t.Errorf("Extracted content does not match. Got %s, want %s", string(content), expectedContent)
		}
	})

	// Test case: No CSV in zip
	t.Run("No CSV in zip", func(t *testing.T) {
		// Create a test zip file without a CSV
		zipPath := createTestZip(map[string]string{
			"test.txt": "Some content",
		})
		defer os.Remove(zipPath)

		// Create a temporary directory for extraction
		destDir, err := os.MkdirTemp("", "extract")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(destDir)

		// Attempt to extract a CSV (which doesn't exist)
		_, err = ExtractCSVFromZip(zipPath, destDir)
		if err == nil {
			t.Fatal("Expected an error, but got nil")
		}
		if err.Error() != "no CSV file found in zip" {
			t.Errorf("Unexpected error message. Got %s, want 'no CSV file found in zip'", err.Error())
		}
	})

	// Test case: Invalid zip file
	t.Run("Invalid zip file", func(t *testing.T) {
		// Create a file that's not a valid zip
		tmpfile, err := os.CreateTemp("", "invalid*.zip")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpfile.Name())
		tmpfile.Write([]byte("This is not a zip file"))
		tmpfile.Close()

		// Create a temporary directory for extraction
		destDir, err := os.MkdirTemp("", "extract")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(destDir)

		// Attempt to extract from the invalid zip file
		_, err = ExtractCSVFromZip(tmpfile.Name(), destDir)
		if err == nil {
			t.Fatal("Expected an error, but got nil")
		}
		if err.Error() != "failed to open zip file" {
			t.Errorf("Unexpected error message. Got %s, want 'failed to open zip file'", err.Error())
		}
	})

	// Test case: Case insensitive CSV extension
	t.Run("Case insensitive CSV extension", func(t *testing.T) {
		// Create a test zip file with a CSV that has uppercase extension
		zipPath := createTestZip(map[string]string{
			"TEST.CSV": "col1,col2\nvalue1,value2",
		})
		defer os.Remove(zipPath)

		// Create a temporary directory for extraction
		destDir, err := os.MkdirTemp("", "extract")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(destDir)

		// Attempt to extract the CSV
		extractedPath, err := ExtractCSVFromZip(zipPath, destDir)
		if err != nil {
			t.Fatalf("ExtractCSVFromZip failed: %v", err)
		}

		// Verify that the extracted file name matches the original (case-sensitive)
		if filepath.Base(extractedPath) != "TEST.CSV" {
			t.Errorf("Extracted file name does not match. Got %s, want TEST.CSV", filepath.Base(extractedPath))
		}
	})
}
