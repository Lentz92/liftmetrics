package services

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ExtractCSVFromZip extract the first CSV file found inside a zip archive
func ExtractCSVFromZip(zipFilePath string, destDir string) (string, error) {
	// Open the zip file
	zipReader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open zip file: %w", err)
	}
	defer zipReader.Close()

	// Loop through the files in the zip archive
	for _, file := range zipReader.File {
		// check if the file is a CSV file
		if filepath.Ext(file.Name) == ".csv" {
			// Open the CSV file inside the zip
			rc, err := file.Open()
			if err != nil {
				return "", fmt.Errorf("failed to open file in zip: %w", err)
			}

			// Define the path to extract the CSV file to
			extractedFilePath := filepath.Join(destDir, filepath.Base(file.Name))

			// Create the file to write the extracted CSV
			outFile, err := os.Create(extractedFilePath)
			if err != nil {
				rc.Close()
				return "", fmt.Errorf("failed to create extracted file: %w", err)
			}

			// Copy the content from the zip to the new file
			_, err = io.Copy(outFile, rc)
			if err != nil {
				rc.Close()
				outFile.Close()
				return "", fmt.Errorf("failed to copy file content: %w", err)
			}

			// Close the files after copying is done
			rc.Close()
			outFile.Close()

			return extractedFilePath, nil
		}
	}

	return "", fmt.Errorf("no CSV file found in zip")
}
