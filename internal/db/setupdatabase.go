package db

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // Import for side-effects to register the SQLite driver
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"strings"
)

type Record struct {
	ID               string  // This field will hold the generated UUID
	Name             string  `csv:"Name"`
	Sex              string  `csv:"Sex"`
	Event            string  `csv:"Event"`
	Equipment        string  `csv:"Equipment"`
	Age              float64 `csv:"Age"`
	AgeClass         string  `csv:"AgeClass"`
	BirthYearClass   string  `csv:"BirthYearClass"`
	Division         string  `csv:"Division"`
	BodyweightKg     float64 `csv:"BodyweightKg"`
	WeightClassKg    string  `csv:"WeightClassKg"`
	Squat1Kg         float64 `csv:"Squat1Kg"`
	Squat2Kg         float64 `csv:"Squat2Kg"`
	Squat3Kg         float64 `csv:"Squat3Kg"`
	Squat4Kg         float64 `csv:"Squat4Kg"`
	Best3SquatKg     float64 `csv:"Best3SquatKg"`
	Bench1Kg         float64 `csv:"Bench1Kg"`
	Bench2Kg         float64 `csv:"Bench2Kg"`
	Bench3Kg         float64 `csv:"Bench3Kg"`
	Bench4Kg         float64 `csv:"Bench4Kg"`
	Best3BenchKg     float64 `csv:"Best3BenchKg"`
	Deadlift1Kg      float64 `csv:"Deadlift1Kg"`
	Deadlift2Kg      float64 `csv:"Deadlift2Kg"`
	Deadlift3Kg      float64 `csv:"Deadlift3Kg"`
	Deadlift4Kg      float64 `csv:"Deadlift4Kg"`
	Best3DeadliftKg  float64 `csv:"Best3DeadliftKg"`
	TotalKg          float64 `csv:"TotalKg"`
	Place            string  `csv:"Place"`
	Dots             float64 `csv:"Dots"`
	Wilks            float64 `csv:"Wilks"`
	Glossbrenner     float64 `csv:"Glossbrenner"`
	Goodlift         float64 `csv:"Goodlift"`
	Tested           string  `csv:"Tested"`
	Country          string  `csv:"Country"`
	State            string  `csv:"State"`
	Federation       string  `csv:"Federation"`
	ParentFederation string  `csv:"ParentFederation"`
	Date             string  `csv:"Date"`
	MeetCountry      string  `csv:"MeetCountry"`
	MeetState        string  `csv:"MeetState"`
	MeetTown         string  `csv:"MeetTown"`
	MeetName         string  `csv:"MeetName"`
	Sanctioned       string  `csv:"Sanctioned"`
}

