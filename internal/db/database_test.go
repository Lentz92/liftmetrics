package db

import (
	_ "github.com/mattn/go-sqlite3" // Import for side-effects to register the SQLite driver
	"os"
	"path/filepath"
	"testing"
)

// TestCreateDatabase tests the CreateDatabase function
func TestCreateDatabase(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	dbFilePath := filepath.Join(tempDir, "test.db")

	// Call the CreateDatabase function
	database, err := CreateDatabase(dbFilePath, false)
	if err != nil {
		t.Fatalf("CreateDatabase failed: %v", err)
	}
	defer database.Close()

	// Check if the database file is created
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		t.Fatalf("Database file does not exist")
	}

	// Check if the table 'records' exists
	tableCheckQuery := `SELECT name FROM sqlite_master WHERE type='table' AND name='records';`
	var tableName string
	err = database.QueryRow(tableCheckQuery).Scan(&tableName)
	if err != nil {
		t.Fatalf("Failed to find table 'records': %v", err)
	}
	if tableName != "records" {
		t.Fatalf("Table 'records' does not exist")
	}
}

// TestPopulateDatabase tests the PopulateDatabase function
func TestPopulateDatabase(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	dbFilePath := filepath.Join(tempDir, "test.db")

	// Call the CreateDatabase function
	database, err := CreateDatabase(dbFilePath, false)
	if err != nil {
		t.Fatalf("CreateDatabase failed: %v", err)
	}
	defer database.Close()

	// Create sample records to populate the database
	records := []*Record{
		{
			Name:             "John Doe",
			Sex:              "M",
			Event:            "Squat",
			Equipment:        "Raw",
			Age:              25,
			AgeClass:         "24-26",
			BirthYearClass:   "1995",
			Division:         "Open",
			BodyweightKg:     90.0,
			WeightClassKg:    "90",
			Squat1Kg:         200.0,
			Squat2Kg:         210.0,
			Squat3Kg:         220.0,
			Best3SquatKg:     220.0,
			Bench1Kg:         120.0,
			Bench2Kg:         125.0,
			Bench3Kg:         130.0,
			Best3BenchKg:     130.0,
			Deadlift1Kg:      230.0,
			Deadlift2Kg:      240.0,
			Deadlift3Kg:      250.0,
			Best3DeadliftKg:  250.0,
			TotalKg:          600.0,
			Place:            "1st",
			Dots:             500.0,
			Wilks:            450.0,
			Glossbrenner:     400.0,
			Goodlift:         550.0,
			Tested:           "Yes",
			Country:          "USA",
			State:            "CA",
			Federation:       "USAPL",
			ParentFederation: "IPF",
			Date:             "2023-07-04",
			MeetCountry:      "USA",
			MeetState:        "CA",
			MeetTown:         "San Diego",
			MeetName:         "Summer Classic",
			Sanctioned:       "Yes",
		},
	}

	// Call the PopulateDatabase function
	err = PopulateDatabase(database, records)
	if err != nil {
		t.Fatalf("PopulateDatabase failed: %v", err)
	}

	// Verify the records are inserted correctly
	row := database.QueryRow("SELECT COUNT(*) FROM records;")
	var count int
	if err := row.Scan(&count); err != nil {
		t.Fatalf("Failed to count records: %v", err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 record, got %d", count)
	}

}
