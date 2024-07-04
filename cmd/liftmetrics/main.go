package main

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"liftmetrics/internal/db"
	"liftmetrics/internal/services"
	"liftmetrics/pkg"
	"log"
	"os"
	"path/filepath"
)

func main() {
	url := "https://openpowerlifting.gitlab.io/opl-csv/files/openipf-latest.zip"
	dataDir := filepath.Join("../../data")
	filePath := filepath.Join(dataDir, "openipf-latest.zip")
	dbFilePath := filepath.Join(dataDir, "db", "openipf.db")

	// Create data and db directories if they don't exist
	err := os.MkdirAll(filepath.Join(dataDir, "db"), os.ModePerm)
	if err != nil {
		return
	}

	// Step 1: Download the ZIP file
	err = services.DownloadFile(url, filePath)
	if err != nil {
		log.Fatalf("Error downloading file: %v\n", err)
	}

	// Step 2: Extract the CSV file from the ZIP
	_, err = services.ExtractCSVFromZip(filePath, dataDir)
	if err != nil {
		log.Fatalf("Failed to extract CSV: %v", err)
	}

	// Step 3: Find the CSV file
	csvFilePath, err := pkg.FindCSVFile(dataDir)
	if err != nil {
		log.Fatalf("Failed to find CSV file: %v", err)
	}

	// Step 4: Open and parse the CSV file
	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer csvFile.Close()

	var records []*db.Record
	if err := gocsv.UnmarshalFile(csvFile, &records); err != nil {
		log.Fatalf("Failed to parse CSV file: %v", err)
	}

	// Step 5: Create the database
	fmt.Println("\nSetting up the database this might take a while!")
	database, err := db.CreateDatabase(dbFilePath, true)
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer database.Close()

	// Step 6: Populate the database with records
	err = db.PopulateDatabase(database, records)
	if err != nil {
		log.Fatalf("Failed to populate database: %v", err)
	}

	// Step 7: Calculate Successful attempts for the database
	fmt.Println("\nCalculating new metrics for the database!")
	err = db.CalculateSuccessfulAttempts(dbFilePath)
	if err != nil {
		log.Fatalf("Failed to update lift attempts in database: %v", err)
	}

	//  Step 8: Create new table and calculate relative diff in attempts
	err = db.LiftDiffTable(database)
	if err != nil {
		log.Fatalf("Failed to ensure lift_diffs table exists: %v", err)
	}

	err = db.CalculateRelativeDiff(database)
	if err != nil {
		log.Fatalf("Failed to calculate and insert relative differences: %v", err)
	}

	fmt.Println("\nAll setup and calculations complete!")

}
