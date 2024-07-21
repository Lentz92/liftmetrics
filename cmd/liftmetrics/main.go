package main

import (
	"liftmetrics/internal/app"
	"liftmetrics/internal/db"
	"liftmetrics/internal/services"
	"log"
	"os"
	"path/filepath"
)

const (
	dataURL    = "https://openpowerlifting.gitlab.io/opl-csv/files/openipf-latest.zip"
	websiteURL = "https://openpowerlifting.gitlab.io/opl-csv/bulk-csv.html"
	dataDir    = "../../data"
	zipFile    = "openipf-latest.zip"
	dbName     = "openipf.db"
)

func main() {
	// Set up logging to include date, time, and file information
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Get the absolute path for the data directory
	absDataDir, err := filepath.Abs(dataDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path for data directory: %v", err)
	}

	// Set up file paths
	filePath := filepath.Join(absDataDir, zipFile)
	dbDir := filepath.Join(absDataDir, "db")
	dbFilePath := filepath.Join(dbDir, dbName)

	// Create necessary directories
	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create directories: %v", err)
	}

	// Set up or update the database
	if err := setupDatabase(dataURL, websiteURL, filePath, absDataDir, dbFilePath); err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}

	// Open the existing database
	database, err := db.OpenDatabase(dbFilePath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	// Create and run the Fyne application
	liftApp := app.New()
	liftApp.Run()
}

// setupDatabase handles the process of setting up or updating the database
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
