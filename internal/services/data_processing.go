package services

import (
	"context"
	"errors"
	"github.com/gocarina/gocsv"
	"liftmetrics/internal/db"
	"liftmetrics/pkg"
	"log"
	"os"
	"time"
)

// Error variables for specific failure scenarios
var (
	ErrDownloadFailed    = errors.New("failed to download file")
	ErrExtractFailed     = errors.New("failed to extract CSV")
	ErrCSVNotFound       = errors.New("failed to find CSV file")
	ErrCSVOpenFailed     = errors.New("failed to open CSV file")
	ErrCSVParseFailed    = errors.New("failed to parse CSV file")
	ErrDBCreateFailed    = errors.New("failed to create database")
	ErrDBPopulateFailed  = errors.New("failed to populate database")
	ErrCalculationFailed = errors.New("failed to perform database calculations")
)

// SetupDatabase handles the entire process of checking for updates, downloading,
// extracting, and setting up the database if necessary.
//
// Parameters:
//   - dataURL: The URL of the ZIP file to download.
//   - websiteURL: The URL of the website containing the current revision number.
//   - filePath: The path where the downloaded ZIP file should be saved.
//   - dataDir: The directory where the CSV should be extracted.
//   - dbFilePath: The path where the database file should be created.
//
// Returns:
//   - error: An error if any step of the process fails, nil otherwise.
func SetupDatabase(dataURL, websiteURL, filePath, dataDir, dbFilePath string) error {
	// Check if we need to update the database
	needsUpdate, err := checkForUpdate(dataDir, websiteURL)
	if err != nil {
		log.Printf("Failed to check for updates: %v", err)
		return err
	}

	// If no update is needed, we can exit early
	if !needsUpdate {
		log.Println("Data is up to date. Skipping download and database setup.")
		return nil
	}

	log.Println("Update needed. Proceeding with download and database setup.")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Download the ZIP file containing the new data
	if err := DownloadFile(ctx, dataURL, filePath); err != nil {
		log.Printf("Error downloading file: %v", err)
		return ErrDownloadFailed
	}

	// Extract the CSV file from the downloaded ZIP
	if _, err := ExtractCSVFromZip(filePath, dataDir); err != nil {
		log.Printf("Failed to extract CSV: %v", err)
		return ErrExtractFailed
	}

	// Locate the extracted CSV file
	csvFilePath, err := pkg.FindCSVFile(dataDir)
	if err != nil {
		log.Printf("Failed to find CSV file: %v", err)
		return ErrCSVNotFound
	}

	// Open the CSV file for reading
	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		log.Printf("Failed to open CSV file: %v", err)
		return ErrCSVOpenFailed
	}
	defer csvFile.Close() // Ensure the file is closed when we're done

	// Parse the CSV file into a slice of Record structs
	var records []*db.Record
	if err := gocsv.UnmarshalFile(csvFile, &records); err != nil {
		log.Printf("Failed to parse CSV file: %v", err)
		return ErrCSVParseFailed
	}

	// Create a new database or open an existing one
	log.Println("Setting up the database. This might take a while!")
	database, err := db.CreateDatabase(dbFilePath, true)
	if err != nil {
		log.Printf("Failed to create database: %v", err)
		return ErrDBCreateFailed
	}
	defer database.Close() // Ensure the database connection is closed when we're done

	// Populate the database with the parsed records
	if err := db.PopulateDatabase(database, records); err != nil {
		log.Printf("Failed to populate database: %v", err)
		return ErrDBPopulateFailed
	}

	// Calculate and update successful attempts for each lift type
	log.Println("Calculating new metrics for the database!")
	fc := db.NewFeatureCalculator()
	if err := fc.UpdateAllMetrics(ctx, database); err != nil {
		log.Printf("Failed to update metrics in database: %v", err)
		return ErrCalculationFailed
	}

	log.Println("All setup and calculations complete!")
	return nil
}

// checkForUpdate compares the local CSV revision with the website revision.
// It returns true if an update is needed, false otherwise.
func checkForUpdate(dataDir, websiteURL string) (bool, error) {
	// Attempt to find the local CSV file
	csvFilePath, err := pkg.FindCSVFile(dataDir)
	if err != nil {
		// If no local file is found, we assume an update is needed
		log.Printf("No local CSV file found: %v", err)
		return true, nil
	}

	// Use the CheckRevision function to compare local and web revisions
	needsUpdate, err := CheckRevision(csvFilePath, websiteURL)
	if err != nil {
		log.Printf("Failed to check revisions: %v", err)
		return false, err
	}

	return needsUpdate, nil
}
