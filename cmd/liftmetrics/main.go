package main

import (
	"liftmetrics/internal/db"
	"liftmetrics/internal/server"
	"liftmetrics/internal/services"
	"log"
	"os"
	"path/filepath"
)

const (
	dataURL    = "https://openpowerlifting.gitlab.io/opl-csv/files/openipf-latest.zip"
	websiteURL = "https://openpowerlifting.gitlab.io/opl-csv/bulk-csv.html"
	dataDir    = "../../data" // Changed to a relative path
	zipFile    = "openipf-latest.zip"
	dbName     = "openipf.db"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	absDataDir, err := filepath.Abs(dataDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path for data directory: %v", err)
	}

	filePath := filepath.Join(absDataDir, zipFile)
	dbDir := filepath.Join(absDataDir, "db")
	dbFilePath := filepath.Join(dbDir, dbName)

	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create directories: %v", err)
	}

	if err := setupDatabase(dataURL, websiteURL, filePath, absDataDir, dbFilePath); err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}

	// Open the existing database
	database, err := db.OpenDatabase(dbFilePath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	// Initialize the server with the existing database connection
	srv, err := server.NewServer(database)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	log.Println("Starting server on http://localhost:8080")
	if err := srv.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

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
