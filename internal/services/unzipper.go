package services

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ExtractCSVFromZip extracts the first CSV file found inside a zip archive and saves it to the specified destination directory.
// Returns the path of the extracted CSV file or an error if one occurs.
func ExtractCSVFromZip(zipFilePath string, destDir string) (string, error) {
	// Open the zip file
	zipReader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open zip file: %w", err)
	}
	defer zipReader.Close()

	// Loop through the files in the zip archive
	for _, file := range zipReader.File {
		// Check if the file is a CSV file
		if filepath.Ext(file.Name) == ".csv" {
			// Open the CSV file inside the zip
			rc, err := file.Open()
			if err != nil {
				return "", fmt.Errorf("failed to open file in zip: %w", err)
			}
			defer rc.Close()

			// Define the path to extract the CSV file to
			extractedFilePath := filepath.Join(destDir, filepath.Base(file.Name))

			// Create the file to write the extracted CSV
			outFile, err := os.Create(extractedFilePath)
			if err != nil {
				return "", fmt.Errorf("failed to create extracted file: %w", err)
			}
			defer outFile.Close()

			// Copy the content from the zip to the new file
			if _, err = io.Copy(outFile, rc); err != nil {
				return "", fmt.Errorf("failed to copy file content: %w", err)
			}

			return extractedFilePath, nil
		}
	}

	return "", fmt.Errorf("no CSV file found in zip")
}