// CreateDatabase Creating database
func CreateDatabase(dbFilePath string, deleteExisting bool) (*sql.DB, error) {
	if deleteExisting {
		if _, err := os.Stat(dbFilePath); err == nil {
			// File exists, remove it
			err := os.Remove(dbFilePath)
			if err != nil {
				return nil, fmt.Errorf("failed to remove existing database file: %w", err)
			}
		} else if !os.IsNotExist(err) {
			// Some other error
			return nil, fmt.Errorf("failed to check if database file exists: %w", err)
		}
	}

	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS records (
        ID TEXT PRIMARY KEY,
        Name TEXT,
        Sex TEXT,
        Event TEXT,
        Equipment TEXT,
        Age REAL,
        AgeClass TEXT,
        BirthYearClass TEXT,
        Division TEXT,
        BodyweightKg REAL,
        WeightClassKg TEXT,
        Squat1Kg REAL,
        Squat2Kg REAL,
        Squat3Kg REAL,
        Squat4Kg REAL,
        Best3SquatKg REAL,
        Bench1Kg REAL,
        Bench2Kg REAL,
        Bench3Kg REAL,
        Bench4Kg REAL,
        Best3BenchKg REAL,
        Deadlift1Kg REAL,
        Deadlift2Kg REAL,
        Deadlift3Kg REAL,
        Deadlift4Kg REAL,
        Best3DeadliftKg REAL,
        TotalKg REAL,
        Place TEXT,
        Dots REAL,
        Wilks REAL,
        Glossbrenner REAL,
        Goodlift REAL,
        Tested TEXT,
        Country TEXT,
        State TEXT,
        Federation TEXT,
        ParentFederation TEXT,
        Date TEXT,
        MeetCountry TEXT,
        MeetState TEXT,
        MeetTown TEXT,
        MeetName TEXT,
        Sanctioned TEXT
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return db, nil
}

// PopulateDatabase inserts records into the database
func PopulateDatabase(db *sql.DB, records []*Record) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	insertSQL := `
    INSERT INTO records (
        ID, Name, Sex, Event, Equipment, Age, AgeClass, BirthYearClass, Division, BodyweightKg, WeightClassKg,
        Squat1Kg, Squat2Kg, Squat3Kg, Squat4Kg, Best3SquatKg, Bench1Kg, Bench2Kg, Bench3Kg, Bench4Kg, Best3BenchKg,
        Deadlift1Kg, Deadlift2Kg, Deadlift3Kg, Deadlift4Kg, Best3DeadliftKg, TotalKg, Place, Dots, Wilks, Glossbrenner,
        Goodlift, Tested, Country, State, Federation, ParentFederation, Date, MeetCountry, MeetState, MeetTown, MeetName, Sanctioned
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	stmt, err := tx.Prepare(insertSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	// Initialize progress bar
	bar := progressbar.NewOptions(len(records), progressbar.OptionSetPredictTime(false))

	for _, record := range records {
		record.ID = uuid.New().String()
		_, err = stmt.Exec(
			record.ID, record.Name, record.Sex, record.Event, record.Equipment, record.Age, record.AgeClass, record.BirthYearClass, record.Division,
			record.BodyweightKg, record.WeightClassKg, record.Squat1Kg, record.Squat2Kg, record.Squat3Kg, record.Squat4Kg, record.Best3SquatKg,
			record.Bench1Kg, record.Bench2Kg, record.Bench3Kg, record.Bench4Kg, record.Best3BenchKg, record.Deadlift1Kg, record.Deadlift2Kg,
			record.Deadlift3Kg, record.Deadlift4Kg, record.Best3DeadliftKg, record.TotalKg, record.Place, record.Dots, record.Wilks,
			record.Glossbrenner, record.Goodlift, record.Tested, record.Country, record.State, record.Federation, record.ParentFederation,
			record.Date, record.MeetCountry, record.MeetState, record.MeetTown, record.MeetName, record.Sanctioned,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert record: %w", err)
		}

		// Update progress bar
		err := bar.Add(1)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func CalculateSuccessfulAttempts(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Add necessary columns if they do not exist
	columnsToAdd := []string{
		"SuccessfulSquatAttempts INTEGER DEFAULT 0",
		"SuccessfulBenchAttempts INTEGER DEFAULT 0",
		"SuccessfulDeadliftAttempts INTEGER DEFAULT 0",
		"TotalSuccessfulAttempts INTEGER DEFAULT 0",
	}
	for _, col := range columnsToAdd {
		alterTableSQL := fmt.Sprintf("ALTER TABLE Records ADD COLUMN %s", col)
		if _, err = db.Exec(alterTableSQL); err != nil {
			if !strings.Contains(err.Error(), "duplicate column name") {
				return fmt.Errorf("failed to add column '%s': %w", col, err)
			}
		}
	}

	// Update individual successful attempt columns
	updateAttemptsSQL := `UPDATE Records SET 
        SuccessfulSquatAttempts = (CASE WHEN Squat1Kg > 0 THEN 1 ELSE 0 END) +
                                  (CASE WHEN Squat2Kg > 0 THEN 1 ELSE 0 END) +
                                  (CASE WHEN Squat3Kg > 0 THEN 1 ELSE 0 END),
        SuccessfulBenchAttempts = (CASE WHEN Bench1Kg > 0 THEN 1 ELSE 0 END) +
                                  (CASE WHEN Bench2Kg > 0 THEN 1 ELSE 0 END) +
                                  (CASE WHEN Bench3Kg > 0 THEN 1 ELSE 0 END),
        SuccessfulDeadliftAttempts = (CASE WHEN Deadlift1Kg > 0 THEN 1 ELSE 0 END) +
                                     (CASE WHEN Deadlift2Kg > 0 THEN 1 ELSE 0 END) +
                                     (CASE WHEN Deadlift3Kg > 0 THEN 1 ELSE 0 END)`
	_, err = db.Exec(updateAttemptsSQL)
	if err != nil {
		log.Printf("Failed to execute updateAttemptsSQL: %s\nError: %s\n", updateAttemptsSQL, err)
		return err
	}

	// Update the total successful attempts column
	updateTotalSQL := `UPDATE Records SET 
        TotalSuccessfulAttempts = SuccessfulSquatAttempts + SuccessfulBenchAttempts + SuccessfulDeadliftAttempts`
	_, err = db.Exec(updateTotalSQL)
	if err != nil {
		log.Printf("Failed to execute updateTotalSQL: %s\nError: %s\n", updateTotalSQL, err)
		return err
	}

	return nil
}

// LiftDiffTable creates the lift_diffs table if it does not already exist
func LiftDiffTable(db *sql.DB) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS lift_diffs (
    	ID TEXT,
        Name TEXT,
        Date TEXT,
        Squat1Perc REAL,
        Squat2Perc REAL,
        Squat3Perc REAL,
        Bench1Perc REAL,
        Bench2Perc REAL,
        Bench3Perc REAL,
        Deadlift1Perc REAL,
        Deadlift2Perc REAL,
        Deadlift3Perc REAL,
		PRIMARY KEY (ID),
    	FOREIGN KEY (ID) REFERENCES Records(ID)
    );`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}
	return nil
}

