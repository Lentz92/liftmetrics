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
        Sanctioned TEXT
    );

    CREATE TABLE IF NOT EXISTS lifter_metrics (
        ID TEXT,
        Name TEXT,
        Date TEXT,
        SuccessfulSquatAttempts INTEGER DEFAULT 0,
        SuccessfulBenchAttempts INTEGER DEFAULT 0,
        SuccessfulDeadliftAttempts INTEGER DEFAULT 0,
        TotalSuccessfulAttempts INTEGER DEFAULT 0,
        Squat1Perc REAL,
        Squat2Perc REAL,
        Squat3Perc REAL,
        Bench1Perc REAL,
        Bench2Perc REAL,
        Bench3Perc REAL,
        Deadlift1Perc REAL,
        Deadlift2Perc REAL,
        Deadlift3Perc REAL,
        Squat1To2Kg REAL,
        Squat2To3Kg REAL,
        Bench1To2Kg REAL,
        Bench2To3Kg REAL,
        Deadlift1To2Kg REAL,
        Deadlift2To3Kg REAL,
        PRIMARY KEY (ID, Date),
        FOREIGN KEY (ID) REFERENCES records(ID)
    );

    CREATE TABLE IF NOT EXISTS aggregated_metrics_sbd (
        Name TEXT PRIMARY KEY,
        AvgSuccessfulSquatAttempts REAL,
        AvgSuccessfulBenchAttempts REAL,
        AvgSuccessfulDeadliftAttempts REAL,
        AvgTotalSuccessfulAttempts REAL,
        AvgSquat1Perc REAL,
        AvgSquat2Perc REAL,
        AvgSquat3Perc REAL,
        AvgBench1Perc REAL,
        AvgBench2Perc REAL,
        AvgBench3Perc REAL,
        AvgDeadlift1Perc REAL,
        AvgDeadlift2Perc REAL,
        AvgDeadlift3Perc REAL,
        AvgSquat1To2Kg REAL,
        AvgSquat2To3Kg REAL,
        AvgBench1To2Kg REAL,
        AvgBench2To3Kg REAL,
        AvgDeadlift1To2Kg REAL,
        AvgDeadlift2To3Kg REAL
    );

    CREATE TABLE IF NOT EXISTS aggregated_metrics_bench (
        Name TEXT PRIMARY KEY,
        AvgSuccessfulBenchAttempts REAL,
        AvgBench1Perc REAL,
        AvgBench2Perc REAL,
        AvgBench3Perc REAL,
        AvgBench1To2Kg REAL,
        AvgBench2To3Kg REAL
    );

    CREATE INDEX IF NOT EXISTS idx_lifter_metrics_name_date ON lifter_metrics(Name, Date);
    CREATE INDEX IF NOT EXISTS idx_records_name_date ON records(Name, Date);
    `

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

func OpenDatabase(dbPath string) (*sql.DB, error) {
	return sql.Open("sqlite3", dbPath)
}
