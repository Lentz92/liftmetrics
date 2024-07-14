// Package services provides utility functions for file operations and data processing.
package services

import (
	"archive/zip"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ExtractCSVFromZip extracts the first CSV file found inside a zip archive and saves it to the specified destination directory.
// It uses logging for error reporting and ensures the destination directory exists before extraction.
//
// Parameters:
//   - zipFilePath: The path to the zip file containing the CSV.
//   - destDir: The directory where the CSV file should be extracted.
//
// Returns:
//   - string: The path of the extracted CSV file.
//   - error: An error if any step of the extraction process fails, nil otherwise.
//
// The function will create the destination directory if it doesn't exist.
// It searches for the first file with a .csv extension (case-insensitive) in the zip archive.
// If no CSV file is found, it returns an error.
func ExtractCSVFromZip(zipFilePath string, destDir string) (string, error) {
	// Open the zip file
	zipReader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		log.Printf("Failed to open zip file %s: %v", zipFilePath, err)
		return "", errors.New("failed to open zip file")
	}
	defer zipReader.Close()

	// Ensure the destination directory exists
	if err := os.MkdirAll(destDir, 0755); err != nil {
		log.Printf("Failed to create destination directory %s: %v", destDir, err)
		return "", errors.New("failed to create destination directory")
	}

	// Loop through the files in the zip archive
	for _, file := range zipReader.File {
		// Check if the file is a CSV file (case-insensitive)
		if strings.EqualFold(filepath.Ext(file.Name), ".csv") {
			return extractCSVFile(file, destDir)
		}
	}

	log.Printf("No CSV file found in zip %s", zipFilePath)
	return "", errors.New("no CSV file found in zip")
}

// extractCSVFile extracts a single CSV file from the zip archive to the destination directory.
// It is a helper function for ExtractCSVFromZip.
//
// Parameters:
//   - file: A pointer to the zip.File to be extracted.
//   - destDir: The directory where the file should be extracted.
//
// Returns:
//   - string: The path of the extracted CSV file.
//   - error: An error if any step of the extraction process fails, nil otherwise.
func extractCSVFile(file *zip.File, destDir string) (string, error) {
	// Open the CSV file inside the zip
	rc, err := file.Open()
	if err != nil {
		log.Printf("Failed to open file %s in zip: %v", file.Name, err)
		return "", errors.New("failed to open file in zip")
	}
	defer rc.Close()

	// Define the path to extract the CSV file to
	extractedFilePath := filepath.Join(destDir, filepath.Base(file.Name))

	// Create the file to write the extracted CSV
	outFile, err := os.Create(extractedFilePath)
	if err != nil {
		log.Printf("Failed to create extracted file %s: %v", extractedFilePath, err)
		return "", errors.New("failed to create extracted file")
	}
	defer outFile.Close()

	// Copy the content from the zip to the new file
	if _, err = io.Copy(outFile, rc); err != nil {
		log.Printf("Failed to copy content to file %s: %v", extractedFilePath, err)
		return "", errors.New("failed to copy file content")
	}

	return extractedFilePath, nil
}