// CalculateRelativeDiff calculates the relative differences between lift attempts and inserts them into the lift_diffs table
func CalculateRelativeDiff(db *sql.DB) error {
	insertSQL := `INSERT INTO lift_diffs (ID, Name, Date, Squat1Perc, Squat2Perc, Squat3Perc, Bench1Perc, Bench2Perc, Bench3Perc, Deadlift1Perc, Deadlift2Perc, Deadlift3Perc)
    SELECT 
        ID,
        Name,
        Date,
        ABS(Squat1Kg) / (CASE WHEN ABS(Squat1Kg) >= ABS(Squat2Kg) AND ABS(Squat1Kg) >= ABS(Squat3Kg) THEN ABS(Squat1Kg)
                              WHEN ABS(Squat2Kg) >= ABS(Squat1Kg) AND ABS(Squat2Kg) >= ABS(Squat3Kg) THEN ABS(Squat2Kg)
                              ELSE ABS(Squat3Kg) END) * 100 AS Squat1Perc,
        ABS(Squat2Kg) / (CASE WHEN ABS(Squat1Kg) >= ABS(Squat2Kg) AND ABS(Squat1Kg) >= ABS(Squat3Kg) THEN ABS(Squat1Kg)
                              WHEN ABS(Squat2Kg) >= ABS(Squat1Kg) AND ABS(Squat2Kg) >= ABS(Squat3Kg) THEN ABS(Squat2Kg)
                              ELSE ABS(Squat3Kg) END) * 100 AS Squat2Perc,
        ABS(Squat3Kg) / (CASE WHEN ABS(Squat1Kg) >= ABS(Squat2Kg) AND ABS(Squat1Kg) >= ABS(Squat3Kg) THEN ABS(Squat1Kg)
                              WHEN ABS(Squat2Kg) >= ABS(Squat1Kg) AND ABS(Squat2Kg) >= ABS(Squat3Kg) THEN ABS(Squat2Kg)
                              ELSE ABS(Squat3Kg) END) * 100 AS Squat3Perc,
        ABS(Bench1Kg) / (CASE WHEN ABS(Bench1Kg) >= ABS(Bench2Kg) AND ABS(Bench1Kg) >= ABS(Bench3Kg) THEN ABS(Bench1Kg)
                              WHEN ABS(Bench2Kg) >= ABS(Bench1Kg) AND ABS(Bench2Kg) >= ABS(Bench3Kg) THEN ABS(Bench2Kg)
                              ELSE ABS(Bench3Kg) END) * 100 AS Bench1Perc,
        ABS(Bench2Kg) / (CASE WHEN ABS(Bench1Kg) >= ABS(Bench2Kg) AND ABS(Bench1Kg) >= ABS(Bench3Kg) THEN ABS(Bench1Kg)
                              WHEN ABS(Bench2Kg) >= ABS(Bench1Kg) AND ABS(Bench2Kg) >= ABS(Bench3Kg) THEN ABS(Bench2Kg)
                              ELSE ABS(Bench3Kg) END) * 100 AS Bench2Perc,
        ABS(Bench3Kg) / (CASE WHEN ABS(Bench1Kg) >= ABS(Bench2Kg) AND ABS(Bench1Kg) >= ABS(Bench3Kg) THEN ABS(Bench1Kg)
                              WHEN ABS(Bench2Kg) >= ABS(Bench1Kg) AND ABS(Bench2Kg) >= ABS(Bench3Kg) THEN ABS(Bench2Kg)
                              ELSE ABS(Bench3Kg) END) * 100 AS Bench3Perc,
        ABS(Deadlift1Kg) / (CASE WHEN ABS(Deadlift1Kg) >= ABS(Deadlift2Kg) AND ABS(Deadlift1Kg) >= ABS(Deadlift3Kg) THEN ABS(Deadlift1Kg)
                                WHEN ABS(Deadlift2Kg) >= ABS(Deadlift1Kg) AND ABS(Deadlift2Kg) >= ABS(Deadlift3Kg) THEN ABS(Deadlift2Kg)
                                ELSE ABS(Deadlift3Kg) END) * 100 AS Deadlift1Perc,
        ABS(Deadlift2Kg) / (CASE WHEN ABS(Deadlift1Kg) >= ABS(Deadlift2Kg) AND ABS(Deadlift1Kg) >= ABS(Deadlift3Kg) THEN ABS(Deadlift1Kg)
                                WHEN ABS(Deadlift2Kg) >= ABS(Deadlift1Kg) AND ABS(Deadlift2Kg) >= ABS(Deadlift3Kg) THEN ABS(Deadlift2Kg)
                                ELSE ABS(Deadlift3Kg) END) * 100 AS Deadlift2Perc,
        ABS(Deadlift3Kg) / (CASE WHEN ABS(Deadlift1Kg) >= ABS(Deadlift2Kg) AND ABS(Deadlift1Kg) >= ABS(Deadlift3Kg) THEN ABS(Deadlift1Kg)
                                WHEN ABS(Deadlift2Kg) >= ABS(Deadlift1Kg) AND ABS(Deadlift2Kg) >= ABS(Deadlift3Kg) THEN ABS(Deadlift2Kg)
                                ELSE ABS(Deadlift3Kg) END) * 100 AS Deadlift3Perc
    FROM
        Records;`
	_, err := db.Exec(insertSQL)
	if err != nil {
		log.Printf("Failed to execute insertSQL: %s\nError: %s\n", insertSQL, err)
		return err
	}
	return nil
}
