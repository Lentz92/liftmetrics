package main

import (
	"context"
	"encoding/json"
	"fmt"
	"liftmetrics/internal/db"
	"liftmetrics/internal/services"
	"liftmetrics/internal/web"
	"log"
	"os"
	"path/filepath"
	"time"
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

	// Get the current working directory
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}
	log.Printf("Current working directory: %s", pwd)

	// Construct the path to interface.html
	interfaceHTMLPath := filepath.Join(pwd, "..", "..", "internal", "web", "interface.html")
	log.Printf("Constructed interface.html path: %s", interfaceHTMLPath)

	// Check if the file exists
	if _, err := os.Stat(interfaceHTMLPath); os.IsNotExist(err) {
		log.Fatalf("interface.html does not exist at path: %s", interfaceHTMLPath)
	}

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

	// Load lifter names
	lifterNames, err := loadLifterNames(filepath.Join(absDataDir, "lifters.json"))
	if err != nil {
		log.Fatalf("Failed to load lifter names: %v", err)
	}

	// Log the number of lifter names loaded
	log.Printf("Loaded %d lifter names", len(lifterNames))

	// Create and start the HTTP server
	server := web.NewServer(lifterNames, database, interfaceHTMLPath)
	log.Fatal(server.Start())
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

	// Generate JSON file of lifter names
	database, err := db.OpenDatabase(dbFilePath)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer database.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	err = db.GenerateLifterJSON(ctx, database, absDataDir)
	if err != nil {
		return fmt.Errorf("generating lifter JSON: %w", err)
	}
	log.Println("Lifter JSON file generated successfully.")

	return nil
}

func loadLifterNames(filePath string) ([]string, error) {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var lifterNames []string
	err = json.Unmarshal(jsonData, &lifterNames)
	if err != nil {
		return nil, err
	}

	return lifterNames, nil
}
