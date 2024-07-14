package main

import (
	"log"
	"os"
	"path/filepath"

	"liftmetrics/internal/services"
)

const (
	dataURL    = "https://openpowerlifting.gitlab.io/opl-csv/files/openipf-latest.zip"
	websiteURL = "https://openpowerlifting.gitlab.io/opl-csv/bulk-csv.html"
	dataDir    = "../../data"
	zipFile    = "openipf-latest.zip"
	dbName     = "openipf.db"
)

func main() {
	// Configure logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Set up file paths
	absDataDir, err := filepath.Abs(dataDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path for data directory: %v", err)
	}

	filePath := filepath.Join(absDataDir, zipFile)
	dbDir := filepath.Join(absDataDir, "db")
	dbFilePath := filepath.Join(dbDir, dbName)

	// Create data and db directories
	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create directories: %v", err)
	}

	// Log the start of the process
	log.Println("Starting database setup process...")

	// Check for updates, download, extract, and process data if needed
	if err := setupDatabase(dataURL, websiteURL, filePath, absDataDir, dbFilePath); err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}

	log.Println("Database setup process completed.")
}

// setupDatabase wraps the services.SetupDatabase function with additional logging.
func setupDatabase(dataURL, websiteURL, filePath, absDataDir, dbFilePath string) error {
	log.Println("Checking for updates and processing data if needed...")
	err := services.SetupDatabase(dataURL, websiteURL, filePath, absDataDir, dbFilePath)
	if err != nil {
		log.Printf("Error during database setup: %v", err)
		return err
	}
	log.Println("Database setup completed successfully.")
	return nil
}
