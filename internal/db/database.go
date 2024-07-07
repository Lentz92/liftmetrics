package db

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/schollz/progressbar/v3"
	"os"
)

type Record struct {
	ID               string
	Name             string
	Sex              string
	Event            string
	Equipment        string
	Age              float64
	AgeClass         string
	BirthYearClass   string
	Division         string
	BodyweightKg     float64
	WeightClassKg    string
	Squat1Kg         float64
	Squat2Kg         float64
	Squat3Kg         float64
	Squat4Kg         float64
	Best3SquatKg     float64
	Bench1Kg         float64
	Bench2Kg         float64
	Bench3Kg         float64
	Bench4Kg         float64
	Best3BenchKg     float64
	Deadlift1Kg      float64
	Deadlift2Kg      float64
	Deadlift3Kg      float64
	Deadlift4Kg      float64
	Best3DeadliftKg  float64
	TotalKg          float64
	Place            string
	Dots             float64
	Wilks            float64
	Glossbrenner     float64
	Goodlift         float64
	Tested           string
	Country          string
	State            string
	Federation       string
	ParentFederation string
	Date             string
	MeetCountry      string
	MeetState        string
	MeetTown         string
	MeetName         string
	Sanctioned       string
}

func CreateDatabase(dbFilePath string, deleteExisting bool) (*sql.DB, error) {
	if deleteExisting {
		if err := os.Remove(dbFilePath); err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to remove existing database file: %w", err)
		}
	}

	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := createTables(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
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
        Sanctioned TEXT,
        SuccessfulSquatAttempts INTEGER DEFAULT 0,
        SuccessfulBenchAttempts INTEGER DEFAULT 0,
        SuccessfulDeadliftAttempts INTEGER DEFAULT 0,
        TotalSuccessfulAttempts INTEGER DEFAULT 0
    );

    CREATE TABLE IF NOT EXISTS lift_diffs (
        ID TEXT PRIMARY KEY,
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
        FOREIGN KEY (ID) REFERENCES Records(ID)
    );`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

func PopulateDatabase(db *sql.DB, records []*Record) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
    INSERT INTO records (
        ID, Name, Sex, Event, Equipment, Age, AgeClass, BirthYearClass, Division, BodyweightKg, WeightClassKg,
        Squat1Kg, Squat2Kg, Squat3Kg, Squat4Kg, Best3SquatKg, Bench1Kg, Bench2Kg, Bench3Kg, Bench4Kg, Best3BenchKg,
        Deadlift1Kg, Deadlift2Kg, Deadlift3Kg, Deadlift4Kg, Best3DeadliftKg, TotalKg, Place, Dots, Wilks, Glossbrenner,
        Goodlift, Tested, Country, State, Federation, ParentFederation, Date, MeetCountry, MeetState, MeetTown, MeetName, Sanctioned
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	bar := progressbar.NewOptions(len(records), progressbar.OptionSetPredictTime(false))

	for _, record := range records {
		record.ID = uuid.New().String()
		if _, err := stmt.Exec(
			record.ID, record.Name, record.Sex, record.Event, record.Equipment, record.Age, record.AgeClass, record.BirthYearClass, record.Division,
			record.BodyweightKg, record.WeightClassKg, record.Squat1Kg, record.Squat2Kg, record.Squat3Kg, record.Squat4Kg, record.Best3SquatKg,
			record.Bench1Kg, record.Bench2Kg, record.Bench3Kg, record.Bench4Kg, record.Best3BenchKg, record.Deadlift1Kg, record.Deadlift2Kg,
			record.Deadlift3Kg, record.Deadlift4Kg, record.Best3DeadliftKg, record.TotalKg, record.Place, record.Dots, record.Wilks,
			record.Glossbrenner, record.Goodlift, record.Tested, record.Country, record.State, record.Federation, record.ParentFederation,
			record.Date, record.MeetCountry, record.MeetState, record.MeetTown, record.MeetName, record.Sanctioned,
		); err != nil {
			return fmt.Errorf("failed to insert record: %w", err)
		}

		if err := bar.Add(1); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func CalculateSuccessfulAttempts(db *sql.DB) error {
	updateAttemptsSQL := `
    UPDATE Records SET 
        SuccessfulSquatAttempts = (CASE WHEN Squat1Kg > 0 THEN 1 ELSE 0 END) +
                                  (CASE WHEN Squat2Kg > 0 THEN 1 ELSE 0 END) +
                                  (CASE WHEN Squat3Kg > 0 THEN 1 ELSE 0 END),
        SuccessfulBenchAttempts = (CASE WHEN Bench1Kg > 0 THEN 1 ELSE 0 END) +
                                  (CASE WHEN Bench2Kg > 0 THEN 1 ELSE 0 END) +
                                  (CASE WHEN Bench3Kg > 0 THEN 1 ELSE 0 END),
        SuccessfulDeadliftAttempts = (CASE WHEN Deadlift1Kg > 0 THEN 1 ELSE 0 END) +
                                     (CASE WHEN Deadlift2Kg > 0 THEN 1 ELSE 0 END) +
                                     (CASE WHEN Deadlift3Kg > 0 THEN 1 ELSE 0 END),
        TotalSuccessfulAttempts = SuccessfulSquatAttempts + SuccessfulBenchAttempts + SuccessfulDeadliftAttempts`

	if _, err := db.Exec(updateAttemptsSQL); err != nil {
		return fmt.Errorf("failed to calculate successful attempts: %w", err)
	}

	return nil
}

func CalculateRelativeDiff(db *sql.DB) error {
	insertSQL := `
    INSERT INTO lift_diffs (ID, Name, Date, Squat1Perc, Squat2Perc, Squat3Perc, Bench1Perc, Bench2Perc, Bench3Perc, Deadlift1Perc, Deadlift2Perc, Deadlift3Perc)
    SELECT 
        ID,
        Name,
        Date,
        ABS(Squat1Kg) / MAX(ABS(Squat1Kg), ABS(Squat2Kg), ABS(Squat3Kg)) * 100 AS Squat1Perc,
        ABS(Squat2Kg) / MAX(ABS(Squat1Kg), ABS(Squat2Kg), ABS(Squat3Kg)) * 100 AS Squat2Perc,
        ABS(Squat3Kg) / MAX(ABS(Squat1Kg), ABS(Squat2Kg), ABS(Squat3Kg)) * 100 AS Squat3Perc,
        ABS(Bench1Kg) / MAX(ABS(Bench1Kg), ABS(Bench2Kg), ABS(Bench3Kg)) * 100 AS Bench1Perc,
        ABS(Bench2Kg) / MAX(ABS(Bench1Kg), ABS(Bench2Kg), ABS(Bench3Kg)) * 100 AS Bench2Perc,
        ABS(Bench3Kg) / MAX(ABS(Bench1Kg), ABS(Bench2Kg), ABS(Bench3Kg)) * 100 AS Bench3Perc,
        ABS(Deadlift1Kg) / MAX(ABS(Deadlift1Kg), ABS(Deadlift2Kg), ABS(Deadlift3Kg)) * 100 AS Deadlift1Perc,
        ABS(Deadlift2Kg) / MAX(ABS(Deadlift1Kg), ABS(Deadlift2Kg), ABS(Deadlift3Kg)) * 100 AS Deadlift2Perc,
        ABS(Deadlift3Kg) / MAX(ABS(Deadlift1Kg), ABS(Deadlift2Kg), ABS(Deadlift3Kg)) * 100 AS Deadlift3Perc
    FROM
        Records;`

	if _, err := db.Exec(insertSQL); err != nil {
		return fmt.Errorf("failed to calculate relative differences: %w", err)
	}

	return nil
}
