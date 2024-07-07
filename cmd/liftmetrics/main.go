package main

import (
	"log"
	"os"
	"path/filepath"

	"liftmetrics/internal/services"
)

const (
	dataURL = "https://openpowerlifting.gitlab.io/opl-csv/files/openipf-latest.zip"
	dataDir = "data"
	zipFile = "openipf-latest.zip"
	dbName  = "openipf.db"
)

func main() {
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

	// Download, extract, and process data
	if err := services.SetupDatabase(dataURL, filePath, absDataDir, dbFilePath); err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}

	log.Println("Database setup completed successfully.")
}
