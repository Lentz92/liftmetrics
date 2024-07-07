package services

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	"liftmetrics/internal/db"
	"liftmetrics/pkg"
)

// SetupDatabase handles the entire process of downloading, extracting, and setting up the database
func SetupDatabase(url, filePath, dataDir, dbFilePath string) error {
	// Download the ZIP file
	if err := DownloadFile(url, filePath); err != nil {
		return fmt.Errorf("error downloading file: %w", err)
	}

	// Extract the CSV file from the ZIP
	if _, err := ExtractCSVFromZip(filePath, dataDir); err != nil {
		return fmt.Errorf("failed to extract CSV: %w", err)
	}

	// Find the CSV file
	csvFilePath, err := pkg.FindCSVFile(dataDir)
	if err != nil {
		return fmt.Errorf("failed to find CSV file: %w", err)
	}

	// Open and parse the CSV file
	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer csvFile.Close()

	var records []*db.Record
	if err := gocsv.UnmarshalFile(csvFile, &records); err != nil {
		return fmt.Errorf("failed to parse CSV file: %w", err)
	}

	// Create the database
	fmt.Println("\nSetting up the database. This might take a while!")
	database, err := db.CreateDatabase(dbFilePath, true)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	defer database.Close()

	// Populate the database with records
	if err := db.PopulateDatabase(database, records); err != nil {
		return fmt.Errorf("failed to populate database: %w", err)
	}

	// Calculate Successful attempts for the database
	fmt.Println("\nCalculating new metrics for the database!")
	if err := db.CalculateSuccessfulAttempts(database); err != nil {
		return fmt.Errorf("failed to update lift attempts in database: %w", err)
	}

	// Calculate relative diff in attempts
	if err := db.CalculateRelativeDiff(database); err != nil {
		return fmt.Errorf("failed to calculate and insert relative differences: %w", err)
	}

	fmt.Println("\nAll setup and calculations complete!")
	return nil
}
