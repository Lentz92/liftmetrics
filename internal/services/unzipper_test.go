package services

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

// TestExtractCSVFromZip tests the ExtractCSVFromZip function
func TestExtractCSVFromZip(t *testing.T) {
	// Create a temporary directory for the zip file and the extracted CSV
	tempDir := t.TempDir()
	zipFilePath := filepath.Join(tempDir, "test.zip")
	csvContent := "col1,col2\nval1,val2"

	// Create a zip file with a single CSV file inside it
	createTestZipFile(zipFilePath, "test.csv", csvContent)

	// Extract the CSV file from the zip archive
	extractedFilePath, err := ExtractCSVFromZip(zipFilePath, tempDir)
	if err != nil {
		t.Fatalf("ExtractCSVFromZip failed: %v", err)
	}

	// Check if the extracted file exists
	if _, err := os.Stat(extractedFilePath); os.IsNotExist(err) {
		t.Fatalf("Extracted file does not exist")
	}

	// Check the file content
	content, err := os.ReadFile(extractedFilePath)
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}
	if string(content) != csvContent {
		t.Fatalf("File content mismatch. Got %q, expected %q", string(content), csvContent)
	}
}

// createTestZipFile creates a zip file for testing purposes
func createTestZipFile(zipFilePath, csvFileName, csvContent string) {
	// Create the zip file
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		panic(err)
	}
	defer zipFile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Create a new file inside the zip
	writer, err := zipWriter.Create(csvFileName)
	if err != nil {
		panic(err)
	}

	// Write the content to the file inside the zip
	_, err = writer.Write([]byte(csvContent))
	if err != nil {
		panic(err)
	}
}
