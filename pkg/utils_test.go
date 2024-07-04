package pkg

import (
	"os"
	"path/filepath"
	"testing"
)

// TestFindCSVFile tests the FindCSVFile function
func TestFindCSVFile(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Create a set of test files in the directory
	files := []struct {
		name    string
		content string
	}{
		{"test1.txt", "This is a test file."},
		{"openipf-2023-07-04-rev1.csv", "col1,col2\nval1,val2"},
		{"openipf-2023-07-05-rev2.csv", "col1,col2\nval3,val4"},
		{"test2.csv", "col1,col2\nval5,val6"},
	}

	for _, file := range files {
		err := os.WriteFile(filepath.Join(tempDir, file.name), []byte(file.content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file.name, err)
		}
	}

	// Call the FindCSVFile function
	filePath, err := FindCSVFile(tempDir)
	if err != nil {
		t.Fatalf("FindCSVFile failed: %v", err)
	}

	// Expected file path (any of the openipf-*.csv files would be valid)
	expectedFilePaths := []string{
		filepath.Join(tempDir, "openipf-2023-07-04-rev1.csv"),
		filepath.Join(tempDir, "openipf-2023-07-05-rev2.csv"),
	}

	// Check if the returned file path matches any of the expected file paths
	found := false
	for _, expectedFilePath := range expectedFilePaths {
		if filePath == expectedFilePath {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("File path mismatch. Got %q, expected one of %q", filePath, expectedFilePaths)
	}
}
